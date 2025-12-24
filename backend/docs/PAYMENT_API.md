# Nyengo Deliveries - Payment & Payout API Documentation

## Overview

This document describes the Payment Verification and Courier Payout APIs for the Nyengo Deliveries system. These APIs allow:
- **Couriers** to verify payment for delivered orders and request payouts
- **Stores** to query payment status and confirm payments
- **Admins** to approve or reject payout requests

## Base URL
```
https://your-domain.com/api/v1
```

## Authentication

### Courier Authentication (JWT)
Use the Bearer token received from `/couriers/login` in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Store Authentication (API Key)
Use your store API key in the X-API-Key header:
```
X-API-Key: <your_api_key>
```

---

## Courier Payment Endpoints

### 1. Verify Order Payment
Check if payment has been received for a delivered order.

**Endpoint:** `GET /payments/verify/{orderId}`

**Authentication:** JWT (Courier)

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| orderId | UUID | The order ID to verify |

**Response:**
```json
{
  "success": true,
  "data": {
    "orderId": "123e4567-e89b-12d3-a456-426614174000",
    "isPaid": true,
    "amountPaid": 150.00,
    "paymentMethod": "mobile_money",
    "paymentRef": "MOMO-TXN-12345678",
    "paidAt": "2025-12-16T08:30:00Z"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "order not yet delivered"
}
```

---

### 2. Get Payable Orders
Get list of delivered orders eligible for payout.

**Endpoint:** `GET /payments/payable-orders`

**Authentication:** JWT (Courier)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "orderId": "123e4567-e89b-12d3-a456-426614174000",
      "orderNumber": "NYG-20251216-A1B2C3D4",
      "deliveredAt": "2025-12-16T08:30:00Z",
      "totalFare": 150.00,
      "courierEarnings": 135.00,
      "paymentStatus": "paid",
      "payoutStatus": "unpaid",
      "customerName": "John Doe"
    },
    {
      "orderId": "223e4567-e89b-12d3-a456-426614174001",
      "orderNumber": "NYG-20251216-E5F6G7H8",
      "deliveredAt": "2025-12-16T10:15:00Z",
      "totalFare": 200.00,
      "courierEarnings": 180.00,
      "paymentStatus": "paid",
      "payoutStatus": "unpaid",
      "customerName": "Jane Smith"
    }
  ]
}
```

---

### 3. Request Payout
Request payment for completed deliveries.

**Endpoint:** `POST /payments/payouts`

**Authentication:** JWT (Courier)

**Request Body:**
```json
{
  "orderIds": [
    "123e4567-e89b-12d3-a456-426614174000",
    "223e4567-e89b-12d3-a456-426614174001"
  ],
  "payoutMethod": "mobile_money",
  "payoutDetails": "+260970123456"
}
```

**Payout Methods:**
| Method | Description | Details Format |
|--------|-------------|----------------|
| `bank_transfer` | Bank account transfer | JSON with bankName, accountNumber, accountName |
| `mobile_money` | Mobile money transfer | Phone number (e.g., +260970123456) |
| `wallet` | Keep in app wallet | Leave empty |

**Response:**
```json
{
  "success": true,
  "message": "Payout request submitted successfully. Processing will begin shortly.",
  "data": {
    "payoutId": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "totalAmount": "K 315.00",
    "netAmount": "K 283.50",
    "message": "Payout request submitted successfully. Processing will begin shortly."
  }
}
```

**Error Responses:**
```json
{
  "success": false,
  "error": "order 123e4567-e89b-12d3-a456-426614174000 payment not confirmed: payment not yet confirmed"
}
```

```json
{
  "success": false,
  "error": "courier account must be verified to request payouts"
}
```

---

### 4. Get Payout History
Retrieve paginated list of payout requests.

**Endpoint:** `GET /payments/payouts`

**Authentication:** JWT (Courier)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| status | string | - | Filter by status: pending, processing, completed, failed, cancelled |
| page | int | 1 | Page number |
| limit | int | 20 | Items per page (max 100) |

**Example:** `GET /payments/payouts?status=completed&page=1&limit=10`

**Response:**
```json
{
  "success": true,
  "data": {
    "payouts": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "courierId": "660e8400-e29b-41d4-a716-446655440000",
        "orderIds": ["123e4567-e89b-12d3-a456-426614174000"],
        "totalAmount": 315.00,
        "platformFee": 31.50,
        "netAmount": 283.50,
        "currency": "ZMW",
        "status": "completed",
        "payoutMethod": "mobile_money",
        "payoutDetails": "+260970123456",
        "transactionRef": "MOMO-TXN-550e8400",
        "processedAt": "2025-12-16T09:00:00Z",
        "completedAt": "2025-12-16T09:05:00Z",
        "createdAt": "2025-12-16T08:45:00Z",
        "updatedAt": "2025-12-16T09:05:00Z"
      }
    ],
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 5. Get Payout Details
Retrieve details of a specific payout.

