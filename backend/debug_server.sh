#!/bin/bash
cd /Users/khaledahmedmohamed/Desktop/trae\ new\ healthy1/nutrition-platform/backend
echo "Debugging nutrition platform backend server..."
echo "Current directory: $(pwd)"
echo "Go files found:"
ls -la *.go 2>/dev/null || echo "No .go files found"
echo "Go module:"
cat go.mod | head -5
echo "Testing go run:"
go run main.go &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"
sleep 5
echo "Checking processes:"
ps aux | grep $SERVER_PID | grep -v grep || echo "Server process not found"
echo "Checking port 8080:"
lsof -i:8080 || echo "Port 8080 not in use"
echo "Testing health endpoint:"
curl -s http://localhost:8080/health || echo "Health endpoint failed"
echo "Killing server process if running:"
kill $SERVER_PID 2>/dev/null || echo "Server already stopped"
