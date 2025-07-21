#!/bin/bash

# =============================================================================
# Cinema Booking API - Endpoint Testing Script
# =============================================================================

echo "🎬 Cinema Booking API - Endpoint Testing"
echo "========================================"

# Configuration
BASE_URL="http://localhost:3000"
ADMIN_EMAIL="admin@cinema.com"
ADMIN_PASSWORD="admin123"
USER_EMAIL="john@example.com"
USER_PASSWORD="password123"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✅ PASS${NC}: $message"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}❌ FAIL${NC}: $message"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    elif [ "$status" = "INFO" ]; then
        echo -e "${BLUE}ℹ️  INFO${NC}: $message"
    elif [ "$status" = "WARN" ]; then
        echo -e "${YELLOW}⚠️  WARN${NC}: $message"
    fi
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local description=$4
    local data=$5
    local token=$6
    
    # Build curl command
    local curl_cmd="curl -s -w '%{http_code}' -X $method"
    
    if [ ! -z "$token" ]; then
        curl_cmd="$curl_cmd -H 'Authorization: Bearer $token'"
    fi
    
    if [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$data'"
    fi
    
    curl_cmd="$curl_cmd '$BASE_URL$endpoint'"
    
    # Execute request
    local response=$(eval $curl_cmd)
    local status_code="${response: -3}"
    
    if [ "$status_code" = "$expected_status" ]; then
        print_status "PASS" "$description ($method $endpoint)"
    else
        print_status "FAIL" "$description ($method $endpoint) - Expected: $expected_status, Got: $status_code"
    fi
}

# Check if server is running
echo ""
print_status "INFO" "Checking if server is running..."
server_check=$(curl -s -w '%{http_code}' "$BASE_URL/api/films" -o /dev/null)
if [ "$server_check" != "200" ]; then
    print_status "FAIL" "Server is not running at $BASE_URL"
    echo "Please start the server with: go run main.go"
    exit 1
fi
print_status "PASS" "Server is running"

# Authentication Tests
echo ""
echo "🔐 Authentication Tests"
echo "----------------------"

# Login as admin
admin_login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" \
    "$BASE_URL/api/auth/login")

if echo "$admin_login_response" | grep -q "token"; then
    ADMIN_TOKEN=$(echo "$admin_login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    print_status "PASS" "Admin login successful"
else
    print_status "FAIL" "Admin login failed"
    ADMIN_TOKEN=""
fi

# Login as user
user_login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\"}" \
    "$BASE_URL/api/auth/login")

if echo "$user_login_response" | grep -q "token"; then
    USER_TOKEN=$(echo "$user_login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    print_status "PASS" "User login successful"
else
    print_status "FAIL" "User login failed"
    USER_TOKEN=""
fi

# Film Endpoints Tests
echo ""
echo "🎬 Film Endpoints Tests"
echo "----------------------"

test_endpoint "GET" "/api/films" "200" "Get all films (public)"
test_endpoint "GET" "/api/films/507f1f77bcf86cd799439011" "404" "Get non-existent film"

if [ ! -z "$ADMIN_TOKEN" ]; then
    test_endpoint "POST" "/api/films" "201" "Create film (admin)" \
        '{"title":"Test Film","description":"Test Description","duration":120,"genre":"Action","rating":"PG-13"}' \
        "$ADMIN_TOKEN"
fi

# Schedule Endpoints Tests
echo ""
echo "📅 Schedule Endpoints Tests"
echo "--------------------------"

test_endpoint "GET" "/api/jadwals" "200" "Get all schedules (public)"
test_endpoint "GET" "/api/jadwals/detail" "200" "Get schedules with film details (public)"
test_endpoint "GET" "/api/jadwals/507f1f77bcf86cd799439011" "404" "Get non-existent schedule"

# Ticket Endpoints Tests
echo ""
echo "🎫 Ticket Endpoints Tests"
echo "------------------------"

if [ ! -z "$ADMIN_TOKEN" ]; then
    test_endpoint "GET" "/api/tikets" "200" "Get all tickets (admin)" "" "$ADMIN_TOKEN"
fi

if [ ! -z "$USER_TOKEN" ]; then
    test_endpoint "GET" "/api/tikets/user/507f1f77bcf86cd799439011" "200" "Get user tickets" "" "$USER_TOKEN"
fi

# Payment Endpoints Tests
echo ""
echo "💳 Payment Endpoints Tests"
echo "-------------------------"

if [ ! -z "$ADMIN_TOKEN" ]; then
    test_endpoint "GET" "/api/pembayarans" "200" "Get all payments (admin)" "" "$ADMIN_TOKEN"
fi

if [ ! -z "$USER_TOKEN" ]; then
    test_endpoint "GET" "/api/pembayarans/user/507f1f77bcf86cd799439011" "200" "Get user payments" "" "$USER_TOKEN"
fi

# Protected Endpoint Tests
echo ""
echo "🛡️  Protected Endpoint Tests"
echo "---------------------------"

test_endpoint "GET" "/api/auth/profile" "401" "Access profile without token"

if [ ! -z "$USER_TOKEN" ]; then
    test_endpoint "GET" "/api/auth/profile" "200" "Access profile with user token" "" "$USER_TOKEN"
fi

# Admin Only Tests
echo ""
echo "👑 Admin Only Tests"
echo "------------------"

test_endpoint "POST" "/api/films" "401" "Create film without token" \
    '{"title":"Test","description":"Test","duration":120,"genre":"Action","rating":"PG-13"}'

if [ ! -z "$USER_TOKEN" ]; then
    test_endpoint "POST" "/api/films" "403" "Create film with user token" \
        '{"title":"Test","description":"Test","duration":120,"genre":"Action","rating":"PG-13"}' \
        "$USER_TOKEN"
fi

# Invalid ObjectID Tests
echo ""
echo "🆔 ObjectID Validation Tests"
echo "---------------------------"

test_endpoint "GET" "/api/films/invalid-id" "400" "Get film with invalid ObjectID"
test_endpoint "GET" "/api/jadwals/invalid-id" "400" "Get schedule with invalid ObjectID"
test_endpoint "GET" "/api/tikets/invalid-id" "400" "Get ticket with invalid ObjectID" "" "$ADMIN_TOKEN"
test_endpoint "GET" "/api/pembayarans/invalid-id" "400" "Get payment with invalid ObjectID" "" "$ADMIN_TOKEN"

# Summary
echo ""
echo "📊 Test Summary"
echo "==============="
echo "Total Tests: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    print_status "PASS" "All tests passed! 🎉"
    echo ""
    echo "✅ Your API is working correctly!"
    echo "✅ All endpoints are accessible"
    echo "✅ Authentication is working"
    echo "✅ ObjectID validation is working"
    echo "✅ Admin permissions are enforced"
    echo ""
    echo "🚀 Your API is ready for production!"
else
    print_status "FAIL" "$FAILED_TESTS test(s) failed"
    echo ""
    echo "❌ Please check the failed tests above"
    echo "💡 Make sure the database is seeded"
    echo "💡 Verify all controller functions exist"
    echo "💡 Check route definitions"
fi

echo ""
echo "🔗 Next Steps:"
echo "- Import Postman collection from docs/postman/"
echo "- Run manual tests using the API Testing Checklist"
echo "- Test the complete user journey end-to-end"
