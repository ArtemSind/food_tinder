# Food Tinder

A Tinder-like backend service for voting on food products.

## Features

- Generate unique session IDs
- Store/update product votes
- Retrieve existing votes for products
- Retrieve aggregated average scores for products across all sessions

## Technologies

- Go 1.19
- PostgreSQL
- Docker
- Docker Compose

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your machine

### Running with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/ArtemSind/food_tinder.git
cd food_tinder
```

2. Start the application:
```bash
docker-compose up -d
```

The API will be available at http://localhost:8080

### Development Setup

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update the values in `.env` if needed

3. Start a PostgreSQL instance:
```bash
docker-compose up -d postgres
```

4. Run the application locally:
```bash
go run cmd/server/main.go
```

## API Endpoints

- `POST /api/sessions` - Create a new session
- `GET /api/sessions/{sessionID}` - Get session details
- `POST /api/sessions/{sessionID}/votes` - Create or update a vote
- `GET /api/sessions/{sessionID}/votes` - Get votes for a session
- `GET /api/votes/aggregated` - Get aggregated scores for all products

## License

MIT 