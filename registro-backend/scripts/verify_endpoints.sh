#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "----------------------------------------"
echo "Cleaning up previous test data..."
echo "----------------------------------------"
# Attempt to delete the test user (requires sudo if docker needs it)
# We use || true to ensure script continues if user doesn't exist
sudo docker-compose -f docker/docker-compose.yml exec -T db psql -U user -d registro -c "DELETE FROM users WHERE email='testuser@example.com';" || true
echo "Cleanup attempted."
echo ""

echo "----------------------------------------"
echo "Testing Auth Service Endpoints"
echo "----------------------------------------"

# 1. Register
echo "1. Registering User (testuser@example.com)..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "Password123!",
    "first_name": "Test",
    "last_name": "User",
    "role": "student"
  }')
echo "Response: $REGISTER_RESPONSE"
echo ""

# 2. Login
echo "2. Logging In..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "Password123!"
  }')

echo "Response: $LOGIN_RESPONSE"

# Extract Token (Simple grep/cut to avoid jq dependency)
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "Login failed, no token received."
else 
    echo -e "\nToken received successfully."
    
    # 3. Protected Route
    echo "3. Accessing Protected Profile /api/v1/users/..."
    USER_RESPONSE=$(curl -s -X GET "$BASE_URL/users/" \
      -H "Authorization: Bearer $TOKEN")
    echo "Response: $USER_RESPONSE"
fi

echo -e "\n----------------------------------------"
echo "Done."
