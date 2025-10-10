#!/bin/sh

# Ce script orchestre le démarrage du frontend et du backend pour le déploiement.

echo "--- Démarrage du serveur Next.js (sur le port 3000) en arrière-plan ---"
# Redirige la sortie standard (stdout) et les erreurs (stderr) vers frontend.log
# pour que nous puissions déboguer les erreurs de démarrage.
yarn --cwd ./frontend start > frontend.log 2>&1 &

echo "--- Pause de 5 secondes pour laisser le temps au serveur Next.js de démarrer ---"
sleep 5

echo "--- Démarrage du reverse proxy Go (sur le port 8080) en avant-plan ---"
./backend/app
