# main.py
from fastapi import FastAPI

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello, FastAPI!"}


@app.post("/ask")
async def ask(data: dict):
    user_input = data.get("query")
    # call your LLM agent logic here
    return {"response": f"You said: {user_input}"}
