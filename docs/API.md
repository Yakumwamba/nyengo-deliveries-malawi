# Nyengo Deliveries API Documentation

## Base URL
```
Production: https://api.nyengo.com/api/v1
Development: http://localhost:8080/api/v1
```

## Authentication
Most endpoints require a JWT Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```

## Courier Endpoints

### Register Courier
```http
POST /couriers/register
Content-Type: application/json

{
  "email": "courier@example.com",
  "password": "securepassword",
  "companyName": "Fast Deliveries",
  "ownerName": "John Doe",
  "phone": "+260970000000",
  "address": "123 Cairo Road",
  "city": "Lusaka",
  "country": "Zambia",
  "serviceAreas": ["Lusaka", "Kitwe"],
  "vehicleTypes": ["motorcycle", "car"]
}
```

### Login
```http
POST /couriers/login
Content-Type: application/json

{
  "email": "courier@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "jwt-token-here",
    "courier": { ... }
  }
}
```

### Get Profile
```http
GET /couriers/profile
Authorization: Bearer <token>
```

### Update Profile
```http
PUT /couriers/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "companyName": "Updated Name",
  "baseRatePerKm": 6.0,
  "minimumFare": 25.0
}
```

## Pricing Endpoints

### Get Price Estimate
```http
POST /pricing/estimate
Content-Type: application/json

{
  "pickupLatitude": -15.4167,
  "pickupLongitude": 28.2833,
  "deliveryLatitude": -15.4101,
  "deliveryLongitude": 28.3122,
  "packageSize": "medium",
  "packageWeight": 5.0,
  "isFragile": false,
  "isExpress": false
}
```

Response:
```json
{
  "success": true,
  "data": {
    "currency": "ZMW",
    "currencySymbol": "K",
    "distance": 5.2,
    "duration": 15,
    "baseFare": 20.0,
    "distanceFare": 26.0,
    "totalFare": 50.6,
    "formattedTotal": "K50.60"
  }
}
```

## Order Endpoints

### Create Order
```http
POST /orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "customerName": "Jane Smith",
  "customerPhone": "+260971234567",
  "pickupAddress": "Cairo Road, Shop 45",
  "pickupLatitude": -15.4167,
  "pickupLongitude": 28.2833,
  "deliveryAddress": "23 Independence Ave",
  "deliveryLatitude": -15.4101,
  "deliveryLongitude": 28.3122,
  "packageDescription": "Electronics",
  "packageSize": "medium",
  "paymentMethod": "cash"
}
```

### List Orders
```http
GET /orders?page=1&pageSize=20&status=pending,in_transit
Authorization: Bearer <token>
```

### Get Order Details
```http
GET /orders/{id}
Authorization: Bearer <token>
```

### Update Order Status
```http
PUT /orders/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "in_transit",
  "note": "Picked up package"
}
```

### Accept/Decline Order
```http
PUT /orders/{id}/accept
PUT /orders/{id}/decline
Authorization: Bearer <token>
```

## Store Integration Endpoints

### List Available Couriers
```http
GET /stores/couriers?area=Lusaka
X-API-Key: <store-api-key>
```

### Create Order (from Store)
```http
POST /stores/orders
X-API-Key: <store-api-key>
Content-Type: application/json

{
  "courierId": "uuid",
  "customerName": "...",
  ...
}
```

### Track Order Status
```http
GET /stores/orders/{id}/status
X-API-Key: <store-api-key>
```

## WebSocket Connection
```
ws://localhost:8080/ws?token=<jwt-token>
```

Events:
- `new_order` - New order received
- `order_update` - Order status changed
- `chat_message` - New chat message

## Error Responses
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```
