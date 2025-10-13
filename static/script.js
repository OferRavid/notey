const baseUrl = "http://localhost:8080/";
const apiUrl = `${baseUrl}api`;
const appUrl = `${baseUrl}app`;
let jwtToken = "";
let refreshToken = "";
let userID = "";
let userEmail = "";


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
        login(email, password)
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

    login(email, password);
});

async function login(email, password) {
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
        refreshToken = data.refresh_token; // Store the refresh token after login
        userID = data.id; // Store the user's ID after login
        userEmail = data.email; // Store the user's email after login
        localStorage.setItem("jwtToken", jwtToken); // Save all of the above to localStorage
        localStorage.setItem("refreshToken", refreshToken);
        localStorage.setItem("userID", userID);
        localStorage.setItem("userEmail", userEmail);
        window.location.href = `${appUrl}/notes/notes.html`;
    } else {
        const errorText = await loginResponse.json();
        alert("Error: " + errorText.Error);
        console.error("Login error:", errorText.Error);
    }
}