#!/bin/bash

# Script to set up the notes-app project structure with empty files

# Prompt for module path
read -p "Enter your Go module path (e.g., github.com/yourusername/notes-app, default: github.com/example/notes-app): " MODULE_PATH
MODULE_PATH=${MODULE_PATH:-github.com/example/notes-app}

# Create project directory structure
mkdir -p notes-app/{cmd/notes,internal/{db,models,ui,app}}
cd notes-app || exit 1
echo "Created project directory: notes-app"

# Create empty Go files
touch cmd/notes/main.go
touch internal/db/db.go
touch internal/db/sync.go
touch internal/models/note.go
touch internal/models/topic.go
touch internal/models/daily_task.go
touch internal/ui/model.go
touch internal/ui/update.go
touch internal/ui/view.go
touch internal/ui/styles.go
touch internal/app/app.go
echo "Created empty Go files"

# Create go.mod with module path
cat > go.mod << EOL
module $MODULE_PATH

go 1.21
EOL
echo "Created go.mod with module path: $MODULE_PATH"

# Run go mod tidy to initialize dependencies
echo "Running go mod tidy..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "Successfully initialized go.mod"
else
    echo "Error running go mod tidy"
    exit 1
fi

echo "Setup complete! Project structure created in notes-app/"