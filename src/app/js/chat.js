const CHAT_ID = "1"; // временно захардкодим

async function loadMessages() {
    const data = await apiRequest(`/api/v1/protected/chats/${CHAT_ID}/messages`);
    const container = document.getElementById("messages");
    container.innerHTML = "";

    if (data && data.messages) {
        data.messages.forEach(msg => {
            container.innerHTML += `
                <div class="message">
                    <b>${msg.sender_id}</b>: ${msg.content}
                </div>
            `;
        });
    }
}

async function sendMessage() {
    const input = document.getElementById("messageInput");
    const content = input.value;

    if (!content) return;

    await apiRequest(
        `/api/v1/protected/chats/${CHAT_ID}/messages`,
        "POST",
        { content }
    );

    input.value = "";
    loadMessages();
}

document.addEventListener("DOMContentLoaded", () => {
    loadMessages();
    setInterval(loadMessages, 3000);
});