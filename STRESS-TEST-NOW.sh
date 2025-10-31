#!/bin/bash

echo "⚡ STRESS TEST - HIGH PERFORMANCE"
echo "================================="

if ! command -v ab &> /dev/null; then
    echo "Installing Apache Bench..."
    brew install httpd 2>/dev/null || echo "Install apache2-utils manually"
fi

URL="http://localhost:8080/health"
REQUESTS=1000
CONCURRENT=50

echo "Testing: $URL"
echo "Requests: $REQUESTS"
echo "Concurrent: $CONCURRENT"
echo ""

ab -n $REQUESTS -c $CONCURRENT $URL

echo ""
echo "✅ Stress test complete"
