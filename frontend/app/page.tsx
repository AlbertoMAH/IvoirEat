'use client';

import { useState, useEffect } from 'react';

export default function Home() {
  const [message, setMessage] = useState('Chargement...');
  const [error, setError] = useState('');

  useEffect(() => {
    // On utilise un chemin relatif pour que l'appel fonctionne
    // aussi bien en local qu'une fois déployé sur Render.
    // Le reverse proxy Go se chargera de transmettre la requête.
    fetch('/api/message')
      .then((res) => {
        if (!res.ok) {
          throw new Error(`Erreur HTTP: ${res.status}`);
        }
        return res.json();
      })
      .then((data) => {
        setMessage(data.message);
      })
      .catch((err) => {
        console.error("Erreur lors de la récupération du message:", err);
        setError(`Impossible de contacter le backend. Détail : ${err.message}`);
        setMessage(''); // Clear loading message
      });
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gray-50">
      <div className="text-center p-10 border rounded-lg shadow-lg bg-white">
        <h1 className="text-3xl font-bold mb-4 text-gray-800">
          Test de Connexion Frontend-Backend
        </h1>
        {error ? (
          <p className="text-lg text-red-600 bg-red-100 p-4 rounded">
            <strong>Erreur :</strong> {error}
          </p>
        ) : (
          <p className="text-lg text-gray-600">
            Message du backend :
            <span className="font-semibold text-blue-600 ml-2 p-2 bg-blue-100 rounded">
              {message}
            </span>
          </p>
        )}
      </div>
    </main>
  );
}
