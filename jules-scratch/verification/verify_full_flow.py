from playwright.sync_api import sync_playwright, expect
import time

def run_verification(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    try:
        # Step 1: Register a new user
        page.goto("http://localhost:3000/register")
        page.get_by_label("Name").fill("Test User")
        page.get_by_label("Email").fill(f"testuser_{int(time.time())}@example.com")
        page.get_by_label("Password").fill("password123")
        page.get_by_role("button", name="Register").click()
        expect(page).to_have_url("http://localhost:3000/login", timeout=10000)
        print("Registration successful.")

        # Step 2: Log in
        page.get_by_label("Email").fill(f"testuser_{int(time.time())}@example.com")
        page.get_by_label("Password").fill("password123")
        page.get_by_role("button", name="Login").click()
        expect(page).to_have_url("http://localhost:3000/", timeout=10000)
        print("Login successful.")

        # Step 3: Navigate to Upload page and upload a receipt
        page.goto("http://localhost:3000/receipts/upload")
        expect(page.get_by_role("heading", name="Upload a Receipt")).to_be_visible()

        # Use a dummy file for the upload
        file_path = "jules-scratch/verification/dummy_receipt.png"
        with open(file_path, "w") as f:
            f.write("This is a dummy receipt file.")

        page.get_by_label("Receipt Image").set_input_files(file_path)
        page.get_by_role("button", name="Upload").click()

        # Wait for the success message
        expect(page.get_by_text("Receipt uploaded successfully!")).to_be_visible(timeout=10000)
        print("Receipt upload successful.")

        # Step 4: Navigate to the list page and verify the receipt is there
        page.goto("http://localhost:3000/receipts/list")
        expect(page.get_by_role("heading", name="Your Receipts")).to_be_visible()

        # Check for the merchant name from the OCR simulation
        expect(page.get_by_text("Supermarché Fictif")).to_be_visible()
        print("Receipt visible in list.")

        # Step 5: Take a screenshot
        page.screenshot(path="jules-scratch/verification/verification.png")
        print("Screenshot taken.")

    except Exception as e:
        print(f"An error occurred: {e}")
        page.screenshot(path="jules-scratch/verification/error.png")

    finally:
        browser.close()

with sync_playwright() as playwright:
    run_verification(playwright)
