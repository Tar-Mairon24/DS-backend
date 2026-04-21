#!/usr/bin/env bash

set -a
source "$(dirname "$0")/../.env"
set +a

BASE_URL=http://localhost:8080/api/v1
COOKIE_JAR="$(dirname "$0")/.cookies.txt"
IMAGES_DIR="$(dirname "$0")/propertiesImages"

# Check if cookie jar exists
if [ ! -f "$COOKIE_JAR" ]; then
    echo "Error: Cookie jar not found at $COOKIE_JAR"
    echo "Please run seed_auth.sh first to generate authentication tokens"
    exit 1
fi

# Check if images directory exists
if [ ! -d "$IMAGES_DIR" ]; then
    echo "Error: Images directory not found at $IMAGES_DIR"
    echo "Please create $IMAGES_DIR and organize images in subdirectories by property ID"
    exit 1
fi

upload_images() {
    local property_id="$1"
    local property_dir="$IMAGES_DIR/$property_id"
    
    if [ ! -d "$property_dir" ]; then
        echo "Warning: No images directory for property $property_id"
        return
    fi
    
    local images=($(find "$property_dir" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.webp" \) | sort))
    
    if [ ${#images[@]} -eq 0 ]; then
        echo "Warning: No images found for property $property_id"
        return
    fi
    
    echo "Uploading ${#images[@]} images for property $property_id..."
    
    for i in "${!images[@]}"; do
        local image_path="${images[$i]}"
        local filename=$(basename "$image_path")
        local description="Propiedad #$property_id - Imagen $((i + 1))"
        
        # First image is main, others are not
        if [ $i -eq 0 ]; then
            local main_image="true"
            echo "  ✓ Uploading $filename (MAIN)"
        else
            local main_image="false"
            echo "  ✓ Uploading $filename"
        fi
        
        response=$(curl -s -o /dev/null -w "%{http_code}" -X 'POST' \
            --location "$BASE_URL/properties/$property_id/images" \
            -b "$COOKIE_JAR" \
            -c "$COOKIE_JAR" \
            -F "image=@\"$image_path\"" \
            -F "description=$description" \
            -F "main_image=$main_image")
        
        if [ "$response" -eq 201 ] || [ "$response" -eq 200 ]; then
            echo "    Status: $response ✓"
        else
            echo "    Status: $response ✗ Failed"
        fi
    done
}

for property_id in 1 2 3 4 5; do
    upload_images "$property_id"
done

echo "Image seeding complete!"
