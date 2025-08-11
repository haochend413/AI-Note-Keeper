#!/bin/bash

echo "Activating virtual environment..."
source venv/bin/activate

# Kill any existing mantis session
if tmux has-session -t mantis 2>/dev/null; then
    echo "Killing existing tmux session: mantis"
    tmux kill-session -t mantis
fi

# Stop and remove existing Qdrant container (if any)
docker rm -f qdrant_container 2>/dev/null

# Start Qdrant container detached with a fixed name
docker run -d --name qdrant_container -p 6333:6333 -v $(pwd)/data/qdrant:/qdrant/storage qdrant/qdrant

echo "Starting new tmux session with four panes..."

# Start tmux session with first pane tailing Qdrant logs to keep pane alive
tmux new-session -d -s mantis "docker logs -f qdrant_container"

# Split vertically from pane 0 → pane 1 (Ollama)
tmux split-window -v -t mantis:0 "source venv/bin/activate && ollama serve"

# Select pane 1 and split horizontally → pane 2 (FastAPI)
tmux select-pane -t 1
tmux split-window -h -t mantis:0 "source venv/bin/activate && uvicorn main:app --host 0.0.0.0 --port 8000 --reload"

# Select pane 1 (or 2) and split vertically → pane 3 (blank shell)
tmux select-pane -t 2
tmux split-window -v -t mantis:0

# Arrange layout nicely
tmux select-layout -t mantis:0 tiled

# Attach session
tmux attach-session -t mantis
