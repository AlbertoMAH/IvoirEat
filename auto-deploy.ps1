# =======================================================
# Script de déploiement local automatisé avec Docker
# =======================================================

# Force la console à interpréter le script en UTF-8 pour afficher correctement les accents et emojis
$PSDefaultParameterValues['*:Encoding'] = 'utf8'

Write-Host "Début du script de mise à jour de l'environnement Docker."

# --- Étape 1 : Sélection de la branche ---
Write-Host "`n1️⃣  Récupération des dernières informations du dépôt..."
git fetch --all > $null 2>&1

$branches = git branch -r | ForEach-Object { $_.Trim() } | Where-Object { $_ -notmatch "->" }
Write-Host "`nBranches distantes disponibles :"
for ($i = 0; $i -lt $branches.Count; $i++) {
    Write-Host "$($i+1). $($branches[$i])"
}

# Boucle de validation pour le choix de l'utilisateur
$choice = $null
while ($true) {
    $input = Read-Host "`nChoisis le numéro de la branche à déployer"
    if ($input -match '^\d+$' -and [int]$input -ge 1 -and [int]$input -le $branches.Count) {
        $choice = [int]$input
        break # Sort de la boucle si le choix est valide
    } else {
        Write-Host "Choix invalide. Merci de taper un numéro entre 1 et $($branches.Count)."
    }
}

$remoteBranchName = $branches[$choice - 1]
$branchName = $remoteBranchName -replace "origin/", ""
Write-Host "`nBranche '$branchName' sélectionnée."

# --- Le reste du script est identique ---

# --- Étape 2 : Mise à jour du code source (Méthode Forcée) ---
Write-Host "`n2️⃣  Mise à jour du code source depuis Git..."
try {
    $beforeReset = git rev-parse HEAD
    git checkout $branchName
    Write-Host "Forçage de la mise à jour pour correspondre à la version distante (git reset --hard)..."
    git reset --hard $remoteBranchName
    $afterReset = git rev-parse HEAD
} catch {
    Write-Host "Erreur lors de la mise à jour depuis Git : $_"
    exit
}

Write-Host "`nCommits récents sur '$branchName' :"
git log -5 --oneline

# --- Étape 3 : Logique de reconstruction Docker ---
Write-Host "`n3️⃣  Analyse des changements pour Docker..."
$filesChanged = git diff --name-only $beforeReset $afterReset

$needsRebuild = $false
if ($filesChanged -match "go.mod" -or `
    $filesChanged -match "package.json" -or `
    $filesChanged -match "Dockerfile" -or `
    $filesChanged -match "docker-compose.yml") {
    $needsRebuild = $true
}

# --- Étape 4 : Lancement avec Docker Compose ---
Write-Host "`n4️⃣  Mise à jour de l'environnement Docker..."
if ($needsRebuild) {
    Write-Host "Changements détectés dans les dépendances ou la configuration. Reconstruction des images..."
    docker-compose up --build -d
} else {
    Write-Host "Aucun changement de dépendances. Redémarrage des services..."
    docker-compose up -d
    Write-Host "Redémarrage du service backend pour appliquer les modifications de code..."
    docker-compose restart backend
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Une erreur est survenue avec Docker Compose. Vérifiez les messages ci-dessus."
    exit
}

Write-Host "`n✅ Environnement démarré avec succès !"
Write-Host "   - Frontend accessible sur http://localhost:3000"
Write-Host "   - Backend accessible sur http://localhost:8080"

# --- Étape 5 : Affichage des logs ---
Write-Host "`n5️⃣  Affichage des logs en direct (appuyez sur Ctrl+C pour quitter)..."
docker-compose logs -f --tail="50"
