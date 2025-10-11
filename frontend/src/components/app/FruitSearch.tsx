"use client";

import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { Input } from "@/components/ui/input";

// Définir le type pour un fruit, correspondant à la structure Go
type Fruit = {
  id: number;
  name: string;
};

// Fonction pour appeler notre API Go
const fetchFruits = async (searchTerm: string): Promise<Fruit[]> => {
  if (!searchTerm.trim()) {
    return [];
  }
  // Utiliser la variable d'environnement pour l'URL de l'API
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const response = await fetch(`${apiUrl}/api/fruits?search=${encodeURIComponent(searchTerm)}`);
  if (!response.ok) {
    throw new Error("La réponse du réseau n'était pas bonne");
  }
  // Attendre 500ms pour simuler un délai réseau et voir l'état de chargement
  await new Promise(resolve => setTimeout(resolve, 500));
  return response.json();
};

const FruitSearch = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState("");

  // Mettre à jour le terme de recherche décalé après 300ms d'inactivité
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearchTerm(searchTerm);
    }, 300);

    return () => {
      clearTimeout(timer);
    };
  }, [searchTerm]);

  // Utiliser React Query pour récupérer les données
  const { data: fruits, isLoading, error } = useQuery<Fruit[], Error>({
    queryKey: ["fruits", debouncedSearchTerm],
    queryFn: () => fetchFruits(debouncedSearchTerm),
    enabled: !!debouncedSearchTerm, // N'exécuter la requête que si le terme n'est pas vide
  });

  return (
    <div className="w-full max-w-md">
      <h1 className="text-3xl font-bold mb-6 text-center">Rechercher un Fruit</h1>
      <Input
        type="text"
        placeholder="Tapez pour rechercher (ex: Pomme)..."
        className="text-lg p-4"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <div className="mt-8 min-h-[10rem] rounded-lg bg-white p-4 shadow-sm border border-gray-200">
        {isLoading && <p className="text-center text-gray-500 animate-pulse">Chargement...</p>}
        {error && <p className="text-center text-red-500 font-semibold">Erreur: {error.message}</p>}
        {!isLoading && !error && fruits && fruits.length > 0 && (
          <ul className="space-y-2">
            {fruits.map((fruit) => (
              <li key={fruit.id} className="p-2 rounded-md bg-gray-50 text-gray-800">
                {fruit.name}
              </li>
            ))}
          </ul>
        )}
        {!isLoading && !error && fruits && fruits.length === 0 && debouncedSearchTerm && (
          <p className="text-center text-gray-500 pt-8">Aucun fruit trouvé pour "{debouncedSearchTerm}".</p>
        )}
         {!isLoading && !error && (!fruits || fruits.length === 0) && !debouncedSearchTerm && (
          <p className="text-center text-gray-500 pt-8">Veuillez taper pour commencer la recherche.</p>
        )}
      </div>
    </div>
  );
};

export default FruitSearch;
