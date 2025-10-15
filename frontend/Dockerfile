# Utilisation d'une image Node.js officielle, version 20-alpine pour sa légèreté
FROM node:20-alpine

# Définition du répertoire de travail dans le conteneur
WORKDIR /app

# Copie des fichiers de gestion de dépendances
COPY package.json package-lock.json ./

# Installation des dépendances. Cette étape sera mise en cache si les fichiers ci-dessus ne changent pas.
RUN npm install

# Copie du reste du code source de l'application frontend
COPY . .

# Exposition du port 3000, le port par défaut pour le serveur de développement Next.js
EXPOSE 3000

# Commande pour démarrer le serveur de développement avec rechargement à chaud
# Idéal pour l'environnement de développement local.
CMD ["npm", "run", "dev"]
