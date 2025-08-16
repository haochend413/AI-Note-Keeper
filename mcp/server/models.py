from typing import List, Optional
from sqlmodel import Field, SQLModel, Relationship, Column, String


class NoteTopicLink(SQLModel, table=True):
    __tablename__ = "note_topics"  # Join table is already plural, keep as is
    note_id: Optional[int] = Field(
        default=None, foreign_key="notes.id", primary_key=True
    )
    topic_id: Optional[int] = Field(
        default=None, foreign_key="topics.id", primary_key=True
    )


class Note(SQLModel, table=True):
    __tablename__ = "notes"  # Match Go's table name
    id: Optional[int] = Field(default=None, primary_key=True)
    content: str = Field(sa_column=Column("content", String))
    topics: List["Topic"] = Relationship(
        back_populates="notes", link_model=NoteTopicLink
    )


class Topic(SQLModel, table=True):
    __tablename__ = "topics"  # Match Go's table name
    id: Optional[int] = Field(default=None, primary_key=True)
    topic: str = Field(sa_column=Column("topic", String))
    notes: List[Note] = Relationship(back_populates="topics", link_model=NoteTopicLink)


class NoteBase(SQLModel):
    content: str = Field(sa_column=Column("content", String))


class TopicBase(SQLModel):
    topic: str = Field(sa_column=Column("topic", String))
