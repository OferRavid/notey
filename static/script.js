const apiUrl = "http://localhost:8080/api"; // Change this if your server runs on a different URL
const authSection = ["registrationSection", "formSection"];
let jwtToken = "";
let refreshToken = "";
let userID = "";
let userEmail = "";
let noteCount = 0;

// Function to show/hide sections
function showSection(sectionsID) {
    const sections = ["registrationSection", "formSection", "notesSection", "metricsSection"];
    sections.forEach(section => {
        document.getElementById(section).classList.add("hidden");
    });
    for (let i = 0; i < sectionsID.length; i++) {
        document.getElementById(sectionsID[i]).classList.remove("hidden");
    }
}

// Check for existing token on page load
window.onload = () => {
    jwtToken = localStorage.getItem("jwtToken"); // Retrieve token from localStorage

    if (jwtToken) {
        showNotes(); // If token exists, show notes section
    } else {
        showSection(authSection); // Show both registration and login sections
    }
};

// Handle registration form submission
document.getElementById("registrationForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const email = document.getElementById("regEmail").value;
    const password = document.getElementById("regPassword").value;

    const regResponse = await fetch(`${apiUrl}/users`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email, password })
    });

    if (regResponse.ok) {
        const loginResponse = await fetch(`${apiUrl}/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password })
        });

        if (loginResponse.ok) {
            const data = await loginResponse.json();
            jwtToken = data.token; // Store the token after login
            refreshToken = data.refreshToken; // Store the refresh token after login
            userID = data.userID; // Store the user's ID after login
            userEmail = data.userEmail; // Store the user's email after login
            localStorage.setItem("jwtToken", jwtToken); // Save the token to localStorage
            localStorage.setItem("refreshToken", refreshToken);
            localStorage.setItem("userID", userID);
            localStorage.setItem("userEmail", userEmail);
            localStorage.setItem("noteCount", String(noteCount)); // Save a note counter to localStorage
            showNotes(); // Go to notes section after successful login
        } else {
            const errorText = await loginResponse.json();
            alert("Error: " + errorText.Error);
            console.error("Login error:", errorText.Error);
        }
    } else {
        const errorText = await regResponse.json();
        alert("Error: " + errorText.Error);
        console.error("Registration error:", errorText.Error);
    }
});

// Handle login form submission
document.getElementById("userForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const loginResponse = await fetch(`${apiUrl}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email, password })
    });

    if (loginResponse.ok) {
        const data = await loginResponse.json();
        jwtToken = data.token; // Store the token after login
        refreshToken = data.refreshToken; // Store the refresh token after login
        userID = data.userID; // Store the user's ID after login
        userEmail = data.userEmail; // Store the user's email after login
        noteCount = getNoteCount(userID)
        localStorage.setItem("jwtToken", jwtToken); // Save the token to localStorage
        localStorage.setItem("refreshToken", refreshToken);
        localStorage.setItem("userID", userID);
        localStorage.setItem("userEmail", userEmail);
        localStorage.setItem("noteCount", String(noteCount));
        showNotes(); // Go to notes section after successful login
    } else {
        const errorText = await loginResponse.json();
        alert("Error: " + errorText.Error);
        console.error("Login error:", errorText.Error);
    }
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
        noteCount = parseInt(localStorage.getItem("noteCount"))
        noteCount += 1
        localStorage.setItem("noteCount", String(noteCount))
        loadNotes(); // Reload notes after adding
        document.getElementById("noteForm").reset(); // Reset form
    } else {
        const errorText = await response.json()
        alert("Error: " + errorText.Error);
        console.error("Note error:", errorText.Error);
    }
});

async function getNoteCount(userID) {
    const response = await fetch(`${apiUrl}/notes/?author_id=${userID}`, {
        headers: {
            "Authorization": `Bearer ${jwtToken}`
        }
    });

    let notes = []; // Initialize notes as an empty array

    // Check if the response is okay
    if (response.ok) {
        // Try to parse the response as JSON
        notes = await response.json();
        return notes.length;
    }
    return 0;
}

// Show notes section
function showNotes() {
    showSection(["notesSection"]);
    document.getElementById("logoutBtn").classList.remove("hidden");
    document.getElementById("notesH2").textContent = `${localStorage.getItem("userEmail")}'s notes`
    loadNotes();
}

// Load notes from the API
async function loadNotes() {
    if (parseInt(localStorage.getItem("noteCount")) === 0) { return; }
    const response = await fetch(`${apiUrl}/notes/?author_id=${localStorage.getItem("userID")}`, {
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
            const li = document.createElement("li");
            li.innerHTML = `<strong>${note.title}</strong><p>${note.content}</p>`;
            notesList.appendChild(li);
        });
    }
}

// Handle logout
document.getElementById("logoutBtn").addEventListener("click", () => {
    jwtToken = ""; // Clear the token
    refreshToken = "";
    userID = "";
    noteCount = 0;
    localStorage.setItem("jwtToken", jwtToken);
    localStorage.setItem("refreshToken", refreshToken);
    localStorage.setItem("userID", userID);
    localStorage.setItem("noteCount", String(noteCount));
    document.getElementById("logoutBtn").classList.add("hidden");
    showSection(authSection); // Show both registration and login sections
});
