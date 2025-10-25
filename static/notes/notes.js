const apiUrl = `/api`;
const appUrl = `/app`;
let jwtToken = "";
let userID = "";
let username = "";
let userEmail = "";
let noteID = "";


// Check for existing token on page load
window.onload = () => {
    jwtToken = localStorage.getItem("jwtToken"); // Retrieve token from localStorage

    if (jwtToken) {
        const isNotesPage = window.location.pathname.endsWith('notes');
        const isEditPage = window.location.pathname.endsWith('note');

        if (isNotesPage) {
            loadNotes();
        } else if (isEditPage) {
            noteID = localStorage.getItem("note-id")
            if (noteID) {
                loadNoteForEditing(noteID);
            } else {
                console.log("Missing note's ID. Going back to previous page.")
                window.location.href = window.history.back();
            }
            
        }

    } else {
        console.log("Missing access token. Going back to previous page.")
        window.location.href = window.history.back();
    }
};

// Handle logout
document.getElementById("logout-button").addEventListener("click", () => {
    logout();
});

// Handle note form submission
document.getElementById("note-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const title = document.getElementById("note-title").value;
    const content = document.getElementById("note-content").value;

    const response = await secureFetch(`${apiUrl}/notes`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ title, content })
    });

    document.getElementById("note-form").reset(); // Reset form
    loadNotes(); // Reload notes after adding new note
});

async function loadNoteForEditing(id) {
    const response = await secureFetch(`${apiUrl}/notes/${noteID}`, {});

    const note = await response.json();
    renderEditableNote(note);

}

// Load notes from the API
async function loadNotes() {
    const response = await secureFetch(`${apiUrl}/notes`, {});

    notesHeader = document.getElementById("notes-h2");
    username = localStorage.getItem("username");
    notesHeader.textContent = `${username}'s notes:`;
    let notes = await response.json();
    const notesList = document.getElementById("notes-list");
    notesList.innerHTML = ""; // Clear the existing list

    if (notes.length === 0) {
        // If there are no notes, display a message
        const noNotesMessage = document.createElement("li");
        noNotesMessage.classList.add("no-notes");
        noNotesMessage.textContent = "No notes available.";
        notesList.appendChild(noNotesMessage);
    } else {
        // If there are notes, display them
        notes.forEach(note => {
            
            const li = createNoteElement(note)
            
            const editButton = getButtonElement("edit-button", "Edit", note.id, editNote);
            const deleteButton = getButtonElement("delete-button", "Delete", note.id, deleteNote);

            // Adding all the elements to the list
            li.appendChild(editButton);
            li.appendChild(deleteButton);
            notesList.appendChild(li);
        });
    }
}

function createNoteElement(note) {
    const li = document.createElement('li');

    // Create a container for the editable text and note ID
    const noteContentDiv = document.createElement('div');
    noteContentDiv.classList.add('note-content-container');
    noteContentDiv.dataset.noteId = note.id;

    // Create editable elements for title and content
    const titleDiv = document.createElement('div');
    titleDiv.classList.add('note-title');
    titleDiv.textContent = note.title;

    const contentDiv = document.createElement('div');
    contentDiv.classList.add('note-body');
    contentDiv.textContent = note.content;

    // Append all parts to the list item
    noteContentDiv.appendChild(titleDiv);
    noteContentDiv.appendChild(contentDiv);
    li.appendChild(noteContentDiv);
    
    return li;
}

function getButtonElement(buttonName, buttonText, noteID, func) {
    const button = document.createElement(buttonName);
    const text = document.createElement(buttonName + "-text");
    button.classList.add(buttonName);
    text.classList.add(buttonName + "-text");
    text.textContent = buttonText;
    button.appendChild(text);
    button.setAttribute("noteID", noteID);
    button.addEventListener('click', (event) => {
        const clickedButton = event.target;
        const noteID = clickedButton.getAttribute("noteID");
        func(noteID);
    });
    return button;
}

