# Nyengo Deliveries - Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                              │
├─────────────────────┬───────────────────────────────────────────┤
│   Mobile App        │        Store Frontend                     │
│   (Expo/React       │        (Web Integration)                  │
│    Native)          │                                           │
└──────────┬──────────┴────────────────────┬──────────────────────┘
           │                                │
           │ REST/WS                        │ REST API
           │                                │
┌──────────▼────────────────────────────────▼──────────────────────┐
│                      API Gateway (Fiber)                         │
│  ┌────────────┐  ┌─────────────┐  ┌──────────────┐              │
│  │ JWT Auth   │  │ Rate Limit  │  │ API Key Auth │              │
│  └────────────┘  └─────────────┘  └──────────────┘              │
└──────────────────────────┬───────────────────────────────────────┘
                           │
┌──────────────────────────▼───────────────────────────────────────┐
│                     Service Layer                                 │
│  ┌──────────────┐ ┌──────────────┐ ┌─────────────────┐          │
│  │ Courier      │ │ Order        │ │ Pricing         │          │
│  │ Service      │ │ Service      │ │ Service         │          │
│  ├──────────────┤ ├──────────────┤ ├─────────────────┤          │
│  │ Notification │ │ Analytics    │ │ Chat            │          │
│  │ Service      │ │ Service      │ │ Service         │          │
│  └──────────────┘ └──────────────┘ └─────────────────┘          │
└──────────────────────────┬───────────────────────────────────────┘
                           │
┌──────────────────────────▼───────────────────────────────────────┐
│                     Data Layer                                    │
│  ┌──────────────┐ ┌──────────────┐ ┌─────────────────┐          │
│  │ PostgreSQL   │ │ Redis        │ │ File Storage    │          │
│  │ (Primary DB) │ │ (Cache/PubSub│ │ (S3/Local)      │          │
│  └──────────────┘ └──────────────┘ └─────────────────┘          │
└──────────────────────────────────────────────────────────────────┘
```

## Backend Structure

```
backend/
├── cmd/server/main.go      # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── database/           # Database connections
│   ├── models/             # Data models & DTOs
│   ├── repository/         # Data access layer
│   ├── services/           # Business logic
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Request middleware
│   ├── websocket/          # Real-time communication
│   └── utils/              # Utility functions
└── pkg/validator/          # Input validation
```

## Mobile Structure

```
mobile/
├── app/                    # Expo Router pages
│   ├── (auth)/             # Auth screens
│   ├── (tabs)/             # Main tab screens
│   └── orders/             # Order detail screens
└── src/
    ├── components/         # Reusable UI components
    ├── hooks/              # Custom React hooks
    ├── services/           # API & WebSocket clients
    ├── stores/             # Zustand state stores
    ├── types/              # TypeScript types
    ├── utils/              # Utility functions
    └── constants/          # App configuration
```

## Key Design Decisions

### 1. Currency Flexibility
- Currency is configurable via environment variables
- Presets for African currencies (ZMW, ZAR, KES, NGN, GHS)
- All prices stored in smallest unit, formatted on output

### 2. Pricing System
- Distance-based pricing using Haversine formula
- Configurable base rate per km
- Support for surge pricing zones
- Express and fragile handling fees

### 3. Authentication
- JWT tokens for mobile app
- API keys for store integrations
- Role-based access control ready

### 4. Real-time Updates
- WebSocket for live order updates
- Redis Pub/Sub for notifications
- Optimistic UI updates on mobile

### 5. Offline Support (Mobile)
- Zustand for state management
- Secure local storage for tokens
- Queue failed requests for retry

## Data Flow

### Order Creation
1. Store/Courier creates order via API
2. Order Service calculates pricing
3. Order saved to PostgreSQL
4. Notification sent via WebSocket
5. Push notification to courier app

### Order Tracking
1. Driver updates location periodically
2. Location saved to tracking table
3. Redis publishes location update
4. WebSocket broadcasts to relevant clients
5. Customer sees real-time map update

## Security Considerations

- All passwords hashed with bcrypt
- JWT tokens expire after 24 hours
- API keys rotatable per store
- CORS configured for known origins
- Rate limiting on all endpoints
- Input validation on all requests
