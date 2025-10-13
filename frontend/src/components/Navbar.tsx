"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "../contexts/AuthContext";

export default function Navbar() {
  const { isAuthenticated, logout } = useAuth();
  const pathname = usePathname();

  if (!isAuthenticated) {
    return null; // Don't render Navbar on public pages
  }

  const navLinks = [
    { href: "/receipts/list", label: "Mes Reçus" },
    { href: "/receipts/upload", label: "Uploader" },
  ];

  return (
    <nav className="bg-white shadow-md">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center py-4">
          <Link href="/" className="text-xl font-bold text-indigo-600">
            IvoirEat
          </Link>
          <div className="hidden md:flex items-center space-x-8">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`text-gray-600 hover:text-indigo-600 transition-colors ${
                  pathname === link.href ? "font-semibold text-indigo-600" : ""
                }`}
              >
                {link.label}
              </Link>
            ))}
          </div>
          <div className="flex items-center">
            <button
              onClick={logout}
              className="px-4 py-2 font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
            >
              Déconnexion
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}
