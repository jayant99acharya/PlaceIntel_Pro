@echo off
setlocal enabledelayedexpansion

REM PlaceIntel Pro - Windows Setup Script
REM Foursquare Places API Hackathon Project

echo.
echo PlaceIntel Pro - Universal Location Intelligence Platform
echo ==================================================
echo Setting up your hackathon project...
echo.

REM Check if required tools are installed
:check_requirements
echo Checking requirements...

REM Check Docker
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not installed. Please install Docker Desktop first.
    echo Visit: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM Check Docker Compose
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker Compose is not installed. Please install Docker Compose first.
    pause
    exit /b 1
)

REM Check Go (optional)
go version >nul 2>&1
if %errorlevel% equ 0 (
    for /f "tokens=3" %%i in ('go version') do echo Go %%i found
) else (
    echo Go not found (optional for local development)
)

REM Check Python (optional)
python --version >nul 2>&1
if %errorlevel% equ 0 (
    for /f "tokens=2" %%i in ('python --version') do echo Python %%i found
) else (
    echo Python not found (optional for local development)
)

for /f "tokens=3" %%i in ('docker --version') do echo Docker %%i found
for /f "tokens=3" %%i in ('docker-compose --version') do echo Docker Compose %%i found
echo.

REM Setup environment file
:setup_environment
echo Setting up environment...

if not exist .env (
    copy .env .env >nul
    echo Created .env file from template
    echo Please edit .env and add your Foursquare API key
    echo.
    echo To get your Foursquare API key:
    echo 1. Visit: https://foursquare.com/developers/signup
    echo 2. Create a free account
    echo 3. Create a new app and copy the API key
    echo 4. Edit .env file and set FOURSQUARE_API_KEY=your_api_key_here
    echo.
    pause
) else (
    echo .env file already exists
)

REM Check if API key is set
findstr /c:"your_foursquare_api_key_here" .env >nul
if %errorlevel% equ 0 (
    echo Please set your Foursquare API key in .env file
    pause
    exit /b 1
)

echo Environment configured
echo.

REM Build and start services
:start_services
echo Building and starting services...
echo This may take a few minutes on first run...
echo.

docker-compose up --build -d
if %errorlevel% neq 0 (
    echo Failed to start services
    pause
    exit /b 1
)

echo.
echo Services started successfully!
echo.

REM Wait for services to be ready
:wait_for_services
echo Waiting for services to be ready...

REM Wait for API service
echo Checking API service...
set /a counter=0
:api_check_loop
set /a counter+=1
curl -s http://localhost:8081/api/v1/health >nul 2>&1
if %errorlevel% equ 0 (
    echo API service is ready
    goto intelligence_check
)
if %counter% geq 30 (
    echo API service failed to start
    docker-compose logs api
    pause
    exit /b 1
)
timeout /t 2 /nobreak >nul
goto api_check_loop

:intelligence_check
REM Wait for Intelligence service
echo Checking Intelligence service...
set /a counter=0
:intelligence_check_loop
set /a counter+=1
curl -s http://localhost:5000/health >nul 2>&1
if %errorlevel% equ 0 (
    echo Intelligence service is ready
    goto test_api
)
if %counter% geq 30 (
    echo Intelligence service failed to start
    docker-compose logs intelligence
    pause
    exit /b 1
)
timeout /t 2 /nobreak >nul
goto intelligence_check_loop

REM Test the API
:test_api
echo.
echo Testing API endpoints...

REM Test health endpoint
echo Testing health endpoint...
curl -s http://localhost:8081/api/v1/health | findstr "healthy" >nul
if %errorlevel% equ 0 (
    echo Health check passed
) else (
    echo Health check failed
    pause
    exit /b 1
)

REM Test search endpoint with sample data
echo Testing search endpoint...
curl -s "http://localhost:8081/api/v1/places/search?lat=40.7128&lng=-74.0060&query=coffee&limit=1" > temp_result.json
findstr "results" temp_result.json >nul
if %errorlevel% equ 0 (
    echo Search endpoint working
    del temp_result.json
) else (
    echo Search endpoint failed
    type temp_result.json
    del temp_result.json
    pause
    exit /b 1
)

echo.

REM Show success message and next steps
:show_success
echo PlaceIntel Pro is now running successfully!
echo.
echo Access your application:
echo • API Documentation: http://localhost:8081/docs
echo • Demo Application: Open examples\demo.html in your browser
echo • API Health Check: http://localhost:8081/api/v1/health
echo • Intelligence Service: http://localhost:5000/health
echo.
echo Sample API Endpoints:
echo • Search Places: GET http://localhost:8081/api/v1/places/search?lat=40.7128^&lng=-74.0060^&query=coffee
echo • Place Details: GET http://localhost:8081/api/v1/places/{place_id}/details
echo • Popular Places: GET http://localhost:8081/api/v1/analytics/popular?lat=40.7128^&lng=-74.0060
echo.
echo Development Commands:
echo • View logs: docker-compose logs -f
echo • Stop services: docker-compose down
echo • Restart services: docker-compose restart
echo • Rebuild services: docker-compose up --build
echo.
echo Tips:
echo • The demo.html file showcases all three intelligence pillars
echo • API responses include processing times and data sources
echo • All services are containerized for easy deployment
echo • Redis caching improves performance for repeated requests
echo.
echo Ready, Working
echo.

REM Open demo in Chrome
echo Opening demo application in Chrome...

REM Try to find and open Chrome
if exist "%ProgramFiles%\Google\Chrome\Application\chrome.exe" (
    start "" "%ProgramFiles%\Google\Chrome\Application\chrome.exe" "examples\demo.html"
) else if exist "%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe" (
    start "" "%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe" "examples\demo.html"
) else if exist "%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe" (
    start "" "%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe" "examples\demo.html"
) else (
    echo Chrome not found. Please open examples\demo.html manually in Chrome
    echo Opening in default browser instead...
    start examples\demo.html
)

pause