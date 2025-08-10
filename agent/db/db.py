from fastapi import Depends, FastAPI, HTTPException, Query, APIRouter
from sqlmodel import Field, Session, SQLModel, create_engine, select
from sqlmodel import SQLModel, create_engine, text
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker
import os
from dotenv import load_dotenv
from sqlalchemy.orm import sessionmaker, selectinload
from typing import AsyncGenerator, Annotated, List
from db.models import Note, Topic, NoteTopicLink, NoteBase, TopicBase
import asyncio
from pathlib import Path

db_router = APIRouter()

sqlite_file_name = (
    Path(__file__).parent.parent.parent / "client" / "main" / "notes.db"
).as_posix()
sqlite_url = f"sqlite:///{sqlite_file_name}"
engine = create_engine(sqlite_url, echo=True)


def get_session():
    with Session(engine) as session:
        yield session


SessionDep = Annotated[Session, Depends(get_session)]


# Input model for creating/updating a Note (no id, no topics by default)
class NoteCreate(NoteBase):
    pass


# Input model for creating/updating a Topic (no id)
class TopicCreate(TopicBase):
    pass


# Output model for Note (includes id and nested topics)
class NoteRead(NoteBase):
    id: int
    topics: List["TopicRead"] = []  # topics nested


# Output model for Topic (includes id and nested notes)
class TopicRead(TopicBase):
    id: int
    notes: List[NoteRead] = []


@db_router.post("/notes/", response_model=NoteRead)
def create_note(note_in: NoteCreate, session: SessionDep):
    note = Note(content=note_in.content)
    session.add(note)
    session.commit()
    session.refresh(note)
    # Eagerly load topics before returning
    note = session.exec(
        select(Note).where(Note.id == note.id).options(selectinload(Note.topics))
    ).one()
    return note


@db_router.post("/topics/", response_model=TopicRead)
def create_topic(topic_in: TopicCreate, session: SessionDep):
    topic = Topic(name=topic_in.name)
    session.add(topic)
    session.commit()
    session.refresh(topic)
    # Eagerly load notes before returning
    topic = session.exec(
        select(Topic).where(Topic.id == topic.id).options(selectinload(Topic.notes))
    ).one()
    return topic


@db_router.post("/notes/{note_id}/topics/{topic_id}")
def link_note_to_topic(note_id: int, topic_id: int, session: SessionDep):
    note = session.get(Note, note_id)
    topic = session.get(Topic, topic_id)
    if not note or not topic:
        raise HTTPException(status_code=404, detail="Note or Topic not found")
    if topic not in note.topics:
        note.topics.append(topic)
        session.add(note)
        session.commit()
    return {"ok": True}


@db_router.put("/notes/{note_id}", response_model=NoteRead)
def update_note(note_id: int, note_in: NoteCreate, session: SessionDep):
    note = session.get(Note, note_id)
    if not note:
        raise HTTPException(status_code=404, detail="Note not found")
    note.content = note_in.content
    session.add(note)
    session.commit()
    session.refresh(note)
    note = session.exec(
        select(Note).where(Note.id == note.id).options(selectinload(Note.topics))
    ).one()
    return note


@db_router.put("/topics/{topic_id}", response_model=TopicRead)
def update_topic(topic_id: int, topic_in: TopicCreate, session: SessionDep):
    topic = session.get(Topic, topic_id)
    if not topic:
        raise HTTPException(status_code=404, detail="Topic not found")
    topic.name = topic_in.name
    session.add(topic)
    session.commit()
    session.refresh(topic)
    topic = session.exec(
        select(Topic).where(Topic.id == topic.id).options(selectinload(Topic.notes))
    ).one()
    return topic


@db_router.delete("/notes/{note_id}")
def delete_note(note_id: int, session: SessionDep):
    note = session.get(Note, note_id)
    if not note:
        raise HTTPException(status_code=404, detail="Note not found")
    session.delete(note)
    session.commit()
    return {"ok": True}


@db_router.delete("/topics/{topic_id}")
def delete_topic(topic_id: int, session: SessionDep):
    topic = session.get(Topic, topic_id)
    if not topic:
        raise HTTPException(status_code=404, detail="Topic not found")
    session.delete(topic)
    session.commit()
    return {"ok": True}


@db_router.get("/notes/", response_model=List[NoteRead])
def read_notes(session: SessionDep):
    notes = session.exec(select(Note).options(selectinload(Note.topics))).all()
    return notes


@db_router.get("/topics/", response_model=List[TopicRead])
def read_topics(session: SessionDep):
    topics = session.exec(select(Topic).options(selectinload(Topic.notes))).all()
    return topics
