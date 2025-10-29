const apiUrl = `../api`;
const appUrl = `/app`;
let jwtToken = "";
let userID = "";
let username = "";
let userEmail = "";

document.getElementById("login-button").addEventListener("click", () => {
    window.location.href = `${appUrl}/login`;
});