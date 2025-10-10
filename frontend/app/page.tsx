'use client';

import { useState, FormEvent } from 'react';

// Définition du type pour la réponse de l'API Fruit
interface Fruit {
  name: string;
  description: string;
}

export default function Home() {
  const [searchTerm, setSearchTerm] = useState('');
  const [fruit, setFruit] = useState<Fruit | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSearch = async (e: FormEvent) => {
    e.preventDefault();
    if (!searchTerm) {
      setError("Veuillez entrer un nom de fruit.");
      return;
    }

    setIsLoading(true);
    setFruit(null);
    setError(null);

    try {
      // Appel à notre API Go backend
      const res = await fetch(`/api/v1/fruits/${encodeURIComponent(searchTerm)}`);

      if (res.status === 404) {
        setError(`Le fruit "${searchTerm}" n'a pas été trouvé.`);
        setFruit(null);
        return;
      }

      if (!res.ok) {
        throw new Error(`Erreur du serveur: ${res.status}`);
      }

      const data: Fruit = await res.json();
      setFruit(data);

    } catch (err) {
      console.error("Erreur lors de la recherche du fruit:", err);
      const errorMessage = err instanceof Error ? err.message : 'Une erreur inconnue est survenue.';
      setError(`Impossible de contacter le backend. Détail : ${errorMessage}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-4 sm:p-8 md:p-24 bg-gray-100 dark:bg-gray-900">
      <div className="w-full max-w-md p-6 sm:p-8 bg-white dark:bg-gray-800 rounded-xl shadow-lg space-y-6">
        <div className="text-center">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-800 dark:text-white">
            Recherche de Fruits
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-2">
            Entrez le nom d&apos;un fruit pour trouver sa description.
          </p>
        </div>

        <form onSubmit={handleSearch} className="flex flex-col sm:flex-row gap-3">
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Ex: Pomme, Banane..."
            className="flex-grow px-4 py-2 text-gray-700 bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600"
          />
          <button
            type="submit"
            disabled={isLoading}
            className="px-6 py-2 bg-blue-600 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors duration-200"
          >
            {isLoading ? 'Recherche...' : 'Rechercher'}
          </button>
        </form>

        <div className="mt-6 min-h-[120px] flex items-center justify-center">
          {error && (
            <div className="text-center text-red-600 bg-red-100 dark:bg-red-900/20 dark:text-red-400 p-4 rounded-lg">
              <p><strong>Erreur :</strong> {error}</p>
            </div>
          )}
          {fruit && (
            <div className="w-full p-5 bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 rounded-r-lg">
              <h2 className="text-xl font-bold text-blue-800 dark:text-blue-300">{fruit.name}</h2>
              <p className="mt-1 text-gray-600 dark:text-gray-300">{fruit.description}</p>
            </div>
          )}
        </div>
      </div>
       <footer className="text-center mt-8 text-gray-500 dark:text-gray-400 text-sm">
        <p>Test d&apos;intégration full-stack</p>
      </footer>
    </main>
  );
}
