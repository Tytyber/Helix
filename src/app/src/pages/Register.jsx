// src/pages/Register.jsx
import React, { useState } from "react";

export default function Register() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleRegister = async () => {
    const res = await fetch("/api/v1/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });
    const data = await res.json();
    console.log(data);
  };

  return (
    <div class="registerPage" style={{ padding: "2rem" }}>
      <h1 class="registerTextRegister">Register now!</h1>
      <input class="userRegisterPage" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} /><br/>
      <input class="emailRegisterPage" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} /><br/>
      <input class="passwordRegisterPage" placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} /><br/>
      <button href="Login.jsx" class="fromRegisterToLoginButton">Have an account</button>
      <button class="registerButton" onClick={handleRegister}>Register</button>
    </div>
  )
}