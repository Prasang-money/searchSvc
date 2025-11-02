# Country Search Service (searchSvc)

A REST API service that provides country information with caching capabilities. The service fetches data from the [REST Countries API](https://restcountries.com/) and implements an LRU cache to optimize repeated queries.

## Features

- Search countries by name
- LRU (Least Recently Used) caching mechanism (not handling collision of key)
- Thread-safe implementation
- RESTful API endpoints

## Setup Instructions

### Prerequisites

- Go 1.25 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Prasang-money/searchSvc.git
cd searchSvc
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build
```

4. Run the service:
```bash
./searchSvc
```

## API Documentation

### Endpoints

#### 1. Health Check
Check if the service is running.

```
GET /health
```

Response:
```json
{
    "status": "OK"
}
```

#### 2. Search Countries
Search for country information by name.

```
GET /search?name={countryName}
```

Parameters:
- `name` (required): The name of the country to search for

Example Request:
```
GET /api/countries/search?name=United States
```

Example Response:
```json
{
    "name": "United States",
    "population": 331002651,
    "capital": "Washington, D.C.",
    "currency": "$"
}
```

Error Response (500 Internal Server Error):
```json
{
    "error": "failed to fetch country data"
}
```

## Project Structure

```
searchSvc/
├── cache/          # LRU cache implementation
├── handler/        # HTTP handlers
├── models/         # Data models
├── route/          # Router configuration
├── service/        # Business logic
└── utils/          # Utility functions
```

## Architecture

- **Cache Layer**: Implements an LRU (Least Recently Used) caching mechanism using a combination of a hash map and doubly linked list.
- **Service Layer**: Handles business logic, external API calls, and cache interactions.
- **Handler Layer**: Manages HTTP request/response handling.
- **Model Layer**: Defines data structures used throughout the application.

## Configuration

The service uses the following default configurations:

- Default port: 8080
- Cache capacity: Configurable via initialization
- External API: REST Countries API (https://restcountries.com/v3.1)
- HTTP client timeout: 10 seconds

## Development

### Running Tests

To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test ./... -cover
```

To run tests for a specific package:
```bash
go test ./cache  # Test cache package
go test ./service # Test service package
go test ./handler # Test handler package
```

### Adding New Features

1. Create necessary model structs in `models/`
2. Implement business logic in `service/`
3. Add HTTP handlers in `handler/`
4. Update routes in `route/`
5. Add unit tests for new functionality

## Error Handling

The service implements the following error handling:

- HTTP 500: Internal Server Error (API failures, parsing errors)
- HTTP 404: Country not found
- HTTP 200: Successful response
- Cache misses are handled gracefully by fetching from the external API

## Performance Considerations

- LRU cache implementation provides O(1) complexity for both read and write operations
- Thread-safe implementation ensures concurrent access safety
- Configurable cache size to balance memory usage and hit ratio
- HTTP connection pooling for external API calls
