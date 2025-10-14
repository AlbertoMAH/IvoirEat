import json
import sys
import re
import easyocr
from datetime import datetime

# Initialize the OCR reader
# This will download the model the first time it's run
reader = easyocr.Reader(['fr', 'en'])

def extract_text_from_image(image_path):
    """Extracts raw text from an image using EasyOCR."""
    try:
        result = reader.readtext(image_path, detail=0, paragraph=True)
        return " ".join(result)
    except Exception as e:
        return f"Error during OCR processing: {e}"

def find_total(text):
    """Finds the total amount in the OCR text using regex."""
    # Regex to find amounts, looking for keywords like TOTAL, TTC, etc.
    # This regex is quite permissive to handle various formats.
    patterns = [
        r"(?:TOTAL|TOTAL\s*TTC|MONTANT\s*TTC|TOTAL\s*A\s*PAYER)\s*[:\s]*([\d\s,.]+[€$]?)",
        r"([\d\s,.]+[€$]?)\s*(?:TOTAL|TTC)"
    ]

    highest_amount = 0.0

    for pattern in patterns:
        matches = re.findall(pattern, text, re.IGNORECASE)
        for match in matches:
            # Clean up the matched string
            amount_str = re.sub(r"[^\d,.]", "", match).replace(",", ".")
            try:
                amount = float(amount_str)
                if amount > highest_amount:
                    highest_amount = amount
            except ValueError:
                continue

    # If no specific total is found, look for the highest amount on the receipt
    if highest_amount == 0.0:
        all_amounts = re.findall(r"(\d+[,.]\d{2})", text)
        for amount_str in all_amounts:
            try:
                amount = float(amount_str.replace(",", "."))
                if amount > highest_amount:
                    highest_amount = amount
            except ValueError:
                continue

    return highest_amount

def find_date(text):
    """Finds the date in the OCR text and normalizes it."""
    # Regex for various date formats (dd/mm/yyyy, dd-mm-yyyy, yyyy-mm-dd)
    match = re.search(r"(\d{2}[-/]\d{2}[-/]\d{4})|(\d{4}[-/]\d{2}[-/]\d{2})", text)
    if match:
        date_str = match.group(0)
        # Normalize the date format to YYYY-MM-DD
        try:
            # Try parsing dd/mm/yyyy or dd-mm-yyyy
            parsed_date = datetime.strptime(date_str, "%d/%m/%Y")
        except ValueError:
            try:
                parsed_date = datetime.strptime(date_str, "%d-%m-%Y")
            except ValueError:
                try:
                    # Try parsing yyyy-mm-dd
                    parsed_date = datetime.strptime(date_str, "%Y-%m-%d")
                except ValueError:
                    return None
        return parsed_date.strftime("%Y-%m-%d")
    return None

def find_merchant(text_lines):
    """Finds the merchant name, often in the first few lines."""
    if text_lines:
        # A simple heuristic: the merchant is often the first non-trivial line
        for line in text_lines.split("\n"):
            if line.strip():
                return line.strip()
    return "Unknown Merchant"

def classify_receipt(text):
    """Classifies the receipt based on keywords."""
    text_lower = text.lower()

    categories = {
        "Transport": ["taxi", "trajet", "course", "station", "carburant", "ticket", "vol"],
        "Restaurant": ["restaurant", "repas", "menu", "table", "boisson", "cafe", "snack"],
        "Hébergement": ["hotel", "chambre", "nuitée", "séjour"],
        "Achats": ["supermarché", "articles", "ttc", "boutique", "fournitures"],
        "Services": ["abonnement", "licence", "honoraires", "consultation", "réparation"],
        "Santé": ["pharmacie", "medicament", "consultation", "soin"]
    }

    for category, keywords in categories.items():
        if any(keyword in text_lower for keyword in keywords):
            return category

    return "Divers"

def process_receipt(image_path):
    """
    Main function to process a receipt image file.
    It orchestrates the text extraction, data finding, and classification.
    """

    # 1. Extract raw text from the image
    raw_text = extract_text_from_image(image_path)
    if "Error" in raw_text:
        return {"error": raw_text}

    # 2. Extract structured data from the raw text
    total_amount = find_total(raw_text)
    receipt_date = find_date(raw_text)
    merchant_name = find_merchant(raw_text) # Simple version
    receipt_type = classify_receipt(raw_text)

    # 3. Assemble the JSON output
    result = {
        "montant": total_amount,
        "date": receipt_date if receipt_date else datetime.now().strftime("%Y-%m-%d"),
        "tvac": 0.0, # VAT extraction is complex, defaulting to 0 for now
        "marchand": merchant_name,
        "type_recu": receipt_type,
        "anomalie": False # Anomaly detection to be implemented later
    }

    return result

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No file path provided"}), file=sys.stderr)
        sys.exit(1)

    image_file_path = sys.argv[1]

    # Process the receipt and get the structured data
    ocr_result = process_receipt(image_file_path)

    # Print the resulting JSON to standard output
    print(json.dumps(ocr_result, indent=2))
