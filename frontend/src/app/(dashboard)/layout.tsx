import React from "react"
import Link from "next/link"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <section className="flex h-screen">
      <aside className="w-64 bg-gray-800 text-white p-4">
        {/* Sidebar content */}
        <h2 className="text-xl font-bold">Dashboard</h2>
        <nav className="mt-8">
          <ul>
            <li>
              <Link href="/" className="block py-2 px-4 rounded hover:bg-gray-700">
                Overview
              </Link>
            </li>
            <li>
              <Link href="/parkings" className="block py-2 px-4 rounded hover:bg-gray-700">
                Parkings
              </Link>
            </li>
            <li>
              <a href="#" className="block py-2 px-4 rounded hover:bg-gray-700">
                Admins
              </a>
            </li>
            <li>
              <a href="#" className="block py-2 px-4 rounded hover:bg-gray-700">
                Reservations
              </a>
            </li>
          </ul>
        </nav>
      </aside>
      <main className="flex-1 flex flex-col">
        <header className="bg-white shadow-md p-4">
          {/* Topbar content */}
          <h1 className="text-xl font-semibold">Super Admin Dashboard</h1>
        </header>
        <div className="flex-1 p-8 overflow-y-auto">{children}</div>
      </main>
    </section>
  )
}
