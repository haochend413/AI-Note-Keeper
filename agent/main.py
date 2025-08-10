# main.py
from fastapi import FastAPI
from contextlib import asynccontextmanager
from db.db import db_router


app = FastAPI()

app.include_router(db_router, prefix="/db")
