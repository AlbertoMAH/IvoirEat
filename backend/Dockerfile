# Utilisation de l'image officielle Go 1.25, basée sur Alpine pour une taille réduite
FROM golang:1.25-alpine

# Définition du répertoire de travail dans le conteneur
WORKDIR /app

# Copie des fichiers de dépendances du module Go
COPY go.mod go.sum ./

# Téléchargement des dépendances. Celles-ci seront mises en cache si les fichiers go.mod/go.sum ne changent pas.
RUN go mod download

# Copie du reste du code source du backend dans le conteneur
COPY . .

# Exposition du port 8080 pour permettre la communication avec l'application
EXPOSE 8080

# Commande pour lancer l'application avec 'go run'.
# 'go run' compile et exécute l'application, ce qui est parfait pour le développement.
# Les modifications du code (montées via un volume dans docker-compose) seront prises en compte au redémarrage du conteneur.
CMD ["go", "run", "main.go"]
