import { useEffect, useRef } from 'react';
import * as Notifications from 'expo-notifications';
import { Platform } from 'react-native';

Notifications.setNotificationHandler({
    handleNotification: async () => ({
        shouldShowAlert: true,
        shouldPlaySound: true,
        shouldSetBadge: true,
    }),
});

export function useNotifications() {
    const notificationListener = useRef<any>();
    const responseListener = useRef<any>();

    useEffect(() => {
        registerForPushNotifications();

        notificationListener.current = Notifications.addNotificationReceivedListener((notification) => {
            console.log('Notification received:', notification);
        });

        responseListener.current = Notifications.addNotificationResponseReceivedListener((response) => {
            console.log('Notification response:', response);
            // Handle navigation based on notification data
        });

        return () => {
            if (notificationListener.current) {
                Notifications.removeNotificationSubscription(notificationListener.current);
            }
            if (responseListener.current) {
                Notifications.removeNotificationSubscription(responseListener.current);
            }
        };
    }, []);

    return {
        scheduleNotification: async (title: string, body: string, data?: any) => {
            await Notifications.scheduleNotificationAsync({
                content: { title, body, data },
                trigger: null,
            });
        },
    };
}

async function registerForPushNotifications() {
    if (Platform.OS === 'android') {
        await Notifications.setNotificationChannelAsync('default', {
            name: 'default',
            importance: Notifications.AndroidImportance.MAX,
            vibrationPattern: [0, 250, 250, 250],
        });
    }

    const { status: existingStatus } = await Notifications.getPermissionsAsync();
    let finalStatus = existingStatus;

    if (existingStatus !== 'granted') {
        const { status } = await Notifications.requestPermissionsAsync();
        finalStatus = status;
    }

    if (finalStatus !== 'granted') {
        console.log('Push notification permission not granted');
        return;
    }

    const token = await Notifications.getExpoPushTokenAsync();
    console.log('Push token:', token.data);
    return token.data;
}
