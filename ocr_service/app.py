from flask import Flask, request, jsonify
import os
from ocr import process_receipt # On importe la logique qu'on a déjà écrite !

app = Flask(__name__)

@app.route('/ocr', methods=['POST'])
def ocr_endpoint():
    if 'receiptImage' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    file = request.files['receiptImage']

    # Sauvegarde temporaire du fichier pour l'analyse
    temp_dir = "temp_uploads"
    os.makedirs(temp_dir, exist_ok=True)
    temp_path = os.path.join(temp_dir, file.filename)
    file.save(temp_path)

    # On appelle notre fonction OCR existante
    try:
        result = process_receipt(temp_path)
        return jsonify(result)
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        # Nettoyage du fichier temporaire
        if os.path.exists(temp_path):
            os.remove(temp_path)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8081)))
