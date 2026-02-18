import API from "./index";

// Личные чаты
export const createChat = (userId: string) =>
  API.post("/protected/chats", { user_id: userId });

export const getChatMessages = (chatId: string) =>
  API.get(`/protected/chats/${chatId}/messages`);

export const sendChatMessage = (chatId: string, content: string) =>
  API.post(`/protected/chats/${chatId}/messages`, { content });

// Групповые чаты
export const createGroup = (name: string) =>
  API.post("/protected/groups", { name });

export const joinGroup = (groupId: string) =>
  API.post(`/protected/groups/${groupId}/join`);

export const leaveGroup = (groupId: string) =>
  API.post(`/protected/groups/${groupId}/leave`);

export const getGroupMessages = (groupId: string) =>
  API.get(`/protected/groups/${groupId}/messages`);

export const sendGroupMessage = (groupId: string, content: string) =>
  API.post(`/protected/groups/${groupId}/messages`, { content });
