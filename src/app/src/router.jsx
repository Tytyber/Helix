import { BrowserRouter, Routes, Route } from "react-router-dom"
import Login from "./pages/Login"
import Register from "./pages/Register"
import Feed from "./pages/Feed"
import Friends from "./pages/Friends"
import Calls from "./pages/Calls"
import Layout from "./components/Layout"

export default function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route element={<Layout />}>
          <Route path="/feed" element={<Feed />} />
          <Route path="/friends" element={<Friends />} />
          <Route path="/calls" element={<Calls />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}