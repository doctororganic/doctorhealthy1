#!/bin/bash

# Kill processes using port 8080
PORT=8080

echo "🔍 Checking for processes on port $PORT..."

# Find processes using the port
PIDS=$(lsof -ti:$PORT 2>/dev/null)

if [ -z "$PIDS" ]; then
    echo "✅ No processes found on port $PORT"
    exit 0
fi

echo "⚠️  Found processes on port $PORT:"
lsof -i:$PORT

echo "🛑 Killing processes..."
for PID in $PIDS; do
    echo "   Killing PID: $PID"
    kill -9 $PID 2>/dev/null
done

sleep 1

# Verify
REMAINING=$(lsof -ti:$PORT 2>/dev/null)
if [ -z "$REMAINING" ]; then
    echo "✅ Port $PORT is now free"
else
    echo "❌ Failed to free port $PORT"
    exit 1
fi

