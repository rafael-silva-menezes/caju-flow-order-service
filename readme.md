Here's the updated documentation for your **Order Service**:

---

# Order Service

The **Order Service** is responsible for creating, querying, updating, and canceling orders, and sending them for processing. It exposes a REST API to interact with clients and publishes orders to a RabbitMQ queue to be processed by another service.

## Features

- **Create Order**: Endpoint to create a new order and publish it to the RabbitMQ queue.
- **Query Order**: Endpoint to query an existing order by its ID.
- **List Orders**: Endpoint to list orders with pagination.
- **Update Order**: Endpoint to update an order when it is in "pending" status.
- **Cancel Orders**: Endpoint to cancel an order when it is in "pending" status.

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

#### Response:

```json
{
  "order_id": "12345",
  "status": "pending"
}
```

### `GET /orders/{id}`

Queries an order by ID.

#### Response:

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

#### Response:

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

### `PUT /orders/{id}`

Updates an order, provided that its status is "pending".

#### Request Body:

```json
{
  "customer_name": "John Doe",
  "items": [
    { "name": "Product A", "quantity": 3, "price": 45.0 }
  ]
}
```

#### Response:

```json
{
  "order_id": "12345",
  "status": "pending"
}
```

### `DELETE /orders/{id}`

Cancels an order, provided that its status is "pending".

#### Response:

```json
{
  "order_id": "12345",
  "status": "canceled"
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

---

