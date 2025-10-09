from playwright.sync_api import sync_playwright, expect
import os

def run_verification():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        try:
            # 1. Navigate to the login page
            page.goto("http://localhost:3000/login")

            # 2. Fill in the login form and submit
            # Using placeholder as a fallback if label is not found
            page.locator('input[name="email"]').fill(os.environ.get("ADMIN_EMAIL", "admin@example.com"))
            page.locator('input[name="password"]').fill(os.environ.get("ADMIN_PASSWORD", "password"))
            page.get_by_role("button", name="Login").click()

            # 3. Wait for navigation to the dashboard and for stats to load
            # The dashboard page is the root page after login in this setup
            expect(page).to_have_url("http://localhost:3000/", timeout=10000)

            # Wait for a specific element that indicates data has loaded.
            # Let's wait for the "Parkings" card title to be visible.
            expect(page.get_by_text("Parkings")).to_be_visible(timeout=15000)

            # 4. Take a screenshot of the dashboard
            screenshot_path = "jules-scratch/verification/dashboard_verification.png"
            page.screenshot(path=screenshot_path)
            print(f"Screenshot saved to {screenshot_path}")

        except Exception as e:
            print(f"An error occurred: {e}")
            page.screenshot(path="jules-scratch/verification/error.png")
        finally:
            browser.close()

if __name__ == "__main__":
    run_verification()
