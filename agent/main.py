# main.py
from fastapi import FastAPI
from contextlib import asynccontextmanager
from db.db import create_db_and_tables, db_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    # onstartup, create;
    print("ffff")
    create_db_and_tables()
    yield


app = FastAPI(lifespan=lifespan)


app.include_router(db_router, prefix="/db")
