#!/bin/bash

# Chemins
SOURCE="/storage/emulated/0/ALBERT/Backend/mangoedit.txt"
DEST="$HOME/GoBackend/main.go"
LOG="$HOME/GoBackend/server.log"

# Créer le dossier si besoin
mkdir -p "$(dirname "$DEST")"
touch "$LOG"

# Copier le contenu
cp "$SOURCE" "$DEST"
echo "✅ main.go mis à jour avec mangoedit.txt"

# Tuer le serveur Go en cours (si lancé)
pkill -f "go run $DEST" 2>/dev/null

# Relancer le serveur Go en arrière-plan
nohup go run "$DEST" > "$LOG" 2>&1 &
echo "🚀 Serveur Go relancé automatiquement"
