export const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://192.168.80.104:8082/api/v1';
export const WS_URL = process.env.EXPO_PUBLIC_WS_URL || 'ws://192.168.80.104:8082/ws';

export const APP_CONFIG = {
    name: process.env.EXPO_PUBLIC_APP_NAME || 'Nyengo Deliveries',
    currency: process.env.EXPO_PUBLIC_CURRENCY || 'ZMW',
    currencySymbol: process.env.EXPO_PUBLIC_CURRENCY_SYMBOL || 'K',
};

// Re-export colors from theme for backward compatibility
export { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS, STATUS_COLORS, COMPONENT_STYLES } from './theme';

// Legacy COLORS export (deprecated - use theme.ts instead)
export const LEGACY_COLORS = {
    primary: '#FFCC00',
    primaryDark: '#E6B800',
    secondary: '#1A1A1A',
    success: '#34C759',
    warning: '#FF9500',
    danger: '#FF3B30',
    background: '#F8F8F8',
    surface: '#FFFFFF',
    text: '#1A1A1A',
    textSecondary: '#666666',
    textMuted: '#999999',
    border: '#E5E5E5',
};

export const ORDER_STATUS = {
    PENDING: 'pending',
    ACCEPTED: 'accepted',
    PICKED_UP: 'picked_up',
    IN_TRANSIT: 'in_transit',
    DELIVERED: 'delivered',
    CANCELLED: 'cancelled',
    FAILED: 'failed',
};

export const PAYMENT_STATUS = {
    PENDING: 'pending',
    PAID: 'paid',
    FAILED: 'failed',
    REFUNDED: 'refunded',
};

export const PACKAGE_SIZES = [
    { value: 'small', label: 'Small', description: 'Documents, small boxes' },
    { value: 'medium', label: 'Medium', description: 'Standard packages' },
    { value: 'large', label: 'Large', description: 'Large items, furniture' },
];

export const VEHICLE_TYPES = [
    { value: 'bicycle', label: 'Bicycle', icon: 'bicycle' },
    { value: 'motorcycle', label: 'Motorcycle', icon: 'speedometer' },
    { value: 'car', label: 'Car', icon: 'car' },
    { value: 'van', label: 'Van', icon: 'bus' },
    { value: 'truck', label: 'Truck', icon: 'cube' },
];
