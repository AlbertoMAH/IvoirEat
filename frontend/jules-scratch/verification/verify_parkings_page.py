from playwright.sync_api import sync_playwright, expect
import os

def run_verification():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        try:
            # 1. Navigate to the login page (running on port 3001)
            page.goto("http://localhost:3001/login")

            # 2. Fill in the login form and submit
            page.get_by_label("Email").fill(os.environ.get("ADMIN_EMAIL", "admin@example.com"))
            page.get_by_label("Password").fill(os.environ.get("ADMIN_PASSWORD", "password"))
            page.get_by_role("button", name="Sign in").click()

            # 3. Wait for navigation to the dashboard
            expect(page).to_have_url("http://localhost:3001/", timeout=10000)

            # 4. Navigate to the parkings page
            page.goto("http://localhost:3001/parkings")

            # 5. Wait for the parkings page to load
            expect(page.get_by_text("Gestion des Parkings")).to_be_visible(timeout=15000)

            # 6. Take a screenshot of the parkings page
            screenshot_path = "jules-scratch/verification/parkings_page_verification.png"
            page.screenshot(path=screenshot_path)
            print(f"Screenshot saved to {screenshot_path}")

        except Exception as e:
            print(f"An error occurred: {e}")
            page.screenshot(path="jules-scratch/verification/error.png")
        finally:
            browser.close()

if __name__ == "__main__":
    run_verification()
