# Nyengo Deliveries - Android Widget

This document explains how to set up and use the Android home screen widget for the Nyengo Deliveries courier app.

## Widget Types

The app includes three widget sizes:

### 1. Small Widget (2x1)
- **Name**: Nyengo Status
- **Features**:
  - Online/Offline status indicator
  - Quick stats (orders count, earnings)
- **Best for**: Quick status glance

### 2. Medium Widget (4x1)
- **Name**: Nyengo Delivery
- **Features**:
  - Online/Offline status with toggle badge
  - Today's orders and earnings stats
  - Active delivery info  with destination and ETA
- **Best for**: Active delivery monitoring

### 3. Large Widget (4x2)
- **Name**: Nyengo Dashboard
- **Features**:
  - Full branding with logo
  - Online/Offline toggle button
  - Stats grid (Orders, Earnings, Pending)
  - Active delivery card with full details
- **Best for**: Complete dashboard experience

## Setup Instructions

### 1. Build the App

Since widgets require native code, you need to run a development build:

```bash
# Install dependencies
npm install

# Generate native Android code
npx expo prebuild --platform android

# Run on Android device/emulator
npx expo run:android
```

### 2. Add Widget to Home Screen

1. Long-press on your Android home screen
2. Tap "Widgets"
3. Find "Nyengo" in the widget list
4. Choose your preferred widget size
5. Drag and drop to your home screen

## Widget Actions

| Action | Description |
|--------|-------------|
| Tap Widget | Opens the Nyengo app |
| Tap "Online/Offline" | Toggles your availability status |
| Tap Active Delivery | Opens order details with navigation |

## Updating Widget Data

The widget automatically updates:
- Every 30 minutes (background refresh)
- When you change online/offline status
- When you accept/complete a delivery
- When you open the app

### Programmatic Updates

Use the `useWidgetUpdates` hook in your components:

```typescript
import { useWidgetUpdates } from '../src/hooks/useWidgetUpdates';

function MyComponent() {
  const { updateOnlineStatus, updateActiveDelivery, updateStats } = useWidgetUpdates();

  // Update online status
  await updateOnlineStatus(true);

  // Update active delivery
  await updateActiveDelivery({
    orderNumber: 'NYG-20231215-ABC123',
    destination: '23 Independence Ave, Northmead',
    eta: '15 min',
  });

  // Update stats
  await updateStats({
    todaysOrders: 12,
    todaysEarnings: 'K 450',
    pending: 4,
  });
}
```

## Theme

The widget uses the Yandex/Yango-inspired theme:
- **Primary**: `#FFCC00` (Yellow)
- **Secondary**: `#1A1A1A` (Dark)
- **Accent**: `#FF3B30` (Red)
- **Success**: `#34C759` (Green - Online status)

## Troubleshooting

### Widget not appearing in list
- Make sure you've run `npx expo prebuild` and `npx expo run:android`
- Reinstall the app if widgets still don't appear

### Widget shows "Loading..." or blank
- Check that the app has been opened at least once
- Verify AsyncStorage is working properly

### Widget not updating
- Force refresh: Remove and re-add the widget
- Check battery optimization settings for the app

## File Structure

```
src/
├── widgets/
│   ├── index.ts                 # Exports
│   ├── DeliveryWidget.tsx       # Widget UI components
│   └── widgetTaskHandler.tsx    # Widget logic & data management
├── hooks/
│   └── useWidgetUpdates.ts      # Hook for app-to-widget updates
└── ...

widget-task-handler.tsx          # Entry point for widget registration
app.json                         # Widget configuration
```

## Dependencies

- `react-native-android-widget`: Widget rendering library
- `@react-native-async-storage/async-storage`: Data persistence

## Notes

- Widgets are Android-only (iOS doesn't support third-party widgets in the same way)
- Deep linking (`nyengo://`) is used for navigation from widgets
- Widget updates are batched to conserve battery
