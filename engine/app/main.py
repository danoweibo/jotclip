from fastapi import FastAPI
from app.routers import health
from app.grpc.server import serve
from dotenv import load_dotenv
from google import genai
import threading
import os

load_dotenv()

client = genai.Client(api_key=os.getenv("GEMINI_API_KEY"))

app = FastAPI(title="Jotclip Engine")

app.include_router(health.router)

@app.on_event("startup")
async def startup():
    # Start gRPC server in background thread
    grpc_thread = threading.Thread(target=serve, daemon=True)
    grpc_thread.start()
    print("✅ Jotclip Engine running")
    print(f"✅ Gemini configured: {bool(os.getenv('GEMINI_API_KEY'))}")