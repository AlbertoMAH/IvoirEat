from playwright.sync_api import sync_playwright, expect
import time

def run_auth_verification():
    # Wait for the dev server to be ready.
    print("Waiting 15 seconds for the dev server to start...")
    time.sleep(15)
    print("Continuing with the script.")

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Use a unique email for each run to avoid registration conflicts
        unique_email = f"testuser_{int(time.time())}@example.com"
        password = "password123"

        try:
            # 1. Go to register page and create a new user
            print("Navigating to register page...")
            page.goto("http://localhost:3000/register")

            # Take an immediate screenshot for debugging
            print("Taking initial screenshot to debug...")
            page.screenshot(path="jules-scratch/verification/debug_register_page.png")

            page.get_by_label("Nom complet").fill("Test User")
            page.get_by_label("Adresse e-mail").fill(unique_email)
            page.get_by_label("Mot de passe").fill(password)
            page.get_by_role("button", name="S'inscrire").click()

            print("Registration form submitted.")

            # 2. Wait for success message and redirection to login
            expect(page.get_by_text("Inscription réussie !")).to_be_visible(timeout=5000)
            print("Registration successful, waiting for redirection...")
            page.wait_for_url("**/login", timeout=5000)
            print("Redirected to login page.")

            # 3. Fill login form and submit
            expect(page.get_by_role("heading", name="Connexion")).to_be_visible()
            page.get_by_label("Adresse e-mail").fill(unique_email)
            page.get_by_label("Mot de passe").fill(password)
            page.get_by_role("button", name="Se connecter").click()
            print("Login form submitted.")

            # 4. Wait for redirection to home page and verify logged-in state
            page.wait_for_url("**/", timeout=5000)
            print("Redirected to home page.")
            expect(page.get_by_text("Bienvenue ! Vous êtes connecté.")).to_be_visible()
            print("Login successful, welcome message is visible.")

            # 5. Test the protected route
            page.get_by_role("button", name="Tester la route /api/validate").click()
            print("Protected route button clicked.")

            # 6. Verify the response from the protected route
            response_element = page.locator("pre")
            expect(response_element).to_be_visible(timeout=5000)
            expect(response_element).to_contain_text(f'"Email": "{unique_email}"')
            print("Protected route tested successfully.")

            # 7. Take a screenshot of the final state
            page.screenshot(path="jules-scratch/verification/auth_flow_success.png")
            print("Screenshot taken successfully.")

        except Exception as e:
            print(f"An error occurred during verification: {e}")
            page.screenshot(path="jules-scratch/verification/auth_flow_error.png")
            raise # Re-raise the exception to fail the script

        finally:
            browser.close()

if __name__ == "__main__":
    run_auth_verification()
