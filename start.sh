#!/bin/sh

# Ce script orchestre le démarrage du frontend et du backend pour le déploiement.

echo "--- Démarrage du serveur Next.js (sur le port 3000) en arrière-plan ---"
yarn --cwd ./frontend start > frontend.log 2>&1 &

echo "--- Pause de 8 secondes pour laisser le temps au serveur Next.js de démarrer et de logger ---"
sleep 8

echo "--- Affichage des logs du frontend (contenu de frontend.log) ---"
# Cette commande affichera le contenu du fichier log dans le flux de logs principal de Render.
cat frontend.log
echo "--- Fin des logs du frontend ---"

echo "--- Démarrage du reverse proxy Go (sur le port 8080) en avant-plan ---"
./backend/app
