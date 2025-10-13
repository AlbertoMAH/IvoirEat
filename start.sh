#!/usr/bin/env sh

# S'assurer que le binaire Go est exécutable
chmod +x ./app

# Démarrer le backend Go en arrière-plan sur le port 8080
./app &

# Démarrer le serveur de production Next.js sur le port fourni par Render
cd frontend
npm start
