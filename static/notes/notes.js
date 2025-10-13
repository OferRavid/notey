const baseUrl = "http://localhost:8080/";
const apiUrl = `${baseUrl}api`;
const appUrl = `${baseUrl}app`;
let jwtToken = "";
let refreshToken = "";
let userID = "";
let userEmail = "";


// Check for existing token on page load
window.onload = () => {
    jwtToken = localStorage.getItem("jwtToken"); // Retrieve token from localStorage

    if (jwtToken) {
        loadNotes(); // If token exists, show notes section
    } else {
        window.location.href = `${appUrl}/`;
    }
};

// Handle logout
document.getElementById("logoutBtn").addEventListener("click", () => {
    // Clear the localStorage
    jwtToken = "";
    refreshToken = "";
    userID = "";
    userEmail = "";
    localStorage.setItem("jwtToken", jwtToken);
    localStorage.setItem("refreshToken", refreshToken);
    localStorage.setItem("userID", userID);
    localStorage.setItem("userEmail", userEmail);
    window.location.href = `${appUrl}/`;
});

// Handle note form submission
document.getElementById("noteForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const title = document.getElementById("noteTitle").value;
    const content = document.getElementById("noteContent").value;

    const response = await fetch(`${apiUrl}/notes`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${localStorage.getItem("jwtToken")}`
        },
        body: JSON.stringify({ title, content })
    });

    if (response.ok) {
        loadNotes(); // Reload notes after adding
        document.getElementById("noteForm").reset(); // Reset form
    } else {
        const errorText = await response.json()
        alert("Error: " + errorText.Error);
        console.error("Note error:", errorText.Error);
    }
});


// Load notes from the API
async function loadNotes() {
    const response = await fetch(`${apiUrl}/notes`, {
        headers: {
            "Authorization": `Bearer ${jwtToken}`
        }
    });

    let notes = []; // Initialize notes as an empty array

    // Check if the response is okay
    if (response.ok) {
        // Try to parse the response as JSON
        notes = await response.json();
    } else {
        // If the response is not okay, handle the error
        const errorText = await response.json();
        console.error("Error loading notes:", errorText.Error);
        alert("Error loading notes: " + errorText.Error);
        return; // Exit the function
    }

    const notesList = document.getElementById("notesList");
    notesList.innerHTML = ""; // Clear the existing list

    if (notes.length === 0) {
        // If there are no notes, display a message
        const noNotesMessage = document.createElement("li");
        noNotesMessage.textContent = "No notes available.";
        notesList.appendChild(noNotesMessage);
    } else {
        // If there are notes, display them
        notes.forEach(note => {
            // Creating elements to design the note in the document
            const li = document.createElement("li");
            const editButton = document.createElement("edit-button");
            const editText = document.createElement("edit-button-text");
            const deleteButton = document.createElement("delete-button");
            const deleteText = document.createElement("del-button-text");

            // Adding the note's content to the new list item
            li.innerHTML = `<strong style='font-size: 20px;'><u>${note.title}</u></strong><p style='font-size: 16px;'>${note.content}</p>`;

            // Designing the edit button
            editButton.classList.add("edit-button");
            editText.classList.add("edit-button-text");
            editText.textContent = "Edit";
            editButton.appendChild(editText);
            editButton.setAttribute("noteID", note.id);
            editButton.addEventListener('click', (event) => {
                const clickedButton = event.target;
                const noteID = clickedButton.getAttribute("noteID")
                editNote(noteID);
            });

            // Designing the delete button
            deleteButton.classList.add("delete-button");
            deleteText.classList.add("del-button-text");
            deleteText.textContent = "Delete";
            deleteButton.appendChild(deleteText);
            deleteButton.setAttribute("noteID", note.id);
            deleteButton.addEventListener('click', (event) => {
                const clickedButton = event.target;
                const noteID = clickedButton.getAttribute("noteID")
                deleteNote(noteID);
            });

            // Adding all the elements to the list
            li.appendChild(editButton);
            li.appendChild(deleteButton);
            notesList.appendChild(li);
        });
    }
}

function deleteNote(noteID) {
    console.log(`deleting note with id: ${noteID}`);
    loadNotes();
}

function editNote(noteID) {
    console.log(`editting note with id: ${noteID}`);
    loadNotes();
}
