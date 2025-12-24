# Nyengo Deliveries - Setup Guide

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)
- Expo CLI

## Backend Setup

### 1. Clone and Navigate
```bash
cd nyengo-deliveries/backend
```

### 2. Configure Environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Start with Docker (Recommended)
```bash
docker-compose up -d
```

### 4. Or Manual Setup

Install Go dependencies:
```bash
go mod tidy
```

Start PostgreSQL and Redis, then run migrations:
```bash
psql -U postgres -d nyengo -f internal/database/migrations/001_initial_schema.sql
```

Run the server:
```bash
make dev
# or
go run ./cmd/server
```

### 5. Verify
```bash
curl http://localhost:8080/health
```

## Mobile App Setup

### 1. Navigate to Mobile Directory
```bash
cd nyengo-deliveries/mobile
```

### 2. Install Dependencies
```bash
npm install
# or
yarn install
```

### 3. Configure Environment
```bash
cp .env.example .env
# Edit .env with your API URL
```

### 4. Start Development Server
```bash
npx expo start
```

### 5. Run on Device/Simulator
- Press `i` for iOS Simulator
- Press `a` for Android Emulator
- Scan QR code with Expo Go app

## Configuration

### Currency Settings
Edit `backend/.env`:
```env
CURRENCY=ZMW           # Currency code
CURRENCY_SYMBOL=K      # Currency symbol
BASE_RATE_PER_KM=5.0   # Rate per kilometer
MINIMUM_FARE=20.0      # Minimum delivery fare
PLATFORM_FEE_PERCENT=0.10  # 10% platform fee
```

### Supported Currencies
- ZMW (Zambian Kwacha) - Default
- USD (US Dollar)
- ZAR (South African Rand)
- KES (Kenyan Shilling)
- NGN (Nigerian Naira)
- GHS (Ghanaian Cedi)

## Production Deployment

### Backend (Docker)
```bash
docker build -t nyengo-api .
docker run -p 8080:8080 --env-file .env nyengo-api
```

### Mobile App
```bash
npx expo build:android
npx expo build:ios
# or use EAS Build
npx eas build --platform all
```

## Troubleshooting

### Database Connection Failed
- Check PostgreSQL is running
- Verify DATABASE_URL in .env
- Ensure database exists

### WebSocket Not Connecting
- Verify WS_URL in mobile .env
- Check CORS settings in backend
- Ensure JWT token is valid

### Price Calculation Issues
- Verify coordinates are valid
- Check BASE_RATE_PER_KM and MINIMUM_FARE settings
