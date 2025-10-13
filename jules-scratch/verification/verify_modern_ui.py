from playwright.sync_api import sync_playwright, expect
import time

def run_verification(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    try:
        # Step 1: Login
        page.goto("http://localhost:3000/login")
        email = f"testuser_{int(time.time())}@example.com"
        # We assume the user is already registered for this test
        page.get_by_label("Email").fill("test@example.com")
        page.get_by_label("Password").fill("password123")
        page.get_by_role("button", name="Login").click()
        expect(page).to_have_url("http://localhost:3000/", timeout=10000)
        print("Login successful.")

        # Step 2: Check for Navbar
        expect(page.get_by_role("navigation")).to_be_visible()
        expect(page.get_by_role("link", name="IvoirEat")).to_be_visible()
        print("Navbar is visible.")

        # Step 3: Go to Upload page and check the new UI
        page.goto("http://localhost:3000/receipts/upload")
        expect(page.get_by_text("Cliquez pour choisir un fichier")).to_be_visible()
        print("Modern upload page UI is visible.")
        page.screenshot(path="jules-scratch/verification/modern_upload_page.png")

        # Step 4: Go to List page and check the new UI
        page.goto("http://localhost:3000/receipts/list")
        expect(page.get_by_role("heading", name="Mes Notes de Frais")).to_be_visible()
        print("Modern receipt list page UI is visible.")
        page.screenshot(path="jules-scratch/verification/modern_list_page.png")

        print("Verification script completed successfully.")

    except Exception as e:
        print(f"An error occurred: {e}")
        page.screenshot(path="jules-scratch/verification/error_modern_ui.png")

    finally:
        browser.close()

with sync_playwright() as playwright:
    run_verification(playwright)
