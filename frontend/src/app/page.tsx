"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "../contexts/AuthContext";

export default function Home() {
  const { isAuthenticated, logout, token } = useAuth();
  const [protectedMessage, setProtectedMessage] = useState("");
  const [error, setError] = useState("");

  const handleProtectedRouteClick = async () => {
    setError("");
    setProtectedMessage("");

    if (!token) {
      setError("You are not logged in.");
      return;
    }

    try {
      const response = await fetch("/api/validate", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || "Failed to fetch protected data");
      }

      const data = await response.json();
      setProtectedMessage(JSON.stringify(data, null, 2));
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  return (
    <main className="flex flex-col items-center justify-center min-h-screen p-8 bg-gray-50">
      <div className="w-full max-w-2xl p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h1 className="text-3xl font-bold text-center text-gray-900">
          IvoirEat - Gestion des Notes de Frais
        </h1>

        {isAuthenticated ? (
          <div className="space-y-4 text-center">
            <p className="text-lg">Bienvenue ! Vous êtes connecté.</p>
            <button
              onClick={logout}
              className="px-4 py-2 font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
            >
              Déconnexion
            </button>
          </div>
        ) : (
          <div className="flex justify-center space-x-4">
            <Link
              href="/login"
              className="px-4 py-2 font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
            >
              Se connecter
            </Link>
            <Link
              href="/register"
              className="px-4 py-2 font-medium text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
            >
              S&apos;inscrire
            </Link>
          </div>
        )}

        <hr />

        <div className="space-y-4">
          <h2 className="text-xl font-semibold">Test de Route Protégée</h2>
          <p className="text-sm text-gray-600">
            Cliquez sur ce bouton pour tester l&apos;accès à une route qui nécessite
            d&apos;être authentifié.
          </p>
          <button
            onClick={handleProtectedRouteClick}
            disabled={!isAuthenticated}
            className="w-full px-4 py-2 font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            Tester la route /api/validate
          </button>
          {protectedMessage && (
            <div className="p-4 mt-4 text-sm text-green-700 bg-green-100 border border-green-200 rounded-md">
              <h3 className="font-bold">Réponse du serveur :</h3>
              <pre className="mt-2 whitespace-pre-wrap">{protectedMessage}</pre>
            </div>
          )}
          {error && (
            <div className="p-4 mt-4 text-sm text-red-700 bg-red-100 border border-red-200 rounded-md">
              <h3 className="font-bold">Erreur :</h3>
              <p>{error}</p>
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
