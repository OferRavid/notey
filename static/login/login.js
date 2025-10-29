const apiUrl = `/api`;
const appUrl = `/app`;
let jwtToken = "";
let userID = "";
let username = "";
let userEmail = "";

// Function to show/hide sections
function showSection(sectionID) {
    const sections = ["signup-section", "login-section"];
    sections.forEach(section => {
        if (sectionID === section) {
            document.getElementById(sectionID).classList.remove("hidden");
        } else {
            document.getElementById(section).classList.add("hidden");
        }
    });
}

window.onload = () => {
    // Make sure localStorage and cookie are empty
    localStorage.clear();
    document.cookie = "refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

    showSection("login-section"); // Show login sections
};

// Handle registration form submission
document.getElementById("signUpForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const username = document.getElementById("signupUsername").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("signupPassword").value;

    const regResponse = await fetch(`${apiUrl}/users`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, email, password })
    });

    if (regResponse.ok) {
        login(username, password)
    } else {
        const errorText = await regResponse.json();
        alert("Error: " + errorText.Error);
        console.error("Sign up error:", errorText.Error);
    }
});

// Handle login form submission
document.getElementById("loginForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const username = document.getElementById("loginUsername").value;
    const password = document.getElementById("loginPassword").value;

    login(username, password);
});

async function login(username, password) {
    const loginResponse = await fetch(`${apiUrl}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, password })
    });

    if (loginResponse.ok) {
        const data = await loginResponse.json();
        jwtToken = data.token; // Store the token after login
        userID = data.id; // Store the user's ID after login
        username = data.username; // Store the user's Username after login
        userEmail = data.email; // Store the user's email after login
        localStorage.setItem("jwtToken", jwtToken); // Save all of the above to localStorage
        localStorage.setItem("userID", userID);
        localStorage.setItem("username", username);
        localStorage.setItem("userEmail", userEmail);
        window.location.href = `${appUrl}/notes`;
    } else {
        const errorText = await loginResponse.json();
        alert("Error: " + errorText.Error);
        console.error("Login error:", errorText.Error);
    }
}

document.getElementById("signup-button").addEventListener("click", (event) => {
    showSection("signup-section");
});
