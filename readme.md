# Order Service

The **Order Service** is responsible for creating, querying, and sending orders for processing. It exposes a REST API to interact with clients and publishes orders to a RabbitMQ queue to be processed by another service.

## Features

- **Create Order**: Endpoint to create a new order.
- **Query Order**: Endpoint to query an existing order by its ID.
- **List Orders**: Endpoint to list orders with pagination.

## Endpoints

### `POST /orders`

Creates a new order and publishes it to the RabbitMQ queue.

#### Request Body:

```json
{
  "customer_name": "John Smith",
  "items": [
    { "name": "Product A", "quantity": 2, "price": 50.0 },
    { "name": "Product B", "quantity": 1, "price": 30.0 }
  ]
}
```

#### Create Order Response:

```json
{
  "order_id": "12345",
  "status": "pending"
}
```

### `GET /orders/{id}`

Queries an order by ID.

#### Query Order Response:

```json
{
  "order_id": "12345",
  "customer_name": "John Smith",
  "items": [
    { "name": "Product A", "quantity": 2, "price": 50.0 },
    { "name": "Product B", "quantity": 1, "price": 30.0 }
  ],
  "status": "pending"
}
```

### `GET /orders?page={page}&size={size}`

Lists orders with pagination.

#### List Orders Response:

```json
{
  "page": 1,
  "size": 10,
  "total": 50,
  "orders": [
    { "order_id": "12345", "customer_name": "John Smith", "status": "pending" },
    { "order_id": "12346", "customer_name": "Mary Sanders", "status": "completed" }
  ]
}
```

## Execution Instructions

### Environment Configuration

Configure environment variables for RabbitMQ and PostgreSQL.

### Build and Run the Service:

```bash
docker build -t order-service .
docker run -p 8080:8080 order-service
```

### Tests

- Use TDD to ensure test coverage.
- Perform unit and integration tests.