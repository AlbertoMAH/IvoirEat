# jules-scratch/verification/register.py
from playwright.sync_api import sync_playwright

def run(playwright):
    browser = playwright.chromium.launch()
    page = browser.new_page()
    page.goto("http://localhost:3000/register")
    page.get_by_label("Name").fill("Test User")
    page.get_by_label("Email").fill("test@example.com")
    page.get_by_label("Password").fill("password")
    page.get_by_role("button", name="S'inscrire").click()
    page.wait_for_url("http://localhost:3000/login")
    browser.close()

with sync_playwright() as playwright:
    run(playwright)
