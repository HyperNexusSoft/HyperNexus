#!/bin/bash
# TormentNexus — Full Platform Test Suite
# Run: bash scripts/test-all.sh

set -e

BASE_URL="https://demo.hypernexus.site"
KERNEL_URL="http://127.0.0.1:7778"
PASS=0
FAIL=0
ERRORS=""

log() {
    echo "[$(date '+%H:%M:%S')] $1"
}

pass() {
    echo "  ✅ $1"
    PASS=$((PASS + 1))
}

fail() {
    echo "  ❌ $1"
    FAIL=$((FAIL + 1))
    ERRORS="$ERRORS\n- $1"
}

test_health() {
    log "Testing Health Checks..."
    
    # Demo health
    RESULT=$(curl -s --max-time 5 "$BASE_URL/health" 2>/dev/null)
    if echo "$RESULT" | grep -q '"ok":true'; then
        pass "Demo health check"
    else
        fail "Demo health check"
    fi
    
    # Kernel health
    RESULT=$(curl -s --max-time 5 "$KERNEL_URL/health" 2>/dev/null)
    if echo "$RESULT" | grep -q '"ok":true'; then
        pass "Kernel health check"
    else
        fail "Kernel health check"
    fi
}

test_api_index() {
    log "Testing API Index..."
    
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/index" 2>/dev/null)
    if echo "$RESULT" | grep -q '"endpoints"'; then
        pass "API index returns endpoints"
    else
        fail "API index returns endpoints"
    fi
}

test_catalog_search() {
    log "Testing Catalog Search..."
    
    # Search by keyword
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/backlog/search?q=postgres&limit=5" 2>/dev/null)
    if echo "$RESULT" | grep -q '"results"'; then
        pass "Catalog search returns results"
    else
        fail "Catalog search returns results"
    fi
    
    # Get total count
    TOTAL=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get('total',0))" 2>/dev/null)
    if [ "$TOTAL" -gt 0 ]; then
        pass "Catalog has $TOTAL entries"
    else
        fail "Catalog has entries"
    fi
    
    # Search by category
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/backlog/search?category=mcp_server&limit=3" 2>/dev/null)
    if echo "$RESULT" | grep -q '"results"'; then
        pass "Category filter works"
    else
        fail "Category filter works"
    fi
    
    # Empty search
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/backlog/search?q=nonexistent12345&limit=5" 2>/dev/null)
    if echo "$RESULT" | grep -q '"results"'; then
        pass "Empty search handled"
    else
        fail "Empty search handled"
    fi
}

test_catalog_stats() {
    log "Testing Catalog Stats..."
    
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/backlog/stats" 2>/dev/null)
    if echo "$RESULT" | grep -q '"total"'; then
        pass "Stats endpoint works"
    else
        fail "Stats endpoint works"
    fi
    
    TOTAL=$(echo "$RESULT" | python3 -c "import sys,json;print(json.load(sys.stdin).get('total',0))" 2>/dev/null)
    if [ "$TOTAL" -gt 10000 ]; then
        pass "Stats shows $TOTAL entries"
    else
        fail "Stats shows sufficient entries (got $TOTAL)"
    fi
}

test_catalog_categories() {
    log "Testing Catalog Categories..."
    
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/backlog/categories" 2>/dev/null)
    if echo "$RESULT" | grep -q '"categories"'; then
        pass "Categories endpoint works"
    else
        fail "Categories endpoint works"
    fi
    
    COUNT=$(echo "$RESULT" | python3 -c "import sys,json;print(len(json.load(sys.stdin).get('categories',{})))" 2>/dev/null)
    if [ "$COUNT" -gt 5 ]; then
        pass "Categories has $COUNT categories"
    else
        fail "Categories has sufficient categories (got $COUNT)"
    fi
}

test_memory_system() {
    log "Testing Memory System..."
    
    RESULT=$(curl -s --max-time 5 "$BASE_URL/api/memory/stats" 2>/dev/null)
    if echo "$RESULT" | grep -q '"success"'; then
        pass "Memory stats endpoint works"
    else
        fail "Memory stats endpoint works"
    fi
}

test_account_system() {
    log "Testing Account System..."
    
    # Register
    RESULT=$(curl -s --max-time 5 -X POST "$BASE_URL/api/account/register" \
        -H "Content-Type: application/json" \
        -d '{"username":"testuser","password":"testpass123"}' 2>/dev/null)
    if echo "$RESULT" | grep -q '"success"'; then
        pass "Account registration works"
    else
        fail "Account registration works"
    fi
    
    # Login
    RESULT=$(curl -s --max-time 5 -X POST "$BASE_URL/api/account/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"testuser","password":"testpass123"}' 2>/dev/null)
    if echo "$RESULT" | grep -q '"token"'; then
        pass "Account login works"
    else
        fail "Account login works"
    fi
}

test_ssl() {
    log "Testing SSL/HTTPS..."
    
    # HTTPS redirect
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "http://demo.hypernexus.site" 2>/dev/null)
    if [ "$RESULT" = "301" ] || [ "$RESULT" = "302" ]; then
        pass "HTTP redirects to HTTPS"
    else
        fail "HTTP redirects to HTTPS (got $RESULT)"
    fi
    
    # SSL certificate
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{ssl_verify_result}" "https://demo.hypernexus.site" 2>/dev/null)
    if [ "$RESULT" = "0" ]; then
        pass "SSL certificate valid"
    else
        fail "SSL certificate valid (got $RESULT)"
    fi
}

