import API from "./index";

export const register = (username: string, email: string, password: string) =>
  API.post("/auth/register", { username, email, password });

export const login = (email: string, password: string) =>
  API.post("/auth/login", { email, password });

export const getMe = () => API.get("/protected/me");

export const updateProfile = (data: any) =>
  API.put("/protected/profile", data);
