#!/usr/bin/env bash

set -a
source "$(dirname "$0")/../.env"
set +a

BASE_URL=http://localhost:8080/api/v1

post_data() {
    local name="$1"
    local body="$2"

    echo "Seeding $name..."
    response=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
        --location "$BASE_URL/$name" \
        -H "Content-Type: application/json" \
        -d "$body")
    
    echo "Response status code: $response"

    if [ "$response" -eq 201 ]; then
        echo "$name seeded successfully."
    else
        echo "Failed to seed $name. HTTP status code: $response"
    fi
}

# Seed Users
post_data "users" '{
    "email": "'"$USER1_EMAIL"'",
    "username": "'"$USER1_NAME"'",
    "password": "'"$USER1_PASSWORD"'",
    "role": "admin"
}'