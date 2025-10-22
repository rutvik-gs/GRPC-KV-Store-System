##################################################################
###             THIS IS AI GENERATED !!!!                      ###
##################################################################

#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080"
MAX_WAIT=60
WAIT_INTERVAL=2

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  KV Store E2E Integration Tests${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print colored status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Function to wait for service
wait_for_service() {
    echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
    elapsed=0
    while [ $elapsed -lt $MAX_WAIT ]; do
        if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Services are ready!${NC}"
            return 0
        fi
        sleep $WAIT_INTERVAL
        elapsed=$((elapsed + WAIT_INTERVAL))
        echo -n "."
    done
    echo -e "${RED}✗ Timeout waiting for services${NC}"
    return 1
}

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}🧹 Cleaning up...${NC}"
    docker-compose down > /dev/null 2>&1
    echo -e "${GREEN}✓ Cleanup complete${NC}"
}

# Trap exit to ensure cleanup
trap cleanup EXIT

# Step 1: Clean up any existing containers
echo -e "${YELLOW}📦 Step 1: Cleaning up existing containers...${NC}"
docker-compose down > /dev/null 2>&1
print_status $? "Cleanup complete"

# Step 2: Build Docker images
echo -e "${YELLOW}🔨 Step 2: Building Docker images...${NC}"
docker-compose build
print_status $? "Docker images built"

# Step 3: Start services
echo -e "${YELLOW}🚀 Step 3: Starting services...${NC}"
docker-compose up -d
print_status $? "Services started"

# Step 4: Wait for services to be ready
wait_for_service
print_status $? "Services health check"

# Give services a moment to stabilize
sleep 2

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Running API Tests${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to run test
run_test() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5

    echo -e "${YELLOW}Testing: ${test_name}${NC}"

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ $test_name (Status: $status_code)${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ $test_name (Expected: $expected_status, Got: $status_code)${NC}"
        echo -e "${RED}  Response: $body${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Run Tests

echo -e "${BLUE}--- Health Check ---${NC}"
run_test "Health Check" "GET" "/health" "" "200"

echo ""
echo -e "${BLUE}--- Set Operations ---${NC}"
run_test "Set key-value (success)" "POST" "/kv" '{"key":"username","value":"alice"}' "201"
run_test "Set without key (error)" "POST" "/kv" '{"key":"","value":"test"}' "400"
run_test "Set with invalid JSON (error)" "POST" "/kv" '{invalid}' "400"

echo ""
echo -e "${BLUE}--- Get Operations ---${NC}"
run_test "Get existing key" "GET" "/kv/username" "" "200"
run_test "Get non-existent key" "GET" "/kv/nonexistent-key-12345" "" "404"

echo ""
echo -e "${BLUE}--- Delete Operations ---${NC}"
run_test "Delete existing key" "DELETE" "/kv/username" "" "200"
run_test "Delete non-existent key" "DELETE" "/kv/nonexistent-key-99999" "" "404"

echo ""
echo -e "${BLUE}--- Verification ---${NC}"
run_test "Verify deleted key is gone" "GET" "/kv/username" "" "404"

echo ""
echo -e "${BLUE}--- Multiple Operations ---${NC}"
run_test "Set namespace key" "POST" "/kv" '{"key":"user:1:email","value":"alice@example.com"}' "201"
run_test "Get namespace key" "GET" "/kv/user:1:email" "" "200"
run_test "Delete namespace key" "DELETE" "/kv/user:1:email" "" "200"

# Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Passed: $TESTS_PASSED${NC}"
echo -e "${RED}✗ Failed: $TESTS_FAILED${NC}"
echo -e "${BLUE}Total: $((TESTS_PASSED + TESTS_FAILED))${NC}"
echo ""

# Show service logs if there were failures
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "${YELLOW}📋 Service Logs:${NC}"
    docker-compose logs --tail=20
fi

# Exit with appropriate code
if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed!${NC}"
    exit 1
fi
