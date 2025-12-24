import { useEffect, useCallback } from 'react';
import { requestWidgetUpdate } from 'react-native-android-widget';
import { WidgetDataManager, WIDGET_NAMES } from '../widgets';

/**
 * Hook to manage widget updates from the app
 * Call this hook in your main app component to keep widgets in sync
 */
export function useWidgetUpdates() {
    /**
     * Update widget with new online status
     */
    const updateOnlineStatus = useCallback(async (isOnline: boolean) => {
        await WidgetDataManager.setOnlineStatus(isOnline);
        await refreshAllWidgets();
    }, []);

    /**
     * Update widget with active delivery info
     */
    const updateActiveDelivery = useCallback(async (delivery: {
        orderNumber: string;
        destination: string;
        eta: string;
    } | null) => {
        await WidgetDataManager.setActiveDelivery(delivery);
        await refreshAllWidgets();
    }, []);

    /**
     * Update widget with today's stats
     */
    const updateStats = useCallback(async (stats: {
        todaysOrders: number;
        todaysEarnings: string;
        pending: number;
    }) => {
        await WidgetDataManager.updateStats(stats);
        await refreshAllWidgets();
    }, []);

    /**
     * Clear widget data (call on logout)
     */
    const clearWidgetData = useCallback(async () => {
        await WidgetDataManager.clearAll();
        await refreshAllWidgets();
    }, []);

    /**
     * Refresh all widget instances
     */
    const refreshAllWidgets = async () => {
        try {
            await Promise.all([
                requestWidgetUpdate({
                    widgetName: WIDGET_NAMES.SMALL,
                    renderWidget: () => Promise.resolve(null as any),
                    widgetNotFound: () => { },
                }),
                requestWidgetUpdate({
                    widgetName: WIDGET_NAMES.MEDIUM,
                    renderWidget: () => Promise.resolve(null as any),
                    widgetNotFound: () => { },
                }),
                requestWidgetUpdate({
                    widgetName: WIDGET_NAMES.LARGE,
                    renderWidget: () => Promise.resolve(null as any),
                    widgetNotFound: () => { },
                }),
            ]);
        } catch (error) {
            // Widget updates may fail silently if no widgets are placed
            console.log('Widget update:', error);
        }
    };

    return {
        updateOnlineStatus,
        updateActiveDelivery,
        updateStats,
        clearWidgetData,
        refreshAllWidgets,
    };
}

export default useWidgetUpdates;
