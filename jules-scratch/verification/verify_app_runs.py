from playwright.sync_api import sync_playwright, Page, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    page = browser.new_page()

    # 1. Arrange: Go to the application homepage.
    # The dev server is typically on port 3000 or 5173.
    page.goto("http://localhost:5173", wait_until="domcontentloaded")

    # 2. Assert: Check for a known element on the page to confirm it loaded.
    # We'll just wait for the body to be visible as a basic check.
    expect(page.locator("body")).to_be_visible()

    # 3. Screenshot: Capture the final result for visual verification.
    page.screenshot(path="jules-scratch/verification/verification.png")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
