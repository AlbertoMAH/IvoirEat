"use client";

import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Search } from "lucide-react";


// Définir le type pour un fruit, correspondant à la structure Go
type Fruit = {
  id: number;
  name: string;
  description: string;
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
    <Card className="w-full max-w-lg mx-auto">
      <CardHeader>
        <CardTitle className="text-2xl font-bold tracking-tight text-center">
          Rechercher un Fruit
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
          <Input
            type="text"
            placeholder="Ex: Banane"
            className="pl-10 text-base"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="mt-6 min-h-[200px] space-y-3">
          {isLoading && (
            <>
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-16 w-full" />
            </>
          )}
          {error && (
            <div className="text-center text-red-500 font-semibold pt-8">
              <p>Erreur: {error.message}</p>
              <p className="text-sm text-gray-500">Le serveur backend est-il bien lancé ?</p>
            </div>
          )}
          {!isLoading && !error && fruits && fruits.length > 0 && (
            fruits.map((fruit) => (
              <div key={fruit.id} className="p-3 rounded-lg border bg-gray-50/50">
                <p className="font-semibold text-gray-800">{fruit.name}</p>
                <p className="text-sm text-gray-600">{fruit.description}</p>
              </div>
            ))
          )}
          {!isLoading && !error && fruits && fruits.length === 0 && debouncedSearchTerm && (
            <p className="text-center text-gray-500 pt-8">Aucun fruit trouvé pour "{debouncedSearchTerm}".</p>
          )}
          {!isLoading && !error && (!fruits || fruits.length === 0) && !debouncedSearchTerm && (
            <p className="text-center text-gray-500 pt-8">Veuillez taper pour commencer la recherche.</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default FruitSearch;
