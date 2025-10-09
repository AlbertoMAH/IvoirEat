import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    // Définit l'URL du backend en fonction de l'environnement
    const backendUrl =
      process.env.NODE_ENV === "production"
        ? "https://ivoireat.onrender.com" // URL du backend en production
        : "http://localhost:8080"; // URL du backend en local

    return [
      {
        // Cette règle intercepte toutes les requêtes commençant par /api/
        source: "/api/:path*",
        // Et les redirige vers le backend approprié
        destination: `${backendUrl}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
