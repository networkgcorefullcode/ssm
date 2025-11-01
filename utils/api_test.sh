#!/bin/bash

# SSM API Test Script
# Usage: ./api_test.sh <endpoint> <method> <json_file> [base_url]
# Example: ./api_test.sh /encrypt POST ../docs/json_examples/requests/encrypt_request.json

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print usage
print_usage() {
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $0 <endpoint> <method> <json_file> [base_url]"
    echo ""
    echo -e "${YELLOW}Arguments:${NC}"
    echo "  endpoint   - API endpoint (e.g., /encrypt, /decrypt, /generate-aes-key)"
    echo "  method     - HTTP method (GET, POST, PUT, DELETE)"
    echo "  json_file  - Path to JSON file with request body"
    echo "  base_url   - Optional base URL (default: http://localhost:8080)"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 /encrypt POST ../docs/json_examples/requests/encrypt_request.json"
    echo "  $0 /health-check GET"
    echo "  $0 /decrypt POST ../docs/json_examples/requests/decrypt_request.json http://localhost:9000"
    echo ""
    echo -e "${YELLOW}Available endpoints:${NC}"
    echo "  /health-check         - Health check (GET/POST)"
    echo "  /encrypt              - Encrypt data (POST)"
    echo "  /decrypt              - Decrypt data (POST)"
    echo "  /generate-aes-key     - Generate AES key (POST)"
    echo "  /generate-des-key     - Generate DES key (POST)"
    echo "  /generate-des3-key    - Generate DES3 key (POST)"
    echo "  /store-key            - Store key (POST)"
    echo "  /get-data-keys        - Get multiple keys by label (POST)"
    echo "  /get-key              - Get single key by label (POST)"
    echo "  /get-all-keys         - Get all keys (POST)"
}

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed${NC}"
    echo "Please install curl to use this script"
    exit 1
fi

# Check if jq is installed (optional but recommended for pretty output)
JQ_INSTALLED=false
if command -v jq &> /dev/null; then
    JQ_INSTALLED=true
fi

# Parse arguments
if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    print_usage
    exit 0
fi

if [ $# -lt 2 ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo ""
    print_usage
    exit 1
fi

ENDPOINT="$1"
METHOD=$(echo "$2" | tr '[:lower:]' '[:upper:]')
JSON_FILE="$3"
BASE_URL="${4:-http://localhost:8080}"

# Validate HTTP method
if [[ ! "$METHOD" =~ ^(GET|POST|PUT|DELETE|PATCH)$ ]]; then
    echo -e "${RED}Error: Invalid HTTP method '$METHOD'${NC}"
    echo "Valid methods: GET, POST, PUT, DELETE, PATCH"
    exit 1
fi

# Build the full URL
FULL_URL="${BASE_URL}${ENDPOINT}"

echo -e "${GREEN}=== SSM API Request ===${NC}"
echo -e "${YELLOW}URL:${NC}     $FULL_URL"
echo -e "${YELLOW}Method:${NC}  $METHOD"

# Prepare curl command
CURL_CMD="curl -k -s -w '\n\nHTTP Status: %{http_code}\nTime: %{time_total}s\n'"

# Add method
CURL_CMD="$CURL_CMD -X $METHOD"

# Add headers
CURL_CMD="$CURL_CMD -H 'Content-Type: application/json'"
CURL_CMD="$CURL_CMD -H 'Accept: application/json'"

# Add data if JSON file is provided
if [ -n "$JSON_FILE" ]; then
    if [ ! -f "$JSON_FILE" ]; then
        echo -e "${RED}Error: JSON file not found: $JSON_FILE${NC}"
        exit 1
    fi
    echo -e "${YELLOW}Body:${NC}    $JSON_FILE"
    echo -e "${YELLOW}Content:${NC}"
    if [ "$JQ_INSTALLED" = true ]; then
        jq '.' "$JSON_FILE"
    else
        cat "$JSON_FILE"
    fi
    CURL_CMD="$CURL_CMD -d @$JSON_FILE"
else
    echo -e "${YELLOW}Body:${NC}    (none)"
fi

echo ""
echo -e "${GREEN}=== Response ===${NC}"

# Execute curl command
RESPONSE=$(eval "$CURL_CMD '$FULL_URL'")

# Pretty print response if jq is available
if [ "$JQ_INSTALLED" = true ]; then
    echo "$RESPONSE" | head -n -2 | jq '.' 2>/dev/null || echo "$RESPONSE" | head -n -2
else
    echo "$RESPONSE" | head -n -2
fi

# Print status and time
echo ""
echo "$RESPONSE" | tail -n 2

# Extract HTTP status code
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP Status:" | awk '{print $3}')

# Color code the result
echo ""
if [ "$HTTP_STATUS" -ge 200 ] && [ "$HTTP_STATUS" -lt 300 ]; then
    echo -e "${GREEN}✓ Request successful${NC}"
    exit 0
elif [ "$HTTP_STATUS" -ge 400 ] && [ "$HTTP_STATUS" -lt 500 ]; then
    echo -e "${YELLOW}⚠ Client error${NC}"
    exit 1
elif [ "$HTTP_STATUS" -ge 500 ]; then
    echo -e "${RED}✗ Server error${NC}"
    exit 1
else
    echo -e "${YELLOW}? Unknown status${NC}"
    exit 1
fi
