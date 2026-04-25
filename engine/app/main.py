from fastapi import FastAPI
from app.routers import health
from dotenv import load_dotenv
from google import genai
import os

load_dotenv()

client = genai.Client(api_key=os.getenv("GEMINI_API_KEY"))

app = FastAPI(title="Jotclip Engine")

app.include_router(health.router)

@app.on_event("startup")
async def startup():
    print("✅ Jotclip Engine running")
    print(f"✅ Gemini configured: {bool(os.getenv('GEMINI_API_KEY'))}")