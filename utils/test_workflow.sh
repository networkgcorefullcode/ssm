#!/bin/bash

# SSM API Test - Practical Example
# This script demonstrates a complete workflow using the SSM API

# Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n${CYAN}=== SSM API Complete Workflow Test ===${NC}"
echo -e "${YELLOW}This will test the complete encryption/decryption workflow${NC}\n"

BASE_URL="${1:-http://localhost:8080}"

# Step 1: Health Check
echo -e "\n${GREEN}[Step 1] Checking service health...${NC}"
./api_test.sh /health-check GET "" "$BASE_URL"

read -p "Press Enter to continue to Step 2..."

# Step 2: Generate AES Key
echo -e "\n${GREEN}[Step 2] Generating AES-256 key...${NC}"
./api_test.sh /generate-aes-key POST ../docs/json_examples/requests/gen_aes_key_request.json "$BASE_URL"

read -p "Press Enter to continue to Step 3..."

# Step 3: Encrypt Data
echo -e "\n${GREEN}[Step 3] Encrypting data...${NC}"
./api_test.sh /encrypt POST ../docs/json_examples/requests/encrypt_request.json "$BASE_URL"

read -p "Press Enter to continue to Step 4..."

# Step 4: Decrypt Data
echo -e "\n${GREEN}[Step 4] Decrypting data...${NC}"
echo -e "${YELLOW}Note: Update the decrypt_request.json with the cipher and IV from Step 3${NC}"
./api_test.sh /decrypt POST ../docs/json_examples/requests/decrypt_request.json "$BASE_URL"

read -p "Press Enter to continue to Step 5..."

# Step 5: Get All Keys
echo -e "\n${GREEN}[Step 5] Retrieving all keys...${NC}"
./api_test.sh /get-all-keys POST "" "$BASE_URL"

echo -e "\n${CYAN}=== Workflow Complete ===${NC}"
echo -e "${GREEN}All steps executed successfully!${NC}\n"
