async function login() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const data = await apiRequest("/api/v1/auth/login", "POST", {
        email,
        password
    });

    if (data && data.access_token) {
        localStorage.setItem("token", data.access_token);
        window.location.href = "index.html";
    } else {
        document.getElementById("error").innerText = "Ошибка входа";
    }
}

async function getMe() {
    const data = await apiRequest("/api/v1/protected/me");
    return data;
}

function logout() {
    localStorage.removeItem("token");
    window.location.href = "login.html";
}