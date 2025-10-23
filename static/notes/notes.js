const apiUrl = `/api`;
const appUrl = `/app`;
let jwtToken = "";
let refreshToken = "";
let userID = "";
let userEmail = "";


// Check for existing token on page load
window.onload = () => {
    jwtToken = localStorage.getItem("jwtToken"); // Retrieve token from localStorage

    if (jwtToken) {
        loadNotes(); // If token exists, show notes
    } else {
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
    // checkRefreshToken();
    const title = document.getElementById("note-title").value;
    const content = document.getElementById("note-content").value;

    const response = await fetch(`${apiUrl}/notes`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${localStorage.getItem("jwtToken")}`
        },
        body: JSON.stringify({ title, content })
    });

    if (response.ok) {
        document.getElementById("note-form").reset(); // Reset form
        loadNotes(); // Reload notes after adding
    } else {
        const errorText = await response.json();
        console.error("Note error:", errorText.Error);
        handleErrorCode(errorText.status);
    }
});


// Load notes from the API
async function loadNotes() {
    // checkRefreshToken();
    const response = await fetch(`${apiUrl}/notes`, {
        headers: {
            "Authorization": `Bearer ${jwtToken}`
        }
    });

    if (!response.ok) {
        // If the response is not okay, handle the error
        const errorText = await response.json();
        console.error("Error loading notes:", errorText.Error);
        handleErrorCode(errorText.status);
    }

    let notes = await response.json();
    const notesList = document.getElementById("notes-list");
    notesList.innerHTML = ""; // Clear the existing list

    if (notes.length === 0) {
        // If there are no notes, display a message
        const noNotesMessage = document.createElement("li");
        noNotesMessage.textContent = "No notes available.";
        notesList.appendChild(noNotesMessage);
    } else {
        // If there are notes, display them
        notes.forEach(note => {
            
            const li = document.createElement("li");
            li.innerHTML = `<strong style='font-size: 20px;'><u>${note.title}</u></strong><p style='font-size: 16px;'>${note.content}</p>`;
            
            const editButton = getButtonElement("edit-button", "Edit", note.id, editNote);
            const deleteButton = getButtonElement("delete-button", "Delete", note.id, deleteNote);

            // Adding all the elements to the list
            li.appendChild(editButton);
            li.appendChild(deleteButton);
            notesList.appendChild(li);
        });
    }
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

function logout() {
// Clear the localStorage
    jwtToken = "";
    refreshToken = "";
    userID = "";
    userEmail = "";
    localStorage.clear()
    window.location.href = `${appUrl}/`;
}

function deleteNote(noteID) {
    // checkRefreshToken();
    console.log(`deleting note with id: ${noteID}`);
    loadNotes();
}

function editNote(noteID) {
    // checkRefreshToken();
    localStorage.setItem("note-id", noteID)
    console.log(`editting note with id: ${noteID}`);
    loadNotes();
}

async function checkRefreshToken() {
    refreshToken = localStorage.getItem("refreshToken");
    if (refreshToken === "") {
        handleErrorCode(401);
    }
    const response = await fetch(`${apiUrl}/refresh`, {
        headers: {
            "Authorization": `Bearer ${refreshToken}`
        }
    });

    if (response.ok) {
        const data = await response.json();
        jwtToken = data.token;
        localStorage.setItem("jwtToken", jwtToken);
    } else {
        const data = await response.json();
        console.error("Authorization error:", data.Error);
        
        if (response.status === 401 && !data.revoked) {
            const revokeResponse = await fetch(`${apiUrl}/revoke`, {
                headers: {
                    "Authorization": `Bearer ${refreshToken}`
                }
            });

            if (!revokeResponse.ok) {
                const errorResponse = await revokeResponse.json();
                console.error("Error while revoking token:", errorResponse.Error);
                handleErrorCode(errorResponse.status);
            }
        }
        handleErrorCode(response.status)
    }
}

function handleErrorCode(errorCode) {
    switch (errorCode) {
        case 400:
        case 401:
        case 403:
        case 404:
        case 500:
    }
}
