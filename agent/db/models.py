from typing import List, Optional
from sqlmodel import Field, SQLModel, Relationship


class NoteTopicLink(SQLModel, table=True):
    note_id: Optional[int] = Field(
        default=None, foreign_key="note.id", primary_key=True
    )
    topic_id: Optional[int] = Field(
        default=None, foreign_key="topic.id", primary_key=True
    )


class NoteBase(SQLModel):
    content: str


class TopicBase(SQLModel):
    name: str


class Note(NoteBase, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    topics: List["Topic"] = Relationship(
        back_populates="notes", link_model=NoteTopicLink
    )


class Topic(TopicBase, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    notes: List[Note] = Relationship(back_populates="topics", link_model=NoteTopicLink)
