from typing import Any
import httpx
from mcp.server.fastmcp import FastMCP
from sqlmodel import Field, Session, SQLModel, create_engine, select
from sqlalchemy.orm import selectinload
from pathlib import Path
from models import Note, Topic, NoteTopicLink, NoteBase, TopicBase

# Initialize FastMCP server
mcp = FastMCP("weather")

# Constants
sqlite_file_name = (
    Path(__file__).parent.parent.parent / "mts" / "cmd" / "notes" / "notes.db"
).as_posix()
sqlite_url = f"sqlite:///{sqlite_file_name}"
USER_AGENT = "mantis/1.0"
engine = create_engine(sqlite_url, echo=True)


def get_session():
    with Session(engine) as session:
        yield session


# now, add functions for READ_ONLY first (Later enable agent to change the database;)

# query helper functions;


@mcp.tool(description="Get a list of all not-deleted notes from the database.")
def read_notes(input: str = ""):
    """Get all notes from the database that haven't been deleted."""
    with Session(engine) as session:
        notes = session.exec(
            select(Note)
            .where(Note.deleted_at.is_(False) | Note.deleted_at.is_(None))
            .options(selectinload(Note.topics))
        ).all()
        return [note.model_dump() for note in notes]


@mcp.tool(
    description="Get a list of all topics from the database. No parameters needed."
)
def read_topics(input: str = ""):
    """Get all topics from the database. No input required."""
    with Session(engine) as session:
        topics = session.exec(select(Topic).options(selectinload(Topic.notes))).all()
        return [topic.model_dump() for topic in topics]


# run server;
if __name__ == "__main__":
    print(sqlite_file_name)
    print("starting...")
    mcp.run(transport="stdio")
    # print("started")
