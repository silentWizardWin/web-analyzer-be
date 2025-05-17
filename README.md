# Web Analyzer Backend

A REST API in Golang that analyzes a webpage and returns information.

## Project Structure

```
web-analyzer-be/
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go             # entry point
├── internal/
│   ├── handler/
│   │   └── analyze.go          # HTTP handler for /analyze
│   ├── service/
│   │   └── analyzer.go         # business logic
│   ├── util/
│   │   └── util.go             # HTML parsing helpers
│   └── model/
│       └── types.go            # request/response structs
├── README.md
```

## Run using Docker
```
# build the image
docker build -t web-analyzer .

# run the container
docker run -p 8080:8080 web-analyzer
```

## Run Locally
```
go mod tidy
go run cmd/server/main.go
```

## Run Tests & Check Coverage

Use these commands:
```
go test ./... -v
go test ./... -cover
```

To generate an HTML coverage report:
```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```
## Code Coverage Report

| Package/File          |  Coverage  |
|-----------------------|------------|
| `main.go`             |    100 %   |
| `handler/analyze.go`  |     84 %   |
| `service/analyzer.go` |     75 %   |
| `util.go`             |     87 %   |

## Endpoint

**POST /analyze**

```
curl --location 'localhost:8080/analyze'
--data '{
    "url": "https://www.gmail.com"
}'
```

### Request Payload

```json
{
  "url": "https://www.example.com"
}
```

### Response Payload
```
{
  "html_version": "HTML5",
  "title": "Example Domain",
  "headings_count": {
    "h1": 1,
    "h2": 0
  },
  "login_form_exists": false,
  "link_analysis": {
    "internal_links": 3,
    "external_links": 5,
    "inaccessible_links": 1
  }
}
```

## Further improvements
### Functional improvements
1. URL Input Enhancements
* Support multiple URLs in one request (bulk analysis).
* Use event-driven for bulk requests (poll from frontend)
* Accept URLs via query parameter.
* Allow scanning entire domains and not just a single page.

2. More HTML Analysis
Check for:
* Broken images and missing alt attributes (for accessibility).
* Analyze JavaScript usage.
* Calculate page size (KB) and resource count.

3. Authentication & History
* Add user login via JWT.
* Save analysis history per user (with timestamps).
* Provide downloadable PDF/JSON reports.

### Non-functional improvements
1. Error Handling & Resilience
* Retry with exponential backoff for temporary network failures.

2. Testing
* Use integration tests with mock web servers.
* Run coverage threshold checks in CI.

3. Security
* Sanitize input URLs to prevent SSRF.
* Add rate-limiting (10 requests/minute/IP).

4. Performance and Scalability
* Use worker queues for async analysis.
* Add caching for already-analyzed URLs.

5. CI/CD & Automation
* GitHub Actions or GitLab CI for
* Pushing Docker images to a registry.
