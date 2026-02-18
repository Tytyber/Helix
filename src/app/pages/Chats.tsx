import { useEffect, useState } from "react";
import { getChatMessages, sendChatMessage } from "../api/chats";

interface Message {
  id: string;
  senderID: string;
  content: string;
  createdAt: string;
}

export default function Chat({ chatId }: { chatId: string }) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMsg, setNewMsg] = useState("");

  useEffect(() => {
    fetchMessages();
    const interval = setInterval(fetchMessages, 3000); // автообновление
    return () => clearInterval(interval);
  }, [chatId]);

  const fetchMessages = async () => {
    const res = await getChatMessages(chatId);
    setMessages(res.data.messages);
  };

  const sendMessage = async () => {
    if (!newMsg) return;
    await sendChatMessage(chatId, newMsg);
    setNewMsg("");
    fetchMessages();
  };

  return (
    <div className="p-4">
      <div className="h-[400px] overflow-y-auto border p-2 mb-2">
        {messages.map((m) => (
          <div key={m.id} className="mb-1">
            <strong>{m.senderID}:</strong> {m.content}
          </div>
        ))}
      </div>
      <div className="flex gap-2">
        <input
          value={newMsg}
          onChange={(e) => setNewMsg(e.target.value)}
          className="flex-1 border p-1 rounded"
        />
        <button onClick={sendMessage} className="bg-purple-500 text-white p-1 rounded">
          Send
        </button>
      </div>
    </div>
  );
}
