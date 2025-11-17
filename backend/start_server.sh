#!/bin/bash
cd /Users/khaledahmedmohamed/Desktop/trae\ new\ healthy1/nutrition-platform/backend
echo "Starting nutrition platform backend server..."
nohup go run main.go > server.log 2>&1 &
echo $! > server.pid
echo "Server started with PID: $(cat server.pid)"
echo "Log file: server.log"
sleep 2
echo "Checking if server is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Server is running on port 8080"
else
    echo "❌ Server failed to start. Check server.log for errors."
fi
