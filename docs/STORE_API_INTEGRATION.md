# Nyengo Deliveries - Store Integration API Documentation

This document provides complete guidance for store application developers who want to integrate with the Nyengo Deliveries courier backend. It covers how to retrieve available couriers, create delivery orders, and implement real-time tracking.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Before Checkout: Courier Selection](#before-checkout-courier-selection)
- [At Checkout: Creating Orders](#at-checkout-creating-orders)
- [After Checkout: Order Tracking](#after-checkout-order-tracking)
- [WebSocket Integration](#websocket-integration)
- [Webhook Notifications](#webhook-notifications)
- [Error Handling](#error-handling)
- [Code Examples](#code-examples)

---

## Overview

The Nyengo Deliveries API provides a complete solution for e-commerce stores to:

1. **Pre-Checkout**: Display available couriers with pricing estimates based on delivery distance
2. **Checkout**: Create delivery orders with selected couriers
3. **Post-Checkout**: Track deliveries in real-time and receive status updates

### Base URL

```
Production: https://api.nyengo-deliveries.com/api/v1
Development: http://localhost:8082/api/v1
```

### Response Format

All responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "message": "Optional message"
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error description",
  "code": "ERROR_CODE"
}
```

---

## Authentication

Store integrations use **API Key authentication**.

### Test API Key (Development)

For development and testing, use the following pre-configured API key:

```
nyg_test_store_api_key_2024_dev
```

> ⚠️ **Important**: This key is for testing only. Contact the Nyengo Deliveries team for a production API key.

### Header Format

```http
X-API-Key: nyg_test_store_api_key_2024_dev
```

All store endpoints (`/stores/*`) require this header. Requests without a valid API key will receive a `401 Unauthorized` response.

### Environment Configuration (Backend)

API keys are configured via the `STORE_API_KEYS` environment variable (comma-separated for multiple keys):

```bash
# .env
STORE_API_KEYS=your-production-key-1,your-production-key-2
```

---

## Before Checkout: Courier Selection

When a customer enters their delivery address, use these endpoints to display available couriers and pricing.

### 1. Get Price Estimate

Get a quick price estimate before showing couriers.

**Endpoint**: `POST /pricing/estimate`

**Note**: This endpoint is public (no authentication required).

#### Request Body

```json
{
  "pickupLatitude": -15.4167,
  "pickupLongitude": 28.2833,
  "pickupAddress": "Cairo Road, Lusaka",
  "deliveryLatitude": -15.3875,
  "deliveryLongitude": 28.3228,
  "deliveryAddress": "Kabulonga, Lusaka",
  "packageSize": "medium",
  "packageWeight": 2.5,
  "isFragile": false,
  "isExpress": false
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pickupLatitude` | float | ✅ | Pickup location latitude |
| `pickupLongitude` | float | ✅ | Pickup location longitude |
| `deliveryLatitude` | float | ✅ | Delivery location latitude |
| `deliveryLongitude` | float | ✅ | Delivery location longitude |
| `pickupAddress` | string | ❌ | Human-readable pickup address |
| `deliveryAddress` | string | ❌ | Human-readable delivery address |
| `packageSize` | string | ❌ | `small`, `medium`, or `large` |
| `packageWeight` | float | ❌ | Weight in kilograms |
| `isFragile` | boolean | ❌ | Requires fragile handling |
| `isExpress` | boolean | ❌ | Express delivery requested |

#### Response

```json
{
  "success": true,
  "data": {
    "currency": "ZMW",
    "currencySymbol": "K",
    "distance": 8.5,
    "duration": 25,
    "baseFare": 15.00,
    "distanceFare": 42.50,
    "weightFare": 5.00,
    "fragileFare": 0,
    "expressFare": 0,
    "surgeFare": 0,
    "subTotal": 62.50,
    "platformFee": 6.25,
    "totalFare": 68.75,
    "formattedTotal": "K 68.75",
    "formattedBreakdown": {
      "baseFare": "K 15.00",
      "distanceFare": "K 42.50",
      "subTotal": "K 62.50",
      "platformFee": "K 6.25",
      "totalFare": "K 68.75"
    },
    "estimatedPickup": "10-15 mins",
    "estimatedDelivery": "35-45 mins",
    "isLocalDelivery": true,
    "deliveryType": "local",
    "disclaimer": "Prices may vary based on traffic and conditions"
  }
}
```

---

### 2. List Available Couriers

Get available couriers based on delivery distance. The API automatically determines whether to return local couriers (for nearby deliveries) or external courier services (for long-distance/inter-city deliveries).

**Endpoint**: `GET /stores/couriers`

**Authentication**: Required (API Key)

#### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pickupLat` | float | ✅ | Pickup location latitude |
| `pickupLon` | float | ✅ | Pickup location longitude |
| `deliveryLat` | float | ✅ | Delivery location latitude |
| `deliveryLon` | float | ✅ | Delivery location longitude |
| `area` | string | ❌ | Fallback: filter by service area name |

#### Example Request

```bash
curl -X GET "https://api.nyengo-deliveries.com/api/v1/stores/couriers?\
pickupLat=-15.4167&pickupLon=28.2833&\
deliveryLat=-15.3875&deliveryLon=28.3228" \
  -H "X-API-Key: your-api-key"
```

#### Response (Local Delivery < 30km)

```json
{
  "success": true,
  "data": {
    "deliveryType": "local",
    "isLocalDelivery": true,
    "distance": 8.5,
    "threshold": 30,
    "couriers": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "local",
        "name": "Swift Riders Zambia",
        "logoUrl": "https://example.com/logo.png",
        "estimatedFare": 65.00,
        "formattedFare": "K 65.00",
        "baseRatePerKm": 5.00,
        "minimumFare": 20.00,
        "rating": 4.8,
        "totalReviews": 156,
        "totalDeliveries": 1250,
        "isVerified": true,
        "isFeatured": true,
        "estimatedTime": "20-35 mins"
      },
      {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
        "type": "local",
        "name": "QuickDrop Express",
        "logoUrl": "https://example.com/quickdrop.png",
        "estimatedFare": 72.50,
        "formattedFare": "K 72.50",
        "baseRatePerKm": 6.50,
        "minimumFare": 25.00,
        "rating": 4.5,
        "totalReviews": 89,
        "totalDeliveries": 680,
        "isVerified": true,
        "isFeatured": false,
        "estimatedTime": "25-40 mins"
      }
    ],
    "recommended": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Swift Riders Zambia",
      "recommendedFor": "recommended"
    },
    "cheapestOption": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Swift Riders Zambia",
      "recommendedFor": "cheapest"
    },
    "fastestOption": {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "name": "QuickDrop Express",
      "recommendedFor": "best_rated"
    }
  }
}
```

#### Response (Inter-City Delivery >= 30km)

```json
{
  "success": true,
  "data": {
    "deliveryType": "intercity",
    "isLocalDelivery": false,
    "distance": 320.5,
    "threshold": 30,
    "couriers": [
      {
        "id": "dhl-express",
        "type": "external",
        "name": "DHL Express",
        "logoUrl": "https://example.com/dhl-logo.png",
        "description": "Fast international courier service",
        "estimatedFare": 450.00,
        "formattedFare": "K 450.00",
        "baseRatePerKm": 1.50,
        "minimumFare": 100.00,
        "estimatedDeliveryDays": "1-2 business days",
        "serviceType": "express",
        "trackingUrl": "https://www.dhl.com/track"
      },
      {
        "id": "fedex-economy",
        "type": "external",
        "name": "FedEx Economy",
        "logoUrl": "https://example.com/fedex-logo.png",
        "description": "Reliable economy shipping",
        "estimatedFare": 280.00,
        "formattedFare": "K 280.00",
        "baseRatePerKm": 0.80,
        "minimumFare": 75.00,
        "estimatedDeliveryDays": "3-5 business days",
        "serviceType": "economy"
      }
    ],
    "recommended": {
      "id": "dhl-express",
      "name": "DHL Express",
      "recommendedFor": "recommended"
    },
    "cheapestOption": {
      "id": "fedex-economy",
      "name": "FedEx Economy",
      "recommendedFor": "cheapest"
    },
    "fastestOption": {
      "id": "dhl-express",
      "name": "DHL Express",
      "recommendedFor": "fastest"
    }
  }
}
```

### Understanding Courier Types

| Type | When Returned | Characteristics |
|------|---------------|-----------------|
| `local` | Distance < 30km | Registered local couriers with ratings, same-day delivery |
| `external` | Distance >= 30km | External services (DHL, FedEx), multi-day delivery with external tracking |

---

## At Checkout: Creating Orders

Once the customer selects a courier and confirms their order, create a delivery request.

### Create Order

**Endpoint**: `POST /stores/orders`

**Authentication**: Required (API Key)

#### Request Body

```json
{
  "courierId": "550e8400-e29b-41d4-a716-446655440000",
  "customerName": "John Banda",
  "customerPhone": "+260977123456",
  "customerEmail": "john@example.com",
  "pickupAddress": "ShopRite, Cairo Road, Lusaka",
  "pickupLatitude": -15.4167,
  "pickupLongitude": 28.2833,
  "pickupNotes": "Collect from customer service desk",
  "pickupContactName": "Store Manager",
  "pickupContactPhone": "+260211234567",
  "deliveryAddress": "Plot 123, Kabulonga, Lusaka",
  "deliveryLatitude": -15.3875,
  "deliveryLongitude": 28.3228,
  "deliveryNotes": "Gate has intercom - buzz unit 5",
  "packageDescription": "Electronics - Laptop",
  "packageSize": "medium",
  "packageWeight": 3.5,
  "isFragile": true,
  "requiresSignature": true,
  "paymentMethod": "card",
  "scheduledPickup": "2024-12-24T14:00:00Z",
  "externalOrderId": "STORE-ORD-12345",
  "storeId": "optional-store-uuid"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `courierId` | string | ✅ | UUID of selected courier from listing |
| `customerName` | string | ✅ | Customer's full name |
| `customerPhone` | string | ✅ | Customer's phone number |
| `customerEmail` | string | ❌ | Customer's email address |
| `pickupAddress` | string | ✅ | Full pickup address |
| `pickupLatitude` | float | ✅ | Pickup GPS latitude |
| `pickupLongitude` | float | ✅ | Pickup GPS longitude |
| `pickupNotes` | string | ❌ | Special pickup instructions |
| `pickupContactName` | string | ❌ | Person to contact at pickup |
| `pickupContactPhone` | string | ❌ | Contact phone at pickup |
| `deliveryAddress` | string | ✅ | Full delivery address |
| `deliveryLatitude` | float | ✅ | Delivery GPS latitude |
| `deliveryLongitude` | float | ✅ | Delivery GPS longitude |
| `deliveryNotes` | string | ❌ | Delivery instructions |
| `packageDescription` | string | ✅ | What's being delivered |
| `packageSize` | string | ✅ | `small`, `medium`, or `large` |
| `packageWeight` | float | ❌ | Weight in kg |
| `isFragile` | boolean | ❌ | Fragile handling required |
| `requiresSignature` | boolean | ❌ | Signature required on delivery |
| `paymentMethod` | string | ✅ | `cash`, `mobile_money`, `card`, `wallet` |
| `scheduledPickup` | datetime | ❌ | ISO 8601 format, for scheduled pickups |
| `externalOrderId` | string | ❌ | Your store's order ID for reference |
| `storeId` | UUID | ❌ | Your store's ID in our system |

#### Response

```json
{
  "success": true,
  "data": {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "orderNumber": "NYG-20241224-A1B2C3",
    "courierId": "550e8400-e29b-41d4-a716-446655440000",
    "externalOrderId": "STORE-ORD-12345",
    "customerName": "John Banda",
    "customerPhone": "+260977123456",
    "status": "pending",
    "paymentStatus": "pending",
    "distance": 8.5,
    "totalFare": 85.00,
    "platformFee": 8.50,
    "courierEarnings": 76.50,
    "estimatedDelivery": "2024-12-24T15:30:00Z",
    "createdAt": "2024-12-24T12:00:00Z"
  }
}
```

> **Important**: Store the `id` (order ID) returned in the response. You'll need this for tracking and status updates.

---

## After Checkout: Order Tracking

After order creation, implement tracking to keep customers informed about their delivery.

### 1. Get Order Status

Poll this endpoint to get the current order status.

**Endpoint**: `GET /stores/orders/:id/status`

**Authentication**: Required (API Key)

#### Response

```json
{
  "success": true,
  "data": {
    "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "orderNumber": "NYG-20241224-A1B2C3",
    "status": "in_transit",
    "updatedAt": "2024-12-24T13:45:00Z"
  }
}
```

#### Order Status Values

| Status | Description | Customer Message |
|--------|-------------|-----------------|
| `pending` | Waiting for courier to accept | "Looking for a courier..." |
| `accepted` | Courier accepted the order | "Courier assigned!" |
| `declined` | Courier declined (will be reassigned) | "Finding another courier..." |
| `picked_up` | Package collected from pickup location | "Package picked up!" |
| `in_transit` | En route to delivery address | "On the way!" |
| `delivered` | Successfully delivered | "Delivered!" |
| `cancelled` | Order was cancelled | "Order cancelled" |
| `failed` | Delivery failed | "Delivery unsuccessful" |

---

### 2. Get Live Tracking

Get real-time location and ETA for active deliveries.

**Endpoint**: `GET /tracking/:orderId`

**Authentication**: Not required (tracking is public for customer convenience)

#### Response

```json
{
  "success": true,
  "data": {
    "id": "tracking-uuid",
    "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "driverName": "Moses Phiri",
    "driverPhone": "+260966789012",
    "vehicleType": "motorcycle",
    "vehiclePlate": "ALU 1234",
    "currentLatitude": -15.4012,
    "currentLongitude": 28.3021,
    "lastLocationAt": "2024-12-24T13:45:30Z",
    "estimatedArrival": "2024-12-24T14:05:00Z",
    "distanceRemaining": 3.2,
    "durationRemaining": 1200,
    "isActive": true
  }
}
```

| Field | Description |
|-------|-------------|
| `currentLatitude/Longitude` | Driver's current GPS coordinates |
| `lastLocationAt` | When location was last updated |
| `estimatedArrival` | Predicted arrival time (ISO 8601) |
| `distanceRemaining` | Kilometers to destination |
| `durationRemaining` | Seconds until arrival |
| `isActive` | Whether tracking is currently active |

---

### 3. Get Location History

Retrieve the path the driver has taken (useful for showing the route on a map).

**Endpoint**: `GET /tracking/:orderId/history`

**Authentication**: Not required

#### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 100 | Maximum points to return (max 1000) |

#### Response

```json
{
  "success": true,
  "data": {
    "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "points": [
      {
        "latitude": -15.4167,
        "longitude": 28.2833,
        "timestamp": "2024-12-24T13:30:00Z",
        "speed": 0,
        "heading": 45
      },
      {
        "latitude": -15.4102,
        "longitude": 28.2901,
        "timestamp": "2024-12-24T13:35:00Z",
        "speed": 35.5,
        "heading": 67
      }
    ],
    "count": 2
  }
}
```

---

## WebSocket Integration

For real-time updates without polling, connect to our WebSocket endpoint.

### Connection

```
wss://api.nyengo-deliveries.com/ws?token=your-jwt-token
```

> **Note**: WebSocket connections require JWT authentication. Contact support for WebSocket access tokens.

### Message Format

#### Subscribe to Order Updates

```json
{
  "type": "subscribe",
  "channel": "order:7c9e6679-7425-40de-944b-e07fc1f90ae7"
}
```

#### Location Update Event

```json
{
  "type": "location_update",
  "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "data": {
    "latitude": -15.4012,
    "longitude": 28.3021,
    "timestamp": "2024-12-24T13:45:30Z",
    "distanceRemaining": 3.2,
    "etaMinutes": 12
  }
}
```

#### Status Change Event

```json
{
  "type": "status_update",
  "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "data": {
    "status": "picked_up",
    "timestamp": "2024-12-24T13:30:00Z",
    "note": "Package collected successfully"
  }
}
```

---

## Webhook Notifications

Configure webhooks to receive push notifications for order events.

### Setup

Contact the Nyengo Deliveries team to configure your webhook URL:

- Webhook URL: `https://your-store.com/webhooks/nyengo`
- Events to subscribe: `order.accepted`, `order.picked_up`, `order.delivered`, etc.

### Webhook Payload

```json
{
  "event": "order.status_changed",
  "timestamp": "2024-12-24T13:45:00Z",
  "data": {
    "orderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "orderNumber": "NYG-20241224-A1B2C3",
    "externalOrderId": "STORE-ORD-12345",
    "previousStatus": "accepted",
    "newStatus": "picked_up",
    "courierName": "Swift Riders Zambia"
  }
}
```

### Webhook Events

| Event | Description |
|-------|-------------|
| `order.accepted` | Courier accepted the delivery |
| `order.picked_up` | Package collected from pickup |
| `order.in_transit` | Delivery is on the way |
| `order.delivered` | Successfully delivered |
| `order.failed` | Delivery attempt failed |
| `order.cancelled` | Order was cancelled |

---

## Error Handling

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| `200` | Success |
| `201` | Created (order created successfully) |
| `400` | Bad Request (invalid parameters) |
| `401` | Unauthorized (invalid/missing API key) |
| `404` | Not Found (order/courier not found) |
| `429` | Rate Limited (slow down requests) |
| `500` | Server Error |

### Error Response Format

```json
{
  "success": false,
  "error": "Invalid courier ID",
  "code": "INVALID_COURIER_ID"
}
```

### Common Errors

| Error Code | Cause | Solution |
|------------|-------|----------|
| `INVALID_API_KEY` | Missing or invalid API key | Check your API key header |
| `INVALID_COURIER_ID` | Courier UUID is malformed | Use the ID from courier listing |
| `COURIER_NOT_FOUND` | Courier no longer available | Refresh courier list |
| `INVALID_COORDINATES` | GPS coordinates are invalid | Verify lat/long values |
| `ORDER_NOT_FOUND` | Order ID doesn't exist | Check the order ID |

---

## Code Examples

### JavaScript/TypeScript

```typescript
// api-client.ts
const API_BASE = 'https://api.nyengo-deliveries.com/api/v1';
const API_KEY = 'your-api-key';

// Get available couriers
async function getCouriers(pickup: {lat: number, lng: number}, delivery: {lat: number, lng: number}) {
  const params = new URLSearchParams({
    pickupLat: pickup.lat.toString(),
    pickupLon: pickup.lng.toString(),
    deliveryLat: delivery.lat.toString(),
    deliveryLon: delivery.lng.toString(),
  });
  
  const response = await fetch(`${API_BASE}/stores/couriers?${params}`, {
    headers: { 'X-API-Key': API_KEY }
  });
  
  return response.json();
}

// Create order
async function createOrder(orderData: OrderRequest) {
  const response = await fetch(`${API_BASE}/stores/orders`, {
    method: 'POST',
    headers: {
      'X-API-Key': API_KEY,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(orderData)
  });
  
  return response.json();
}

// Get live tracking
async function getTracking(orderId: string) {
  const response = await fetch(`${API_BASE}/tracking/${orderId}`);
  return response.json();
}
```

### Python

```python
import requests

API_BASE = "https://api.nyengo-deliveries.com/api/v1"
API_KEY = "your-api-key"

headers = {"X-API-Key": API_KEY}

# Get available couriers
def get_couriers(pickup_lat, pickup_lon, delivery_lat, delivery_lon):
    params = {
        "pickupLat": pickup_lat,
        "pickupLon": pickup_lon,
        "deliveryLat": delivery_lat,
        "deliveryLon": delivery_lon
    }
    response = requests.get(f"{API_BASE}/stores/couriers", 
                           headers=headers, params=params)
    return response.json()

# Create order
def create_order(order_data):
    response = requests.post(f"{API_BASE}/stores/orders",
                            headers=headers, json=order_data)
    return response.json()

# Get live tracking
def get_tracking(order_id):
    response = requests.get(f"{API_BASE}/tracking/{order_id}")
    return response.json()
```

---

## Integration Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CHECKOUT FLOW                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐    ┌──────────────────┐    ┌────────────────────────────┐  │
│  │ Customer    │───▶│ GET /stores/     │───▶│ Display courier options    │  │
│  │ enters      │    │ couriers         │    │ with prices & delivery     │  │
│  │ address     │    │                  │    │ times                      │  │
│  └─────────────┘    └──────────────────┘    └────────────────────────────┘  │
│         │                                              │                     │
│         │                                              ▼                     │
│         │                                   ┌────────────────────────────┐  │
│         │                                   │ Customer selects courier   │  │
│         │                                   └────────────────────────────┘  │
│         │                                              │                     │
│         ▼                                              ▼                     │
│  ┌─────────────┐    ┌──────────────────┐    ┌────────────────────────────┐  │
│  │ Customer    │───▶│ POST /stores/    │───▶│ Store Order ID for        │  │
│  │ confirms    │    │ orders           │    │ tracking                   │  │
│  │ checkout    │    │                  │    │                            │  │
│  └─────────────┘    └──────────────────┘    └────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                          POST-CHECKOUT TRACKING                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────┐     ┌────────────────┐     ┌──────────────────────────┐ │
│  │ GET /stores/   │────▶│ Display status │────▶│  pending → accepted →   │ │
│  │ orders/:id/    │     │ to customer    │     │  picked_up → in_transit │ │
│  │ status         │     │                │     │  → delivered            │ │
│  └────────────────┘     └────────────────┘     └──────────────────────────┘ │
│          │                                                                   │
│          ▼                                                                   │
│  ┌────────────────┐     ┌────────────────┐     ┌──────────────────────────┐ │
│  │ GET /tracking/ │────▶│ Show live map  │────▶│  Driver location + ETA  │ │
│  │ :orderId       │     │ to customer    │     │                          │ │
│  └────────────────┘     └────────────────┘     └──────────────────────────┘ │
│          │                                                                   │
│          ▼                                                                   │
│  ┌────────────────┐     ┌────────────────┐                                  │
│  │ WebSocket or   │────▶│ Real-time      │                                  │
│  │ Webhook        │     │ push updates   │                                  │
│  └────────────────┘     └────────────────┘                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Support

For integration support or to obtain API keys:

- **Email**: <support@nyengo-deliveries.com>
- **Developer Portal**: <https://developers.nyengo-deliveries.com>
- **Status Page**: <https://status.nyengo-deliveries.com>

---

*Document Version: 1.0.0*  
*Last Updated: December 2024*
