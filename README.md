# Service Customer

Customer profile, addresses, wishlist, and order history management service for the niaga-platform e-commerce system.

## Features

- **Customer Profiles**: Manage customer information (name, email, phone, etc.)
- **Address Management**: CRUD operations for shipping/billing addresses with default address support
- **Wishlist**: Save and manage favorite products
- **Order History**: View customer order history (placeholder for cross-service integration)

## Tech Stack

- **Go 1.23.0**
- **Gin** - HTTP web framework
- **GORM** - ORM for PostgreSQL
- **PostgreSQL** - Database with `customer` schema
- **JWT** - Authentication middleware

## Project Structure

```
service-customer/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loader
│   ├── handlers/
│   │   ├── profile_handler.go   # Profile endpoints
│   │   ├── address_handler.go   # Address endpoints
│   │   ├── wishlist_handler.go  # Wishlist endpoints
│   │   └── order_history_handler.go
│   ├── middleware/
│   │   └── auth.go              # JWT authentication
│   ├── models/
│   │   ├── profile.go           # Customer profile model
│   │   ├── address.go           # Address model
│   │   └── wishlist.go          # Wishlist item model
│   └── repository/
│       ├── profile_repository.go
│       ├── address_repository.go
│       └── wishlist_repository.go
├── .env.example
├── go.mod
└── README.md
```

## API Endpoints

### Profile Management
- `GET /api/v1/customer/profile` - Get customer profile
- `PUT /api/v1/customer/profile` - Update customer profile

### Address Management
- `GET /api/v1/customer/addresses` - List all addresses
- `POST /api/v1/customer/addresses` - Create new address
- `PUT /api/v1/customer/addresses/:id` - Update address
- `DELETE /api/v1/customer/addresses/:id` - Delete address
- `PUT /api/v1/customer/addresses/:id/default` - Set default address

### Wishlist
- `GET /api/v1/customer/wishlist` - Get wishlist items
- `POST /api/v1/customer/wishlist` - Add product to wishlist
- `DELETE /api/v1/customer/wishlist/:productId` - Remove from wishlist

### Order History
- `GET /api/v1/customer/orders` - Get order history

All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

## Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Configure database settings in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=customer_db
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the service:
```bash
go run cmd/server/main.go
```

The service will start on port 8004 (configurable via `APP_PORT`).

## Database

The service uses a PostgreSQL database with a `customer` schema:

- `customer.profiles` - Customer profile information
- `customer.addresses` - Shipping/billing addresses
- `customer.wishlist_items` - Saved products

Schema is automatically created and migrated on startup.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| APP_PORT | 8004 | Server port |
| APP_ENV | development | Environment (development/production) |
| DB_HOST | localhost | Database host |
| DB_PORT | 5432 | Database port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | customer_db | Database name |
| JWT_SECRET | - | JWT signing secret |

## Health Check

```bash
curl http://localhost:8004/health
```

## Development

### Adding New Features

1. Define model in `internal/models/`
2. Create repository in `internal/repository/`
3. Implement handler in `internal/handlers/`
4. Register routes in `cmd/server/main.go`

### TODO

- [ ] Implement unit tests for repositories
- [ ] Implement integration tests for handlers
- [ ] Add cross-service communication with service-order
- [ ] Add NATS event publishing for profile/address changes
- [ ] Add Redis caching for frequently accessed profiles
- [ ] Add profile picture upload with MinIO integration

## License

Part of the niaga-platform e-commerce system.
