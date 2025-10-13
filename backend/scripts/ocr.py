import json
import sys
import datetime

def process_receipt_simulation():
    """
    Simulates OCR processing on a receipt image and returns structured data.
    In a real-world scenario, this function would use a library like Tesseract
    to extract text from an image file.
    """
    # Hardcoded data for simulation purposes
    ocr_data = {
      "montant": 105.00,
      "date": "2025-10-13",
      "tvac": 15.00,
      "marchand": "Supermarché Fictif",
      "type_recu": "Supermarché",
      "anomalie": False
    }

    return ocr_data

if __name__ == "__main__":
    # In this simulation, we don't need to read the file path from command-line arguments,
    # as the result is hardcoded.
    #
    # if len(sys.argv) < 2:
    #     print(json.dumps({"error": "No file path provided"}), file=sys.stderr)
    #     sys.exit(1)

    result = process_receipt_simulation()

    # Print the resulting JSON to standard output.
    # The Go application will capture this output.
    print(json.dumps(result))
