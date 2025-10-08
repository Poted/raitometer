import os
import google.generativeai as genai
from fastapi import FastAPI, UploadFile, File, HTTPException
from PIL import Image
from dotenv import load_dotenv

load_dotenv()

API_KEY = os.getenv('GOOGLE_API_KEY')
if not API_KEY:
    raise ValueError("GOOGLE_API_KEY environment variable not set")

genai.configure(api_key=API_KEY)
model = genai.GenerativeModel("gemini-1.5-flash-latest")

app = FastAPI(title="raitometer AI Service")

@app.get("/healthcheck")
def healthcheck():
    return {"status": "ok", "service": "AI Service"}

@app.post("/analyze-image")
async def analyze_image(file: UploadFile = File(...)):
    try:
        image = Image.open(file.file)
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid image file")

    prompt = [
        "You are an expert in construction and renovation quotes for the European market.",
        "Analyze this image of a room. Your response MUST be in JSON format only.",
        "Your task is to identify and count specific objects, estimate surface areas, and assess the condition of the walls.",
        "Do not add any explanations or markdown formatting. The JSON output must be structured as follows:",
        """
        {
          "countable_items": [
            {"item": "window", "count": <number>},
            {"item": "door", "count": <number>},
            {"item": "light_point", "count": <number>}
          ],
          "estimated_areas": [
            {"area_type": "wall", "sqm_estimate": <number>}
          ],
          "condition_assessment": {
            "walls": "<'good', 'average', or 'poor'>"
          }
        }
        """,
        "If you cannot determine a value, use 0 for counts and 'unknown' for condition.",
        image,
    ]

    try:
        response = model.generate_content(prompt)
        # We assume the model returns a clean JSON string.
        # In a production app, we would add more robust parsing and error handling here.
        return response.text
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error during AI analysis: {str(e)}")