function renderEditableNote(note) {
    const container = document.getElementById('single-note-container');
    const li = createNoteElement(note);

    // Get the title and body elements within the new list item
    const titleDiv = li.querySelector('.note-title');
    const contentDiv = li.querySelector('.note-body');

    // Make them editable
    titleDiv.contentEditable = true;
    contentDiv.contentEditable = true;

    // Add styling for the editable state
    li.classList.add('editable-note');

    // Create Save and Cancel buttons
    const saveBtn = document.createElement('button');
    saveBtn.textContent = 'Save';
    saveBtn.classList.add('save-btn');
    
    const cancelBtn = document.createElement('button');
    cancelBtn.textContent = 'Cancel';
    cancelBtn.classList.add('cancel-btn');
    
    // Append new buttons
    const buttonGroup = document.createElement('div');
    buttonGroup.classList.add('button-group');
    buttonGroup.appendChild(saveBtn);
    buttonGroup.appendChild(cancelBtn);
    li.appendChild(buttonGroup);

    // Append the entire list item to the container
    container.appendChild(li);

    // Event listeners for Save and Cancel
    saveBtn.addEventListener('click', () => {
        const updatedNote = {
            title: titleDiv.textContent,
            content: contentDiv.textContent,
        };
        secureFetch(`/api/notes/${note.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updatedNote),
        });

        window.location.href = '/app/notes';
    });

    cancelBtn.addEventListener('click', () => {
        window.location.href = '/app/notes';
    });
}

async function logout() {
// Clear the localStorage
    jwtToken = "";
    userID = "";
    username = "";
    userEmail = "";

    const response = await secureFetch(`${apiUrl}/logout`, {
        method: "DELETE"
    })

    localStorage.clear()
    window.location.href = `${appUrl}/`;
}

async function deleteNote(noteID) {
    console.log(`deleting note with id: ${noteID}`);
    jwtToken = localStorage.getItem("jwtToken")
    const response = await secureFetch(`${apiUrl}/notes/${noteID}`, {
        method: "DELETE"
    });

    loadNotes();
}

function editNote(noteID) {
    console.log(`editing note with id: ${noteID}`);
    localStorage.setItem("note-id", noteID)
    window.location.href = `${appUrl}/note`;
}

async function attemptRefresh() {
    const refreshResponse = await fetch(`${apiUrl}/refresh`, {
        method: "POST"
    });

    if (refreshResponse.ok) {
        const data = await refreshResponse.json();
        localStorage.setItem("jwtToken", data.token);
        return true;
    } else {
        return false; 
    }
}

// Global wrapper for making authenticated requests
async function secureFetch(url, options = {}) {
    jwtToken = localStorage.getItem("jwtToken");

    // Add current token to headers
    options.headers = {
        ...options.headers,
        "Authorization": `Bearer ${jwtToken}`
    };

    let response = await fetch(url, options);

    // If the token is expired (401)
    if (response.status === 401) {
        console.log("Access Token expired, attempting refresh...");
        const refreshed = await attemptRefresh();

        if (refreshed) {
            // Get the new token and retry the original request
            jwtToken = localStorage.getItem("jwtToken");
            options.headers["Authorization"] = `Bearer ${jwtToken}`;
            
            // SECOND ATTEMPT
            response = await fetch(url, options); 
        }
    }

    // Handle the final response (either original success, retried success, or final failure)
    if (!response.ok) {
        const errorData = await response.json();
        console.error(errorData.Error);
        handleErrorCode(response.status)
        return;
    }

    return response;
}

function handleErrorCode(errorCode) {
    switch (errorCode) {
        case 400:
            window.location.href = `${appUrl}/400`;
            break;
        case 401:
            window.location.href = `${appUrl}/401`;
            break;
        case 403:
            window.location.href = `${appUrl}/403`;
            break;
        case 404:
            window.location.href = `${appUrl}/404`;
            break;
        case 500:
            window.location.href = `${appUrl}/500`;
            break;
    }
}
