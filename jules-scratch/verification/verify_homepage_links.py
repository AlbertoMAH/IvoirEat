from playwright.sync_api import sync_playwright, expect
import time

def run_verification(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    try:
        # Step 1: Register and Log in
        page.goto("http://localhost:3000/register")
        email = f"testuser_{int(time.time())}@example.com"
        page.get_by_label("Name").fill("Test User Nav")
        page.get_by_label("Email").fill(email)
        page.get_by_label("Password").fill("password123")
        page.get_by_role("button", name="Register").click()
        expect(page).to_have_url("http://localhost:3000/login", timeout=10000)

        page.get_by_label("Email").fill(email)
        page.get_by_label("Password").fill("password123")
        page.get_by_role("button", name="Login").click()
        expect(page).to_have_url("http://localhost:3000/", timeout=10000)
        print("Login successful.")

        # Step 2: Verify the navigation links are visible
        expect(page.get_by_role("link", name="Uploader un Reçu")).to_be_visible()
        expect(page.get_by_role("link", name="Voir mes Reçus")).to_be_visible()
        expect(page.get_by_role("button", name="Déconnexion")).to_be_visible()
        print("Navigation links for authenticated user are visible.")

        # Step 3: Take a screenshot
        page.screenshot(path="jules-scratch/verification/homepage_auth_links.png")
        print("Screenshot taken of the homepage with auth links.")

    except Exception as e:
        print(f"An error occurred: {e}")
        page.screenshot(path="jules-scratch/verification/error_homepage.png")

    finally:
        browser.close()

with sync_playwright() as playwright:
    run_verification(playwright)
