from typing import List, Optional
from sqlmodel import Field, SQLModel, Relationship, Column, String
from datetime import datetime, timezone


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

    # gorm.Model equivalent fields
    id: Optional[int] = Field(default=None, primary_key=True)
    created_at: Optional[datetime] = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        sa_column=Column("created_at"),
    )
    updated_at: Optional[datetime] = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        sa_column=Column("updated_at"),
    )
    deleted_at: Optional[datetime] = Field(default=None, sa_column=Column("deleted_at"))

    # Your custom fields
    content: str = Field(sa_column=Column("content", String))

    # Many-to-many relationship with Topic
    topics: List["Topic"] = Relationship(
        back_populates="notes", link_model=NoteTopicLink
    )


class Topic(SQLModel, table=True):
    __tablename__ = "topics"  # Match Go's table name

    # gorm.Model equivalent fields
    id: Optional[int] = Field(default=None, primary_key=True)
    created_at: Optional[datetime] = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        sa_column=Column("created_at"),
    )
    updated_at: Optional[datetime] = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        sa_column=Column("updated_at"),
    )
    deleted_at: Optional[datetime] = Field(default=None, sa_column=Column("deleted_at"))

    # Your custom fields
    topic: str = Field(sa_column=Column("topic", String))

    # Many-to-many relationship with Note
    notes: List[Note] = Relationship(back_populates="topics", link_model=NoteTopicLink)


class NoteBase(SQLModel):
    content: str = Field(sa_column=Column("content", String))


class TopicBase(SQLModel):
    topic: str = Field(sa_column=Column("topic", String))
