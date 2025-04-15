# Food Tinder

A Tinder-like backend service for voting on food products with automated product data updates from Foodji API.

## Features

- Generate unique session IDs
- Store/update product votes
- Retrieve existing votes for products
- Retrieve aggregated average scores for products across all sessions
- Automatic product data updates from Foodji API every 24 hours
- Redis caching for fast product access

## Technologies

- Go 1.23
- PostgreSQL
- Redis
- Docker
- Docker Compose

## Architecture

The application follows a clean architecture approach:

- **Domain Layer**: Core business logic and entity definitions
- **Application Layer**: Use cases and service implementations
- **Infrastructure Layer**: External interfaces like databases, HTTP handlers, and API clients
- **Presentation Layer**: HTTP API endpoints

### Key Components

- **Session Repository**: Manages user sessions in PostgreSQL
- **Vote Repository**: Handles vote storage and retrieval
- **Product Repository**: Fetches product data from Foodji API and caches it in Redis
- **Graceful Shutdown**: Properly closes all connections when the application stops

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


## API Endpoints

- `POST /api/sessions` - Create a new session
- `GET /api/sessions/{sessionID}` - Get session details
- `POST /api/sessions/{sessionID}/votes` - Create or update a vote
- `GET /api/sessions/{sessionID}/votes` - Get votes for a session
- `GET /api/votes/aggregated` - Get aggregated scores for all products

## Testing

Run the tests with:

```bash
go test ./...
```

## Folder Structure

```
.
├── cmd/
│   └── server/             # Application entry point
├── internal/
│   ├── application/        # Application services
│   ├── domain/             # Domain entities and interfaces
│   └── infrastructure/
│       ├── external/       # External API clients
│       ├── http/           # HTTP handlers and routes
│       └── persistence/    # Database repositories
├── migrations/             # SQL migration files
└── docker-compose.yml      # Docker configuration
```

## License

MIT 