**Endpoint:** `GET /payments/payouts/{payoutId}`

**Authentication:** JWT (Courier)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "courierId": "660e8400-e29b-41d4-a716-446655440000",
    "orderIds": ["123e4567-e89b-12d3-a456-426614174000"],
    "totalAmount": 315.00,
    "platformFee": 31.50,
    "netAmount": 283.50,
    "currency": "ZMW",
    "status": "pending",
    "payoutMethod": "mobile_money",
    "payoutDetails": "+260970123456",
    "createdAt": "2025-12-16T08:45:00Z",
    "updatedAt": "2025-12-16T08:45:00Z"
  }
}
```

---

### 6. Get Earnings Summary
Get summary of courier's earnings.

**Endpoint:** `GET /payments/earnings`

**Authentication:** JWT (Courier)

**Response:**
```json
{
  "success": true,
  "data": {
    "totalEarnings": 5000.00,
    "availableBalance": 1500.00,
    "pendingBalance": 500.00,
    "totalPaidOut": 3000.00,
    "unpaidOrders": 5,
    "pendingPayouts": 1,
    "currency": "ZMW",
    "formattedTotal": "K 5,000.00",
    "formattedAvailable": "K 1,500.00"
  }
}
```

---

### 7. Get Wallet Transactions
Retrieve wallet transaction history.

**Endpoint:** `GET /payments/wallet/transactions`

**Authentication:** JWT (Courier)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 20 | Items per page (max 100) |

**Response:**
```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "courierId": "660e8400-e29b-41d4-a716-446655440000",
        "orderId": "123e4567-e89b-12d3-a456-426614174000",
        "type": "earning",
        "amount": 135.00,
        "balanceBefore": 1365.00,
        "balanceAfter": 1500.00,
        "description": "Earnings for order NYG-20251216-A1B2C3D4",
        "reference": "NYG-20251216-A1B2C3D4",
        "createdAt": "2025-12-16T08:35:00Z"
      },
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "courierId": "660e8400-e29b-41d4-a716-446655440000",
        "payoutId": "550e8400-e29b-41d4-a716-446655440000",
        "type": "payout",
        "amount": -283.50,
        "balanceBefore": 1500.00,
        "balanceAfter": 1216.50,
        "description": "Payout for 2 orders",
        "reference": "550e8400-e29b-41d4-a716-446655440000",
        "createdAt": "2025-12-16T09:05:00Z"
      }
    ],
    "page": 1,
    "pageSize": 20
  }
}
```

---

## Store Integration Endpoints

### 1. Verify Payment (Store)
Allows stores to check payment status for their orders.

**Endpoint:** `GET /stores/payments/verify/{orderId}`

**Authentication:** API Key

**Response:** Same as Courier Payment Verification

---

### 2. Payment Confirmation Webhook
Stores can send payment confirmations via webhook.

**Endpoint:** `POST /webhooks/payment-confirm`

**Authentication:** API Key

**Request Body:**
```json
{
  "orderId": "123e4567-e89b-12d3-a456-426614174000",
  "externalOrderId": "STORE-ORDER-12345",
  "status": "paid",
  "amount": 150.00,
  "currency": "ZMW",
  "transactionRef": "PAY-REF-98765",
  "paymentMethod": "mobile_money",
  "paidAt": "2025-12-16T08:30:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payment status updated",
  "data": {
    "orderId": "123e4567-e89b-12d3-a456-426614174000",
    "isPaid": true,
    "amountPaid": 150.00,
    "paymentMethod": "mobile_money",
    "paymentRef": "PAY-REF-98765",
    "paidAt": "2025-12-16T08:30:00Z"
  }
}
```

---

## Admin Endpoints

### 1. Process Payout Request
Approve or reject a pending payout request.

**Endpoint:** `POST /admin/payments/payouts/{payoutId}/process`

**Authentication:** JWT (Admin)

**Request Body:**
```json
{
  "approve": true,
  "reason": ""
}
```

Or to reject:
```json
{
  "approve": false,
  "reason": "Insufficient verification documents"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payout approved successfully"
}
```

---

## Status Codes

| Status | Description |
|--------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid or missing authentication |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## Payout Status Flow

```
pending → processing → completed
     ↘            ↘
      cancelled    failed
