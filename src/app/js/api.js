
const API_URL = "http://localhost:4000/api/v1"; // поменяй под свой API

// ===== Регистрация =====
const registerForm = document.getElementById("registerForm");
if (registerForm) {
    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(registerForm);
        const data = {
            username: formData.get("username"),
            email: formData.get("email"),
            password: formData.get("password")
        };

        const res = await fetch(`${API_URL}/auth/register`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(data)
        });

        const json = await res.json();
        if (res.ok) {
            localStorage.setItem("access_token", json.access_token);
            window.location.href = "chat.html";
        } else {
            alert(json.error || "Ошибка регистрации");
        }
    });
}

// ===== Вход =====
const loginForm = document.getElementById("loginForm");
if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(loginForm);
        const data = {
            email: formData.get("email"),
            password: formData.get("password")
        };

        const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(data)
        });

        const json = await res.json();
        if (res.ok) {
            localStorage.setItem("access_token", json.access_token);
            window.location.href = "chat.html";
        } else {
            alert(json.error || "Ошибка входа");
        }
    });
}

// ===== Чат =====
const chatForm = document.getElementById("chatForm");
const chatWindow = document.getElementById("chatWindow");
if (chatForm && chatWindow) {
    const token = localStorage.getItem("access_token");
    if (!token) window.location.href = "login.html";

    async function loadMessages() {
        const res = await fetch(`${API_URL}/protected/chats/1/messages`, {
            headers: {"Authorization": `Bearer ${token}`}
        });
        if (res.ok) {
            const json = await res.json();
            chatWindow.innerHTML = json.data.map(m => `<p><b>${m.username}:</b> ${m.text}</p>`).join("");
            chatWindow.scrollTop = chatWindow.scrollHeight;
        }
    }

    chatForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const msgInput = document.getElementById("message");
        const text = msgInput.value.trim();
        if (!text) return;
        const res = await fetch(`${API_URL}/protected/chats/1/messages`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({text})
        });
        if (res.ok) {
            msgInput.value = "";
            loadMessages();
        } else {
            alert("Ошибка отправки сообщения");
        }
    });

    setInterval(loadMessages, 3000);
    loadMessages();
}

// ===== Выход =====
const logoutBtn = document.getElementById("logout");
if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
        localStorage.removeItem("access_token");
        window.location.href = "login.html";
    });
}