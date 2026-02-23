import { Outlet } from "react-router-dom"

export default function Layout() {
  return (
    <div className="h-screen flex flex-col">

      {/* TOP BAR */}
      <div className="h-14 flex items-center justify-between px-6
        bg-glass/70 backdrop-blur-md border-b border-neon/30
        shadow-[0_0_20px_rgba(0,234,255,0.2)]">

        <div className="flex gap-6">
          <button className="hover:text-neon transition">Лента</button>
          <button className="hover:text-neon transition">Поиск</button>
        </div>

        <div className="font-bold text-neon tracking-widest">
          CHATS
        </div>

        <div>
          <div className="w-8 h-8 rounded-full bg-neon shadow-[0_0_10px_#00eaff]" />
        </div>
      </div>


      <div className="flex flex-1 overflow-hidden">

        {/* LEFT HUB COLUMN */}
        <div className="w-20 bg-glass/60 backdrop-blur-lg
          border-r border-neon/20
          flex flex-col items-center py-4 gap-4">

          {[1,2,3,4].map(i => (
            <div
              key={i}
              className="w-12 h-12 rounded-full
              bg-gradient-to-br from-purple-800 to-purple-900
              border border-neon/40
              shadow-[0_0_10px_rgba(0,234,255,0.4)]
              hover:scale-110 transition cursor-pointer"
            />
          ))}
        </div>


        {/* CENTER CHAT AREA */}
        <div className="flex-1 flex flex-col bg-glass/40 backdrop-blur-xl">

          <div className="flex-1 overflow-y-auto p-6">
            <Outlet />
          </div>

        </div>

      </div>


      {/* BOTTOM BAR */}
      <div className="h-16 flex items-center justify-center gap-12
        bg-glass/70 backdrop-blur-md border-t border-neon/30
        shadow-[0_0_20px_rgba(0,234,255,0.2)]">

        <button className="hover:text-neon transition">Друзья</button>
        <button className="hover:text-neon transition">Создать</button>
        <button className="hover:text-neon transition">Настройки</button>
      </div>

    </div>
  )
}