import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { useEffect } from 'react';
import { Platform } from 'react-native';
import { useAuthStore } from '../src/stores/authStore';
import { COLORS } from '../src/constants/theme';

// Import widget registration (Android only)
if (Platform.OS === 'android') {
  require('../widget-task-handler');
}

export default function RootLayout() {
  const { checkAuth, courier } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, []);

  // Update widgets when auth state changes (Android only)
  useEffect(() => {
    if (Platform.OS === 'android' && courier) {
      // Widget updates will be handled by individual screens
      // This is just a placeholder for future implementation
    }
  }, [courier]);

  return (
    <>
      <StatusBar style="light" />
      <Stack
        screenOptions={{
          headerStyle: { backgroundColor: COLORS.secondary },
          headerTintColor: COLORS.white,
          headerTitleStyle: { fontWeight: 'bold' },
          contentStyle: { backgroundColor: COLORS.background },
        }}
      >
        <Stack.Screen name="(auth)" options={{ headerShown: false }} />
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen 
          name="orders/[id]" 
          options={{ 
            title: 'Order Details',
            presentation: 'card' 
          }} 
        />
      </Stack>
    </>
  );
}
