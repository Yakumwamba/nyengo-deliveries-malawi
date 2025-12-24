import React from 'react';
import type { WidgetTaskHandlerProps } from 'react-native-android-widget';
import { SmallDeliveryWidget, MediumDeliveryWidget, LargeDeliveryWidget } from './DeliveryWidget';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Storage keys
const STORAGE_KEYS = {
  IS_ONLINE: '@nyengo_is_online',
  ACTIVE_DELIVERY: '@nyengo_active_delivery',
  STATS: '@nyengo_widget_stats',
};

// Widget names matching app.json config
export const WIDGET_NAMES = {
  SMALL: 'NyengoSmallWidget',
  MEDIUM: 'NyengoMediumWidget',
  LARGE: 'NyengoLargeWidget',
};

/**
 * Get widget data from storage
 */
async function getWidgetData() {
  try {
    const [isOnlineStr, activeDeliveryStr, statsStr] = await Promise.all([
      AsyncStorage.getItem(STORAGE_KEYS.IS_ONLINE),
      AsyncStorage.getItem(STORAGE_KEYS.ACTIVE_DELIVERY),
      AsyncStorage.getItem(STORAGE_KEYS.STATS),
    ]);

    const isOnline = isOnlineStr === 'true';
    const activeDelivery = activeDeliveryStr ? JSON.parse(activeDeliveryStr) : null;
    const stats = statsStr
      ? JSON.parse(statsStr)
      : { todaysOrders: 0, todaysEarnings: 'K 0', pending: 0 };

    return { isOnline, activeDelivery, stats };
  } catch (error) {
    console.error('Error getting widget data:', error);
    return {
      isOnline: false,
      activeDelivery: null,
      stats: { todaysOrders: 0, todaysEarnings: 'K 0', pending: 0 },
    };
  }
}

/**
 * Toggle courier online status
 */
async function toggleOnlineStatus(): Promise<boolean> {
  try {
    const currentStatus = await AsyncStorage.getItem(STORAGE_KEYS.IS_ONLINE);
    const newStatus = currentStatus !== 'true';
    await AsyncStorage.setItem(STORAGE_KEYS.IS_ONLINE, String(newStatus));
    return newStatus;
  } catch (error) {
    console.error('Error toggling status:', error);
    return false;
  }
}

/**
 * Widget Task Handler - Called by Android when widget actions occur
 */
export async function widgetTaskHandler(props: WidgetTaskHandlerProps): Promise<React.ReactElement> {
  const { widgetName, widgetAction, clickAction, clickActionData, widgetInfo } = props;

  // Handle click actions
  if (widgetAction === 'WIDGET_CLICK') {
    switch (clickAction) {
      case 'TOGGLE_STATUS':
        await toggleOnlineStatus();
        break;
      case 'NAVIGATE':
        // The app will handle navigation when opened with this deep link
        // Intent: nyengo://orders/{orderId}
        break;
      case 'OPEN_APP':
      default:
        // Just opens the app
        break;
    }
  }

  // Get latest data
  const widgetData = await getWidgetData();

  // Render appropriate widget based on name
  switch (widgetName) {
    case WIDGET_NAMES.SMALL:
      return <SmallDeliveryWidget {...widgetData} />;
    case WIDGET_NAMES.MEDIUM:
      return <MediumDeliveryWidget {...widgetData} />;
    case WIDGET_NAMES.LARGE:
    default:
      return <LargeDeliveryWidget {...widgetData} />;
  }
}

/**
 * Helper functions to update widget data from the app
 */
export const WidgetDataManager = {
  /**
   * Update online status and refresh widgets
   */
  async setOnlineStatus(isOnline: boolean) {
    await AsyncStorage.setItem(STORAGE_KEYS.IS_ONLINE, String(isOnline));
  },

  /**
   * Update active delivery info
   */
  async setActiveDelivery(delivery: {
    orderNumber: string;
    destination: string;
    eta: string;
  } | null) {
    if (delivery) {
      await AsyncStorage.setItem(STORAGE_KEYS.ACTIVE_DELIVERY, JSON.stringify(delivery));
    } else {
      await AsyncStorage.removeItem(STORAGE_KEYS.ACTIVE_DELIVERY);
    }
  },

  /**
   * Update today's stats
   */
  async updateStats(stats: {
    todaysOrders: number;
    todaysEarnings: string;
    pending: number;
  }) {
    await AsyncStorage.setItem(STORAGE_KEYS.STATS, JSON.stringify(stats));
  },

  /**
   * Clear all widget data (on logout)
   */
  async clearAll() {
    await Promise.all([
      AsyncStorage.removeItem(STORAGE_KEYS.IS_ONLINE),
      AsyncStorage.removeItem(STORAGE_KEYS.ACTIVE_DELIVERY),
      AsyncStorage.removeItem(STORAGE_KEYS.STATS),
    ]);
  },
};

export default widgetTaskHandler;
