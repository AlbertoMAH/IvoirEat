#!/bin/sh

# --- DEBUT SECTION DE DEBOGAGE ---
echo "--- [DEBUG] Début du script start.sh ---"
echo "--- [DEBUG] Répertoire de travail actuel (au début) :"
pwd
echo "--- [DEBUG] Contenu du répertoire racine :"
ls -la /app

echo "--- Changement de répertoire vers le backend ---"
cd backend

echo "--- [DEBUG] Répertoire de travail actuel (après cd) :"
pwd
echo "--- [DEBUG] Contenu du répertoire backend :"
ls -la
echo "--- [DEBUG] Contenu du répertoire public (si il existe) :"
ls -la public
echo "--- [DEBUG] Fin de la section de débogage ---"


echo "--- Démarrage du serveur Go ---"
./app
