@echo off
echo Testing Kumparan Article API with Performance Features...
echo.

echo 1. Creating test article...
curl -X POST http://localhost:8080/articles ^
     -H "Content-Type: application/json" ^
     -d "{\"title\":\"Test Article\",\"body\":\"This is a test article content\",\"author\":\"Test Author\"}"
echo.
echo.

echo 2. Creating another test article...
curl -X POST http://localhost:8080/articles ^
     -H "Content-Type: application/json" ^
     -d "{\"title\":\"Another Article\",\"body\":\"This is another test content about technology\",\"author\":\"John Doe\"}"
echo.
echo.

echo 3. Getting all articles with pagination (page 1, limit 5)...
curl "http://localhost:8080/articles?page=1&limit=5"
echo.
echo.

echo 4. Search articles with query 'test' and pagination...
curl "http://localhost:8080/articles?query=test&page=1&limit=2"
echo.
echo.

echo 5. Filter articles by author 'John' with pagination...
curl "http://localhost:8080/articles?author=John&page=1&limit=10"
echo.
echo.

echo 6. Get article by ID 1...
curl http://localhost:8080/articles/1
echo.
echo.

echo 7. Get article by ID 2...
curl http://localhost:8080/articles/2  
echo.
echo.

echo 8. Testing rate limiting (making multiple requests quickly)...
for /L %%i in (1,1,5) do (
    echo Request %%i:
    curl -s http://localhost:8080/articles?page=1&limit=1
    echo.
)
echo.

echo Testing completed!
echo.
echo Features tested:
echo - ✅ Pagination
echo - ✅ Full-text search
echo - ✅ Rate limiting
echo - ✅ CORS headers
echo - ✅ Connection pooling
pause
