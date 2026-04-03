#!/usr/bin/env bash

set -a
source "$(dirname "$0")/../.env"
set +a

BASE_URL=http://localhost:8080/api/v1
COOKIE_JAR="$(dirname "$0")/.cookies.txt"

echo "Creating initial admin user..."
curl -s -X 'POST' \
    --location "$BASE_URL/users" \
    -H "Content-Type: application/json" \
    -c "$COOKIE_JAR" \
    -d '{
        "email": "'"$USER1_EMAIL"'",
        "username": "'"$USER1_NAME"'",
        "password": "'"$USER1_PASSWORD"'",
        "role": "admin",
        "phone": "8444428728"
    }'

echo "Logging in and saving session..."
curl -s -X 'POST' \
    --location "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -c "$COOKIE_JAR" \
    -d '{
        "email": "'"$USER1_EMAIL"'",
        "password": "'"$USER1_PASSWORD"'"
    }'

echo "Authentication tokens saved to $COOKIE_JAR"