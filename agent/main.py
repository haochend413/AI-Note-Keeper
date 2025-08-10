# main.py
from fastapi import FastAPI
from contextlib import asynccontextmanager
from db.db import db_router, engine
from sqlmodel import Session
from LLM.agent import llm_router, ingest_notes, fetchData


@asynccontextmanager
async def lifespan(app: FastAPI):
    print("fetching embeddings ... ")
    fetchData(Session(engine))
    ingest_notes(Session(engine))
    yield


app = FastAPI(lifespan=lifespan)

app.include_router(db_router, prefix="/db")
app.include_router(llm_router, prefix="/llm")
