#!/bin/sh

# Ce script orchestre le démarrage du frontend et du backend pour le déploiement.

echo "--- Démarrage du serveur Next.js (sur le port 3000) en arrière-plan ---"
yarn --cwd ./frontend start &

echo "--- Démarrage du reverse proxy Go (sur le port 8080) en avant-plan ---"
./backend/app
