import easyocr
import io
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse
import torch

app = FastAPI()

# Initialize the OCR reader. This will download the model on the first run.
# We specify French and English as languages.
reader = easyocr.Reader(['fr', 'en'], gpu=False) # Set gpu=False as we might not have a GPU

@app.get("/")
def read_root():
    return {"message": "OCR Service is running"}

@app.post("/ocr")
async def perform_ocr(file: UploadFile = File(...)):
    """
    Performs OCR on an uploaded image file.
    """
    try:
        # Read the file content
        contents = await file.read()

        # Use easyocr to read text from the image bytes
        # The detail=0 parameter returns a list of strings
        results = reader.readtext(contents, detail=0, paragraph=True)

        # Join the results into a single block of text
        extracted_text = "\n".join(results)

        return JSONResponse(content={
            "filename": file.filename,
            "text": extracted_text
        })
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
