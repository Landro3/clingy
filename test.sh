#!/bin/bash
echo "Setting up client folders..."
mkdir -p client-1
mkdir -p client-2
cp -rf clingy-client/* client-1/
cp -rf clingy-client/* client-2/
echo "Finished copying files"

echo "Setting UI env variables"
touch client-1/ui/.env
echo "API_URL=http://localhost:8888/api" >client-1/ui/.env
touch client-2/ui/.env
echo "API_URL=http://localhost:8989/api" >client-2/ui/.env

echo "Opening tmux with 2x3 layout..."

SESSION_NAME="clingy-session"

# Kill existing session if it exists
tmux has-session -t $SESSION_NAME 2>/dev/null
if [ $? -eq 0 ]; then
    echo "Killing existing tmux session: $SESSION_NAME"
    tmux kill-session -t $SESSION_NAME
fi

# Create new detached session
echo "Creating new tmux session: $SESSION_NAME"
tmux new-session -d -s $SESSION_NAME

# Setup the layout: 2 panes in top row, 3 panes in bottom row
echo "Creating layout..."

# Split horizontally to create top and bottom rows
tmux split-window -v -t $SESSION_NAME:0

# Split top pane horizontally to create 2 panes in top row
tmux split-window -h -t $SESSION_NAME:0.0

# Split bottom-left pane horizontally to create first 2 panes in bottom row
tmux split-window -h -t $SESSION_NAME:0.2

# Split bottom-middle pane horizontally to create third pane in bottom row
# tmux split-window -h -t $SESSION_NAME:0.3

# Bottom-left pane: client-1 Go server
tmux send-keys -t $SESSION_NAME:0.2 'cd client-1/api && go run .' C-m

# Bottom-middle pane: client-2 Go server on port 8989
tmux send-keys -t $SESSION_NAME:0.3 'cd client-2/api && go run . -port 8989' C-m

# Top-left pane: client-1 UI
tmux send-keys -t $SESSION_NAME:0.0 'cd client-1/ui && npm run dev' C-m

# Top-right pane: client-2 UI  
tmux send-keys -t $SESSION_NAME:0.1 'cd client-2/ui && npm run dev' C-m

# Bottom-right pane: helper
# tmux send-keys -t $SESSION_NAME:0.4 'echo "Helper pane ready"' C-m

# Select the top-left pane as active
tmux select-pane -t $SESSION_NAME:0.0

# Attach to the session
tmux attach-session -t $SESSION_NAME

