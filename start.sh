#!/bin/sh

# Ce script orchestre le démarrage du frontend et du backend pour le déploiement.

echo "--- Démarrage du serveur Next.js (sur le port 3000) en arrière-plan ---"
# La commande 'yarn start' lance le serveur de production Next.js après un 'yarn build'.
# Le '&' à la fin exécute la commande en arrière-plan, ce qui permet au script de continuer.
yarn --cwd ./frontend start &

echo "--- Démarrage du reverse proxy Go (sur le port 8080) en avant-plan ---"
# L'exécutable 'app' a été compilé dans le dossier 'backend' par la commande de build de Render.
# C'est le processus principal. Si ce script se termine, le service Render s'arrête.
./backend/app
