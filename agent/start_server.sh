
# setup folders if non-exist

echo "Activating virtual environment..."
source venv/bin/activate
# mkcert will not work for distribution.
echo "Starting local server..."
uvicorn main:app \
    --host 0.0.0.0 \
    --port 8000 \
    --reload



