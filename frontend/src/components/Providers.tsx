"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode, useState } from "react";

// Ce composant "Providers" enveloppe notre application avec le QueryClientProvider
// pour que nous puissions utiliser React Query partout.
export default function Providers({ children }: { children: ReactNode }) {
  // Nous utilisons useState pour nous assurer que le QueryClient n'est créé qu'une seule fois
  // côté client, pour éviter les fuites de données entre les utilisateurs.
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