test_headers() {
    log "Testing Security Headers..."
    
    HEADERS=$(curl -s --max-time 5 -I "$BASE_URL" 2>/dev/null)
    
    if echo "$HEADERS" | grep -qi "x-frame-options"; then
        pass "X-Frame-Options header present"
    else
        fail "X-Frame-Options header present"
    fi
    
    if echo "$HEADERS" | grep -qi "x-content-type-options"; then
        pass "X-Content-Type-Options header present"
    else
        fail "X-Content-Type-Options header present"
    fi
}

test_performance() {
    log "Testing Performance..."
    
    # Health check response time
    TIME=$(curl -s --max-time 5 -o /dev/null -w "%{time_total}" "$BASE_URL/health" 2>/dev/null)
    if (( $(echo "$TIME < 1.0" | bc -l) )); then
        pass "Health check response time: ${TIME}s"
    else
        fail "Health check response time too slow: ${TIME}s"
    fi
    
    # Search response time
    TIME=$(curl -s --max-time 5 -o /dev/null -w "%{time_total}" "$BASE_URL/api/backlog/search?q=postgres&limit=5" 2>/dev/null)
    if (( $(echo "$TIME < 2.0" | bc -l) )); then
        pass "Search response time: ${TIME}s"
    else
        fail "Search response time too slow: ${TIME}s"
    fi
}

test_landing_pages() {
    log "Testing Landing Pages..."
    
    # Main landing page
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Landing page loads"
    else
        fail "Landing page loads (got $RESULT)"
    fi
    
    # Blog
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/blog/" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Blog page loads"
    else
        fail "Blog page loads (got $RESULT)"
    fi
    
    # Catalog
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/catalog" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Catalog page loads"
    else
        fail "Catalog page loads (got $RESULT)"
    fi
    
    # Pricing
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/pricing" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Pricing page loads"
    else
        fail "Pricing page loads (got $RESULT)"
    fi
    
    # Newsletter
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/newsletter" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Newsletter page loads"
    else
        fail "Newsletter page loads (got $RESULT)"
    fi
    
    # RSS Feed
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/blog/feed.xml" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "RSS feed loads"
    else
        fail "RSS feed loads (got $RESULT)"
    fi
    
    # Sitemap
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/sitemap.xml" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Sitemap loads"
    else
        fail "Sitemap loads (got $RESULT)"
    fi
    
    # Robots.txt
    RESULT=$(curl -s --max-time 5 -o /dev/null -w "%{http_code}" "https://tormentnexus.site/robots.txt" 2>/dev/null)
    if [ "$RESULT" = "200" ]; then
        pass "Robots.txt loads"
    else
        fail "Robots.txt loads (got $RESULT)"
    fi
}

test_docker() {
    log "Testing Docker..."
    
    # Check if Docker image exists
    if docker images | grep -q "tormentnexus/tormentnexus"; then
        pass "Docker image exists"
    else
        fail "Docker image exists"
    fi
    
    # Check if demo container is running
    if docker ps | grep -q "tormentnexus-demo"; then
        pass "Demo container running"
    else
        fail "Demo container running"
    fi
}

# Run all tests
echo "=========================================="
echo "  TormentNexus Full Platform Test Suite"
echo "=========================================="
echo ""

test_health
echo ""
test_api_index
echo ""
test_catalog_search
echo ""
test_catalog_stats
echo ""
test_catalog_categories
echo ""
test_memory_system
echo ""
test_account_system
echo ""
test_ssl
echo ""
test_headers
echo ""
test_performance
echo ""
test_landing_pages
echo ""
test_docker

echo ""
echo "=========================================="
echo "  Test Results"
echo "=========================================="
echo ""
echo "  Passed: $PASS"
echo "  Failed: $FAIL"
echo "  Total:  $((PASS + FAIL))"
echo ""

if [ $FAIL -gt 0 ]; then
    echo "  ❌ SOME TESTS FAILED"
    echo ""
    echo "  Failed tests:"
    echo -e "$ERRORS"
    echo ""
    exit 1
else
    echo "  ✅ ALL TESTS PASSED"
    echo ""
    exit 0
fi
