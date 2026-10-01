# Event Explorer

Event Explorer is a web application for searching cities and discovering Music and Sports events. It is built with Go and Beego, with HTML, CSS, and vanilla JavaScript on the frontend.

Event data comes from the Ticketmaster Discovery API, while city search and location details are powered by the Google Places API.

## Setup

### 1. Clone the repository
```bash
git clone https://github.com
cd event-explorer
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Configure environment variables
Copy the example environment file:
```bash
cp .env.example .env
```

Add your API keys:
```env
GOOGLE_PLACES_API_KEY=your_google_places_api_key
TICKETMASTER_API_KEY=your_ticketmaster_api_key
```
Never commit the `.env` file.

### 4. Run the application
```bash
bee run
```
Then open:
```text
http://localhost:8080
```

## Routes

### Pages

| Method | Route                                 | Description               |
| :---   | :---                                  | :---                      |
| GET    | `/`                                   | Home page and city search |
| GET    | `/events?city=Toronto&countryCode=CA` | Event listing             |
| GET    | `/events/:eventId`                    | Event details             |

### API

| Method | Route                                                    | Description                         |
| :---   | :---                                                     | :---                                |
| GET    | `/api/locations/autocomplete?input=Tor&sessionToken=...` | City autocomplete suggestions       |
| GET    | `/api/locations/:placeId?sessionToken=...`               | Selected city details               |
| POST   | `/api/cache/clear`                                       | Clear all or specific cached events |

### Cache Management
The cache endpoint supports clearing the entire cache or clearing a specific cache entry dynamically via POST requests.

Clear all cached events:
```text
POST /api/cache/clear
```

Clear a specific city's category cache:
```text
POST /api/cache/clear?city=Toronto&countryCode=CA&category=Music
```

Supported categories:
* `Music`
* `Sports`

The `city`, `countryCode`, and `category` parameters must either all be provided together or all be omitted.

## Testing
Unit tests use mocked API responses, so live Google Places or Ticketmaster API keys are not required for testing.

Run all tests:
```bash
go test ./...
```

Run tests with race detection:
```bash
go test -race ./...
```

Show package-level coverage:
```bash
go test ./... -cover
```

Generate a coverage profile:
```bash
go test ./... -coverprofile=coverage.out
```

Show detailed coverage:
```bash
go tool cover -func=coverage.out
```

Open the coverage report in a browser:
```bash
go tool cover -html=coverage.out
```

Current test coverage: **91.3%**

## Project Structure
```text
event-explorer/
├── controllers/
├── models/
├── routers/
├── services/
├── validators/
├── tests/
├── static/
├── views/
├── .env.example
├── go.mod
└── main.go
```

## Environment Variables

| Variable                | Description                            |
| :---                    | :---                                   |
| `GOOGLE_PLACES_API_KEY` | API key for Google Places API          |
| `TICKETMASTER_API_KEY`  | API key for Ticketmaster Discovery API |
