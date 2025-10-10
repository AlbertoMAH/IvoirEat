#!/bin/sh

# Ce script orchestre le démarrage du frontend et du backend pour le déploiement.

echo "--- Démarrage du serveur Next.js (sur le port 3000) en arrière-plan ---"
yarn --cwd ./frontend start &

echo "--- Attente intelligente : vérification que le port 3000 est prêt ---"
WAIT_TIMEOUT=30
WAIT_COUNTER=0
# Boucle jusqu'à ce que la commande 'nc' (netcat) réussisse à se connecter au port 3000.
# L'option '-z' demande à nc de ne pas envoyer de données, juste de vérifier si le port est ouvert.
while ! nc -z localhost 3000; do
  # Si le compteur dépasse le timeout, on abandonne.
  if [ ${WAIT_COUNTER} -ge ${WAIT_TIMEOUT} ]; then
    echo "Erreur : Le serveur Next.js n'a pas répondu sur le port 3000 après ${WAIT_TIMEOUT} secondes."
    exit 1
  fi
  # Incrémente le compteur et attend 1 seconde avant de réessayer.
  WAIT_COUNTER=$((WAIT_COUNTER + 1))
  sleep 1
done

echo "--- Le serveur Next.js est prêt ! Démarrage du reverse proxy Go ---"
./backend/app
