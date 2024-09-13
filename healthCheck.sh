#!/bin/bash

# Perform the first GET request
if ! curl -f http://localhost:8080/health/router > /dev/null 2>&1; then
    echo "Health check failed for /health/router"
    exit 1
fi

# Perform the second GET request
if ! curl -f http://localhost:8080/health/database > /dev/null 2>&1; then
    echo "Health check failed for /health/database"
    exit 1
fi

echo "Health check passed"
exit 0
