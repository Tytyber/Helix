import { useState } from "react"
import { useNavigate } from "react-router-dom"
import api from "../api/axios"

export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const navigate = useNavigate()

  const handleLogin = async () => {
    try {
      const res = await api.post("/auth/login", { email, password })
      localStorage.setItem("token", res.data.token)
      navigate("/feed")
    } catch (err) {
      alert("Login failed")
    }
  }

  return (
<div class = "loginPage">
      <h2 class = "loginTextLogin">Login to account</h2>
      <input class="emailRegisterLogin" placeholder="Email" onChange={e => setEmail(e.target.value)} />
      <br />
      <input class="passwordRegisterLogin" type="password" placeholder="Password" onChange={e => setPassword(e.target.value)} />
      <br />
      <button class="loginButton" onClick={handleLogin}>Login now</button>
    </div>
  )
}