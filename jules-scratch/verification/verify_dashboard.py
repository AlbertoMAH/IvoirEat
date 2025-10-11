from playwright.sync_api import sync_playwright, Page, expect

def verify_dashboard(page: Page):
    """
    Vérifie que la page du dashboard se charge correctement et prend une capture d'écran.
    """
    # 1. Naviguer vers la page du dashboard
    # L'application redirige de / vers /dashboard
    page.goto("http://localhost:3000/")

    # 2. Attendre un élément clé du layout pour s'assurer que la page est chargée
    # On attend que le titre principal de la navbar soit visible
    heading = page.get_by_role("heading", name="Tableau de Bord Super Admin")
    expect(heading).to_be_visible(timeout=15000) # Augmentation du timeout

    # 3. Attendre que les données des KPI soient chargées (une des cartes)
    kpi_card = page.get_by_text("Parkings Actifs")
    expect(kpi_card).to_be_visible(timeout=10000)

    # 4. Prendre une capture d'écran de la page complète
    page.screenshot(path="jules-scratch/verification/dashboard.png")

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        verify_dashboard(page)
        browser.close()

if __name__ == "__main__":
    main()
