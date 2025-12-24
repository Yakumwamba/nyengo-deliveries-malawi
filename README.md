# Nyengo Deliveries

A complete courier/delivery management platform for African markets, featuring a Go backend API and Expo React Native mobile app.

## 🚀 Features

- **Courier Dashboard**: Track orders, earnings, and analytics
- **Real-time Tracking**: WebSocket-powered live delivery updates
- **Flexible Pricing**: Distance-based pricing with surge support
- **Multi-Currency**: Built for African currencies (ZMW, ZAR, KES, NGN, etc.)
- **Store Integration**: API for e-commerce platforms to submit orders
- **In-app Chat**: Customer-courier communication

## 📁 Project Structure

```
nyengo-deliveries/
├── backend/          # Go Fiber Backend API
├── mobile/           # Expo React Native App
└── docs/             # Documentation
```

## 🛠 Tech Stack

### Backend
- **Go** with Fiber framework
- **PostgreSQL** for primary database
- **Redis** for caching & real-time
- **JWT** authentication

### Mobile
- **Expo** / React Native
- **Expo Router** for navigation
- **Zustand** for state management
- **TypeScript**

## 🚦 Quick Start

### Backend
```bash
cd backend
cp .env.example .env
docker-compose up -d
# or
make dev
```

### Mobile
```bash
cd mobile
npm install
npx expo start
```

## 📖 Documentation

- [API Documentation](docs/API.md)
- [Setup Guide](docs/SETUP.md)
- [Architecture Overview](docs/ARCHITECTURE.md)

## 💰 Currency Configuration

Default currency is Zambian Kwacha (ZMW). Configure in `.env`:

```env
CURRENCY=ZMW
CURRENCY_SYMBOL=K
BASE_RATE_PER_KM=5.0
MINIMUM_FARE=20.0
PLATFORM_FEE_PERCENT=0.10
```

## 📱 App Screenshots

*Coming soon*

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines.

## 📄 License

MIT License - see LICENSE file for details.

---

Built with ❤️ for African delivery businesses
