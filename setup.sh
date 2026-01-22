#!/bin/bash

# Signal Bot Setup Script
# This script helps you set up and configure the Signal bot

set -e

echo "================================================"
echo "  Signal Bot Setup"
echo "================================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ -f .env ]; then
    echo -e "${YELLOW}Warning: .env file already exists${NC}"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled."
        exit 0
    fi
fi

# Copy .env.example to .env
echo "Creating .env file..."
cp .env.example .env

# Prompt for phone number
echo ""
echo "Enter your Signal phone number in international format (e.g., +1234567890):"
read -p "Phone number: " phone_number

# Validate phone number format
if [[ ! $phone_number =~ ^\+[0-9]{10,15}$ ]]; then
    echo -e "${RED}Error: Phone number must be in international format (+1234567890)${NC}"
    exit 1
fi

# Update .env file using printf to safely escape the value
phone_number_escaped=$(printf '%s\n' "$phone_number" | sed -e 's/[\/&]/\\&/g')
sed -i.bak "s/SIGNAL_PHONE_NUMBER=/SIGNAL_PHONE_NUMBER=$phone_number_escaped/" .env
rm .env.bak 2>/dev/null || true

echo -e "${GREEN}✓ Configuration saved to .env${NC}"
echo ""

# Ask about command prefix
echo "What command prefix would you like to use? (default: !)"
read -p "Prefix: " prefix
if [ ! -z "$prefix" ]; then
    prefix_escaped=$(printf '%s\n' "$prefix" | sed -e 's/[\/&]/\\&/g')
    sed -i.bak "s/BOT_COMMAND_PREFIX=!/BOT_COMMAND_PREFIX=$prefix_escaped/" .env
    rm .env.bak 2>/dev/null || true
fi

echo ""
echo -e "${GREEN}✓ Setup complete!${NC}"
echo ""
echo "================================================"
echo "  Next Steps"
echo "================================================"
echo ""
echo "1. Install signal-cli if you haven't already:"
echo "   - Linux: See README.md for instructions"
echo "   - macOS: brew install signal-cli"
echo ""
echo "2. Register your Signal number:"
echo "   signal-cli -a $phone_number register"
echo ""
echo "3. Verify with the code you receive via SMS:"
echo "   signal-cli -a $phone_number verify CODE"
echo ""
echo "4. Build and run the bot:"
echo "   go build -o signal-bot"
echo "   ./signal-bot"
echo ""
echo "Alternatively, use Docker:"
echo "   docker-compose run --rm signal-bot signal-cli -a $phone_number register"
echo "   docker-compose run --rm signal-bot signal-cli -a $phone_number verify CODE"
echo "   docker-compose up -d"
echo ""
echo "================================================"
echo "For more information, see README.md"
echo "================================================"
