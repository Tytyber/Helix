import { useEffect, useState } from "react"
import api from "../api/axios"

export default function Friends() {
  const [friends, setFriends] = useState([])

  useEffect(() => {
    const fetchFriends = async () => {
      const res = await api.get("/protected/friends")
      setFriends(res.data.friends)
    }
    fetchFriends()
  }, [])

  return (
    <div>
      <h2>Friends</h2>
      {friends.map(friend => (
        <div key={friend.id}>
          {friend.username} ({friend.email})
        </div>
      ))}
    </div>
  )
}