```

| Status | Description |
|--------|-------------|
| `pending` | Payout request submitted, awaiting processing |
| `processing` | Payout is being processed |
| `completed` | Payout successfully sent to courier |
| `failed` | Payout failed (see failureReason) |
| `cancelled` | Payout was cancelled/rejected |

---

## Payment Verification Logic

1. **Cash Payments**: Verified automatically when delivery proof (photo/signature) is uploaded
2. **Store Orders**: Verified by querying store's payment API or via webhook
3. **Mobile Money/Card**: Verified via payment reference

---

## Store Payment API Integration

To integrate your store's payment system:

### 1. Configure Store Payment API

Contact Nyengo support to configure your payment API endpoint:
- **API URL**: Your payment verification endpoint
- **API Key**: Your secret key for authentication
- **Webhook Secret**: For securing webhook callbacks

### 2. Implement Payment Verification Endpoint

Your API should implement:

**Endpoint:** `POST /verify-payment`

**Request Headers:**
```
Content-Type: application/json
X-API-Key: <your_api_key>
X-Timestamp: 2025-12-16T08:30:00Z
X-Signature: <hmac_sha256_signature>
```

**Request Body:**
```json
{
  "orderId": "123e4567-e89b-12d3-a456-426614174000",
  "externalOrderId": "STORE-ORDER-12345",
  "amount": 150.00,
  "currency": "ZMW"
}
```

**Expected Response:**
```json
{
  "isPaid": true,
  "amountPaid": 150.00,
  "paymentMethod": "card",
  "transactionRef": "STORE-TXN-12345",
  "paidAt": "2025-12-16T08:25:00Z"
}
```

### 3. Signature Verification

The `X-Signature` header contains an HMAC-SHA256 signature of the request body using your webhook secret:

```python
import hmac
import hashlib

signature = hmac.new(
    webhook_secret.encode(),
    request_body,
    hashlib.sha256
).hexdigest()
```

---

## Testing Examples

### cURL Examples

**1. Verify Payment:**
```bash
curl -X GET "https://api.nyengo.com/api/v1/payments/verify/123e4567-e89b-12d3-a456-426614174000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**2. Get Payable Orders:**
```bash
curl -X GET "https://api.nyengo.com/api/v1/payments/payable-orders" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**3. Request Payout:**
```bash
curl -X POST "https://api.nyengo.com/api/v1/payments/payouts" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "orderIds": ["123e4567-e89b-12d3-a456-426614174000"],
    "payoutMethod": "mobile_money",
    "payoutDetails": "+260970123456"
  }'
```

**4. Store Verify Payment:**
```bash
curl -X GET "https://api.nyengo.com/api/v1/stores/payments/verify/123e4567-e89b-12d3-a456-426614174000" \
  -H "X-API-Key: YOUR_STORE_API_KEY"
```

**5. Admin Process Payout:**
```bash
curl -X POST "https://api.nyengo.com/api/v1/admin/payments/payouts/550e8400-e29b-41d4-a716-446655440000/process" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "approve": true,
    "reason": ""
  }'
```

---

## JavaScript/TypeScript Examples

```typescript
// Using fetch API

const API_BASE = 'https://api.nyengo.com/api/v1';

// 1. Verify Payment
async function verifyPayment(orderId: string, token: string) {
  const response = await fetch(`${API_BASE}/payments/verify/${orderId}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}

// 2. Request Payout
async function requestPayout(
  orderIds: string[],
  payoutMethod: string,
  payoutDetails: string,
  token: string
) {
  const response = await fetch(`${API_BASE}/payments/payouts`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      orderIds,
      payoutMethod,
      payoutDetails
    })
  });
  return response.json();
}

// 3. Get Earnings Summary
async function getEarnings(token: string) {
  const response = await fetch(`${API_BASE}/payments/earnings`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}

// Store Integration - Payment Webhook
async function sendPaymentConfirmation(orderData: any, apiKey: string) {
  const response = await fetch(`${API_BASE}/webhooks/payment-confirm`, {
    method: 'POST',
    headers: {
      'X-API-Key': apiKey,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      orderId: orderData.nyengoOrderId,
      externalOrderId: orderData.storeOrderId,
      status: 'paid',
      amount: orderData.amount,
      currency: 'ZMW',
      transactionRef: orderData.transactionRef,
      paymentMethod: orderData.paymentMethod,
      paidAt: new Date().toISOString()
    })
  });
  return response.json();
}
```

---

## Support

For integration support, contact:
- **Email**: developers@nyengo.com
- **Documentation**: https://docs.nyengo.com/api
- **API Status**: https://status.nyengo.com
