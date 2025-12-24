import { Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SHADOWS, TYPOGRAPHY } from '../../src/constants/theme';
import { View } from 'react-native';

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: COLORS.primary,
        tabBarInactiveTintColor: COLORS.textMuted,
        tabBarStyle: {
          backgroundColor: COLORS.secondary,
          borderTopWidth: 0,
          height: 70,
          paddingBottom: 12,
          paddingTop: 8,
          ...SHADOWS.xl,
        },
        tabBarLabelStyle: {
          fontSize: TYPOGRAPHY.fontSize.xs,
          fontWeight: TYPOGRAPHY.fontWeight.semibold,
          marginTop: 4,
        },
        headerStyle: { 
          backgroundColor: COLORS.secondary,
          borderBottomWidth: 0,
          elevation: 0,
          shadowOpacity: 0,
        },
        headerTintColor: COLORS.white,
        headerTitleStyle: { 
          fontWeight: TYPOGRAPHY.fontWeight.bold,
          fontSize: TYPOGRAPHY.fontSize.lg,
        },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Dashboard',
          headerShown: false,
          tabBarIcon: ({ color, size, focused }) => (
            <View style={{
              backgroundColor: focused ? COLORS.primary + '20' : 'transparent',
              padding: 8,
              borderRadius: 12,
            }}>
              <Ionicons name={focused ? "home" : "home-outline"} size={size} color={color} />
            </View>
          ),
        }}
      />
      <Tabs.Screen
        name="orders"
        options={{
          title: 'Orders',
          tabBarIcon: ({ color, size, focused }) => (
            <View style={{
              backgroundColor: focused ? COLORS.primary + '20' : 'transparent',
              padding: 8,
              borderRadius: 12,
            }}>
              <Ionicons name={focused ? "receipt" : "receipt-outline"} size={size} color={color} />
            </View>
          ),
        }}
      />
      <Tabs.Screen
        name="earnings"
        options={{
          title: 'Earnings',
          headerShown: false,
          tabBarIcon: ({ color, size, focused }) => (
            <View style={{
              backgroundColor: focused ? COLORS.primary + '20' : 'transparent',
              padding: 8,
              borderRadius: 12,
            }}>
              <Ionicons name={focused ? "wallet" : "wallet-outline"} size={size} color={color} />
            </View>
          ),
        }}
      />
      <Tabs.Screen
        name="chat"
        options={{
          title: 'Chat',
          tabBarIcon: ({ color, size, focused }) => (
            <View style={{
              backgroundColor: focused ? COLORS.primary + '20' : 'transparent',
              padding: 8,
              borderRadius: 12,
            }}>
              <Ionicons name={focused ? "chatbubbles" : "chatbubbles-outline"} size={size} color={color} />
            </View>
          ),
          tabBarBadge: 3,
          tabBarBadgeStyle: {
            backgroundColor: COLORS.accent,
            color: COLORS.white,
            fontSize: 10,
            fontWeight: '700',
            minWidth: 18,
            height: 18,
            borderRadius: 9,
          },
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
          headerShown: false,
          tabBarIcon: ({ color, size, focused }) => (
            <View style={{
              backgroundColor: focused ? COLORS.primary + '20' : 'transparent',
              padding: 8,
              borderRadius: 12,
            }}>
              <Ionicons name={focused ? "person" : "person-outline"} size={size} color={color} />
            </View>
          ),
        }}
      />
    </Tabs>
  );
}
