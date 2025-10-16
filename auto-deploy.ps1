# =======================================================
# Script de déploiement et de gestion de branches Docker
# =======================================================

# Force la console à interpréter le script en UTF-8 pour afficher correctement les accents et emojis
$PSDefaultParameterValues['*:Encoding'] = 'utf8'

# --- Boucle principale du menu ---
while ($true) {
    Write-Host "----------------------------------------------------" -ForegroundColor Green
    Write-Host "  Script de gestion de l'environnement Docker"
    Write-Host "----------------------------------------------------" -ForegroundColor Green

    # --- Étape 1 : Menu principal ---
    Write-Host "`n1️⃣  Récupération des dernières informations du dépôt..."
    git fetch --all --prune > $null 2>&1 # --prune nettoie les branches qui n'existent plus sur le distant

    $branches = git branch -r | ForEach-Object { $_.Trim() } | Where-Object { $_ -notmatch "->" }

    Write-Host "`nQue souhaitez-vous faire ?"
    Write-Host "   1. Mettre à jour et déployer une branche"
    Write-Host "   2. Supprimer une branche distante"
    Write-Host "   3. Quitter"

    $mainChoice = Read-Host "`nVotre choix"

    switch ($mainChoice) {
        "1" {
            # --- DÉPLOIEMENT ---
            Write-Host "`nBranches distantes disponibles pour le déploiement :"
            for ($i = 0; $i -lt $branches.Count; $i++) {
                Write-Host "   $($i+1). $($branches[$i])"
            }

            $deployChoice = $null
            while ($true) {
                $input = Read-Host "`nChoisis le numéro de la branche à déployer"
                if ($input -match '^\d+$' -and [int]$input -ge 1 -and [int]$input -le $branches.Count) {
                    $deployChoice = [int]$input
                    break
                } else { Write-Host "Choix invalide." }
            }

            $remoteBranchName = $branches[$deployChoice - 1]
            $branchName = $remoteBranchName -replace "origin/", ""
            Write-Host "`nBranche '$branchName' sélectionnée pour le déploiement."

            # --- Mise à jour du code (Méthode Forcée) ---
            Write-Host "`n2️⃣  Mise à jour du code source depuis Git..."
            try {
                git checkout $branchName
                git reset --hard $remoteBranchName
            } catch { Write-Host "Erreur lors de la mise à jour : $_"; exit }

            Write-Host "`nCommits récents sur '$branchName' :"
            git log -5 --oneline

            # --- Lancement avec Docker Compose ---
            Write-Host "`n4️⃣  Mise à jour de l'environnement Docker..."
            Write-Host "Reconstruction des images et redémarrage des services..."
            docker-compose up --build -d

            if ($LASTEXITCODE -ne 0) { Write-Host "❌ Erreur Docker Compose."; exit }

            Write-Host "`n✅ Environnement démarré avec succès !"
            Write-Host "   - Frontend: http://localhost:3000"
            Write-Host "   - Backend: http://localhost:8080"

            # --- Affichage des logs ---
            Write-Host "`n5️⃣  Affichage des logs (Ctrl+C pour revenir au menu)..."
            docker-compose logs -f --tail="50"
            # Après Ctrl+C, la boucle recommence
        }
        "2" {
            # --- SUPPRESSION ---
            Write-Host "`nBranches distantes disponibles pour la suppression :"
            for ($i = 0; $i -lt $branches.Count; $i++) {
                Write-Host "   $($i+1). $($branches[$i])"
            }

            $deleteChoice = $null
            while ($true) {
                $input = Read-Host "`nChoisis le numéro de la branche à supprimer"
                if ($input -match '^\d+$' -and [int]$input -ge 1 -and [int]$input -le $branches.Count) {
                    $deleteChoice = [int]$input
                    break
                } else { Write-Host "Choix invalide." }
            }

            $remoteBranchNameToDelete = $branches[$deleteChoice - 1]
            $branchNameToDelete = $remoteBranchNameToDelete -replace "origin/", ""

            $confirmation = Read-Host "Êtes-vous sûr de vouloir supprimer la branche '$branchNameToDelete' ? Cette action est irréversible. (o/n)"
            if ($confirmation -eq 'o') {
                Write-Host "Suppression de la branche distante '$branchNameToDelete'..."
                git push origin --delete $branchNameToDelete

                if (git branch --list $branchNameToDelete) {
                    Write-Host "Suppression de la branche locale '$branchNameToDelete'..."
                    git branch -d $branchNameToDelete
                }

                Write-Host "`n✅ Branche supprimée avec succès."
            } else {
                Write-Host "Suppression annulée."
            }
            Read-Host "Appuyez sur Entrée pour continuer..."
        }
        "3" {
            # --- QUITTER ---
            Write-Host "Fin du script."
            exit
        }
        default {
            Write-Host "Choix invalide, veuillez réessayer."
        }
    }
}
