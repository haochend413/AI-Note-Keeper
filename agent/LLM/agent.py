from fastapi import APIRouter
from db.models import Note, Topic, NoteTopicLink
from typing import List
from sqlmodel import Session, select
from sqlalchemy.orm import selectinload
from db.db import engine, SessionDep
from langchain_ollama import OllamaEmbeddings
from langchain_community.vectorstores import Qdrant
from qdrant_client import QdrantClient
from langchain.chains import RetrievalQA
from langchain_community.llms import Ollama


raw_data = {"notes": [], "topics": [], "links": []}

llm_router = APIRouter()


# fetch everything;
def fetchData(session: SessionDep):
    global raw_data
    with session:
        raw_data["notes"] = session.exec(
            select(Note).options(selectinload(Note.topics))
        ).all()
        raw_data["topics"] = session.exec(select(Topic)).all()
        raw_data["links"] = session.exec(select(NoteTopicLink)).all()


# vectorization;
embedding_model = OllamaEmbeddings(model="nomic-embed-text")
qdrant_client = QdrantClient(host="localhost", port=6333)
vectorstore = None  # Will init on ingest


# These will be set after vectorstore is initialized
llm = None
retriever = None
qa_chain = None


# run at startup;


def ingest_notes(session: SessionDep):
    # erase previous record before re-starting;
    qdrant_client.delete_collection(collection_name="notes")
    notes = raw_data["notes"]
    docs = []
    payload = []
    for note in notes:
        topics = ", ".join(topic.topic for topic in note.topics)
        t_ids = ", ".join(str(topic.id) for topic in note.topics)
        text = f"Note: {note.content}\nTopics: {topics}"
        docs.append(text)
        payload.append({"note_id": note.id, "topic_ids": t_ids, "topic_names": topics})

    global vectorstore, llm, retriever, qa_chain
    vectorstore = Qdrant.from_texts(
        texts=docs,
        embedding=embedding_model,
        location="http://localhost:6333",
        collection_name="notes",
        metadatas=payload,
    )
    llm = Ollama(model="mistral")
    retriever = vectorstore.as_retriever(search_kwargs={"k": 10})
    qa_chain = RetrievalQA.from_chain_type(llm=llm, retriever=retriever)
    return {"status": "indexed"}


# query API
@llm_router.get("/query")
def query_llm(q: str):
    if vectorstore is None or qa_chain is None:
        return {"error": "Vectorstore not initialized. Run /ingest first."}
        # Get retrieved docs for debugging
    docs = retriever.get_relevant_documents(q)
    print("Retrieved docs:", docs)

    result = qa_chain.invoke({"query": q})
    return {"answer": result["result"]}


# Summary API;
