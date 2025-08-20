#!/bin/bash

# PlaceIntel Pro - Setup Script
# Foursquare Places API Hackathon Project

set -e

echo "🏆 PlaceIntel Pro - Universal Location Intelligence Platform"
echo "=================================================="
echo "Setting up your winning hackathon project..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if required tools are installed
check_requirements() {
    echo -e "${BLUE}Checking requirements...${NC}"
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
        exit 1
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
        exit 1
    fi
    
    # Check Go (optional, for local development)
    if command -v go &> /dev/null; then
        echo -e "${GREEN}✅ Go $(go version | cut -d' ' -f3) found${NC}"
    else
        echo -e "${YELLOW}⚠️  Go not found (optional for local development)${NC}"
    fi
    
    # Check Python (optional, for local development)
    if command -v python3 &> /dev/null; then
        echo -e "${GREEN}✅ Python $(python3 --version | cut -d' ' -f2) found${NC}"
    else
        echo -e "${YELLOW}⚠️  Python3 not found (optional for local development)${NC}"
    fi
    
    echo -e "${GREEN}✅ Docker $(docker --version | cut -d' ' -f3 | sed 's/,//') found${NC}"
    echo -e "${GREEN}✅ Docker Compose $(docker-compose --version | cut -d' ' -f3 | sed 's/,//') found${NC}"
    echo ""
}

# Setup environment file
setup_environment() {
    echo -e "${BLUE}Setting up environment...${NC}"
    
    if [ ! -f .env ]; then
        cp .env.example .env
        echo -e "${YELLOW}⚠️  Created .env file from template${NC}"
        echo -e "${YELLOW}⚠️  Please edit .env and add your Foursquare API key${NC}"
        echo ""
        echo -e "${BLUE}To get your Foursquare API key:${NC}"
        echo "1. Visit: https://foursquare.com/developers/signup"
        echo "2. Create a free account"
        echo "3. Create a new app and copy the API key"
        echo "4. Edit .env file and set FOURSQUARE_API_KEY=your_api_key_here"
        echo ""
        read -p "Press Enter after you've added your API key to .env file..."
    else
        echo -e "${GREEN}✅ .env file already exists${NC}"
    fi
    
    # Check if API key is set
    if grep -q "your_foursquare_api_key_here" .env; then
        echo -e "${RED}❌ Please set your Foursquare API key in .env file${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Environment configured${NC}"
    echo ""
}

# Build and start services
start_services() {
    echo -e "${BLUE}Building and starting services...${NC}"
    echo "This may take a few minutes on first run..."
    echo ""
    
    # Build and start services
    docker-compose up --build -d
    
    echo ""
    echo -e "${GREEN}✅ Services started successfully!${NC}"
    echo ""
}

# Wait for services to be ready
wait_for_services() {
    echo -e "${BLUE}Waiting for services to be ready...${NC}"
    
    # Wait for API service
    echo "Checking API service..."
    for i in {1..30}; do
        if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ API service is ready${NC}"
            break
        fi
        if [ $i -eq 30 ]; then
            echo -e "${RED}❌ API service failed to start${NC}"
            docker-compose logs api
            exit 1
        fi
        sleep 2
    done
    
    # Wait for Intelligence service
    echo "Checking Intelligence service..."
    for i in {1..30}; do
        if curl -s http://localhost:5000/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Intelligence service is ready${NC}"
            break
        fi
        if [ $i -eq 30 ]; then
            echo -e "${RED}❌ Intelligence service failed to start${NC}"
            docker-compose logs intelligence
            exit 1
        fi
        sleep 2
    done
    
    echo ""
}

# Test the API
test_api() {
    echo -e "${BLUE}Testing API endpoints...${NC}"
    
    # Test health endpoint
    echo "Testing health endpoint..."
    if curl -s http://localhost:8081/api/v1/health | grep -q "healthy"; then
        echo -e "${GREEN}✅ Health check passed${NC}"
    else
        echo -e "${RED}❌ Health check failed${NC}"
        exit 1
    fi
    
    # Test search endpoint with sample data
    echo "Testing search endpoint..."
    SEARCH_RESULT=$(curl -s "http://localhost:8081/api/v1/places/search?lat=40.7128&lng=-74.0060&query=coffee&limit=1")
    if echo "$SEARCH_RESULT" | grep -q "results"; then
        echo -e "${GREEN}✅ Search endpoint working${NC}"
    else
        echo -e "${RED}❌ Search endpoint failed${NC}"
        echo "Response: $SEARCH_RESULT"
        exit 1
    fi
    
    echo ""
}

# Show success message and next steps
show_success() {
    echo -e "${GREEN}🎉 PlaceIntel Pro is now running successfully!${NC}"
    echo ""
    echo -e "${BLUE}🚀 Access your application:${NC}"
    echo "• API Documentation: http://localhost:8081/docs"
    echo "• Demo Application: Open examples/demo.html in your browser"
    echo "• API Health Check: http://localhost:8081/api/v1/health"
    echo "• Intelligence Service: http://localhost:5000/health"
    echo ""
    echo -e "${BLUE}📊 Sample API Endpoints:${NC}"
    echo "• Search Places: GET http://localhost:8081/api/v1/places/search?lat=40.7128&lng=-74.0060&query=coffee"
    echo "• Place Details: GET http://localhost:8081/api/v1/places/{place_id}/details"
    echo "• Popular Places: GET http://localhost:8081/api/v1/analytics/popular?lat=40.7128&lng=-74.0060"
    echo ""
    echo -e "${BLUE}🛠️  Development Commands:${NC}"
    echo "• View logs: docker-compose logs -f"
    echo "• Stop services: docker-compose down"
    echo "• Restart services: docker-compose restart"
    echo "• Rebuild services: docker-compose up --build"
    echo ""
    echo -e "${YELLOW}💡 Hackathon Tips:${NC}"
    echo "• The demo.html file showcases all three intelligence pillars"
    echo "• API responses include processing times and data sources"
    echo "• All services are containerized for easy deployment"
    echo "• Redis caching improves performance for repeated requests"
    echo ""
    echo -e "${GREEN}🏆 Ready to win the hackathon! Good luck!${NC}"
    echo ""
    echo -e "${BLUE}Opening demo in Chrome...${NC}"
    
    # Try to open in Chrome (different commands for different OS)
    if command -v google-chrome &> /dev/null; then
        google-chrome examples/demo.html
    elif command -v google-chrome-stable &> /dev/null; then
        google-chrome-stable examples/demo.html
    elif command -v chromium-browser &> /dev/null; then
        chromium-browser examples/demo.html
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if [ -d "/Applications/Google Chrome.app" ]; then
            open -a "Google Chrome" examples/demo.html
        else
            echo -e "${YELLOW}Chrome not found. Please open examples/demo.html manually in Chrome${NC}"
        fi
    else
        echo -e "${YELLOW}Chrome not found. Please open examples/demo.html manually in Chrome${NC}"
    fi
}

# Main execution
main() {
    check_requirements
    setup_environment
    start_services
    wait_for_services
    test_api
    show_success
}

# Handle script interruption
trap 'echo -e "\n${RED}Setup interrupted. Cleaning up...${NC}"; docker-compose down; exit 1' INT

# Run main function
main