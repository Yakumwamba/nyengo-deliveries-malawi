import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Alert, Switch } from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useDriverTracking } from '../../src/hooks/useTracking';
import { useAuthStore } from '../../src/stores/authStore';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

export default function DriverTrackingScreen() {
  const { orderId } = useLocalSearchParams<{ orderId: string }>();
  const { courier } = useAuthStore();
  const { isTracking, lastUpdate, error, startTracking, stopTracking } = useDriverTracking(orderId);
  const [isStarting, setIsStarting] = useState(false);

  const handleStartTracking = async () => {
    setIsStarting(true);
    try {
      await startTracking({
        name: courier?.ownerName || 'Driver',
        phone: courier?.phone || '',
        vehicleType: courier?.vehicleTypes?.[0] || 'motorcycle',
      });
      Alert.alert('Tracking Started', 'Your location is now being shared with the customer.');
    } catch (e: any) {
      Alert.alert('Error', e.message);
    } finally {
      setIsStarting(false);
    }
  };

  const handleStopTracking = async () => {
    Alert.alert(
      'Stop Tracking',
      'Are you sure you want to stop sharing your location?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Stop',
          style: 'destructive',
          onPress: async () => {
            await stopTracking('delivery_completed');
            router.back();
          },
        },
      ]
    );
  };

  return (
    <View style={styles.container}>
      {/* Status Header */}
      <View style={[styles.statusCard, isTracking ? styles.trackingActive : styles.trackingInactive]}>
        <View style={[styles.statusIconContainer, isTracking && styles.statusIconActive]}>
          <Ionicons
            name={isTracking ? 'location' : 'location-outline'}
            size={40}
            color={isTracking ? COLORS.secondary : COLORS.textMuted}
          />
        </View>
        <Text style={[styles.statusTitle, isTracking && styles.statusTitleActive]}>
          {isTracking ? 'Tracking Active' : 'Tracking Off'}
        </Text>
        <Text style={[styles.statusSubtitle, isTracking && styles.statusSubtitleActive]}>
          {isTracking
            ? 'Your location is being shared with the customer'
            : 'Start tracking to share your live location'}
        </Text>

        {lastUpdate && isTracking && (
          <View style={styles.lastUpdateBadge}>
            <Ionicons name="time" size={12} color={COLORS.secondary} />
            <Text style={styles.lastUpdateText}>
              Last update: {lastUpdate.toLocaleTimeString()}
            </Text>
          </View>
        )}
      </View>

      {/* Order Info */}
      <View style={styles.orderCard}>
        <View style={styles.orderIcon}>
          <Ionicons name="cube" size={20} color={COLORS.primary} />
        </View>
        <View style={styles.orderInfo}>
          <Text style={styles.orderLabel}>Order ID</Text>
          <Text style={styles.orderNumber}>{orderId}</Text>
        </View>
        <Ionicons name="chevron-forward" size={20} color={COLORS.textMuted} />
      </View>

      {/* Tracking Controls */}
      <View style={styles.controlsCard}>
        <View style={styles.controlRow}>
          <View style={styles.controlLeft}>
            <View style={[styles.controlIcon, { backgroundColor: COLORS.primary + '15' }]}>
              <Ionicons name="navigate" size={20} color={COLORS.primary} />
            </View>
            <View>
              <Text style={styles.controlLabel}>Location Sharing</Text>
              <Text style={styles.controlSubLabel}>
                {isTracking ? 'Customer can see your location' : 'Turn on to share location'}
              </Text>
            </View>
          </View>
          <Switch
            value={isTracking}
            onValueChange={(value) => {
              if (value) {
                handleStartTracking();
              } else {
                handleStopTracking();
              }
            }}
            trackColor={{ false: COLORS.border, true: COLORS.primary }}
            thumbColor={COLORS.white}
            disabled={isStarting}
          />
        </View>
      </View>

      {/* Action Button */}
      {!isTracking ? (
        <TouchableOpacity
          style={[styles.startButton, isStarting && styles.buttonDisabled]}
          onPress={handleStartTracking}
          disabled={isStarting}
        >
          <Ionicons name="navigate" size={24} color={COLORS.secondary} />
          <Text style={styles.startButtonText}>
            {isStarting ? 'Starting...' : 'Start Tracking'}
          </Text>
        </TouchableOpacity>
      ) : (
        <TouchableOpacity style={styles.stopButton} onPress={handleStopTracking}>
          <Ionicons name="checkmark-circle" size={24} color={COLORS.white} />
          <Text style={styles.stopButtonText}>Complete Delivery</Text>
        </TouchableOpacity>
      )}

      {/* Tips */}
      <View style={styles.tipsCard}>
        <View style={styles.tipsHeader}>
          <Ionicons name="bulb" size={18} color={COLORS.primary} />
          <Text style={styles.tipsTitle}>Tips</Text>
        </View>
        <View style={styles.tipsList}>
          <View style={styles.tipItem}>
            <Ionicons name="checkmark-circle" size={14} color={COLORS.success} />
            <Text style={styles.tipText}>Keep GPS enabled for accurate tracking</Text>
          </View>
          <View style={styles.tipItem}>
            <Ionicons name="checkmark-circle" size={14} color={COLORS.success} />
            <Text style={styles.tipText}>Battery saver mode may affect accuracy</Text>
          </View>
          <View style={styles.tipItem}>
            <Ionicons name="checkmark-circle" size={14} color={COLORS.success} />
            <Text style={styles.tipText}>Tracking stops when you complete delivery</Text>
          </View>
        </View>
      </View>

      {error && (
        <View style={styles.errorCard}>
          <Ionicons name="warning" size={20} color={COLORS.accent} />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
    padding: SPACING.base,
  },
  
  // Status Card
  statusCard: {
    alignItems: 'center',
    padding: SPACING['2xl'],
    borderRadius: RADIUS.xl,
    marginBottom: SPACING.base,
    ...SHADOWS.md,
  },
  trackingActive: {
    backgroundColor: COLORS.primary,
  },
  trackingInactive: {
    backgroundColor: COLORS.surface,
  },
  statusIconContainer: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: COLORS.background,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  statusIconActive: {
    backgroundColor: COLORS.secondary,
  },
  statusTitle: {
    fontSize: TYPOGRAPHY.fontSize['2xl'],
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.text,
  },
  statusTitleActive: {
    color: COLORS.secondary,
  },
  statusSubtitle: {
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.textSecondary,
    marginTop: SPACING.sm,
    textAlign: 'center',
  },
  statusSubtitleActive: {
    color: COLORS.secondary + 'CC',
  },
  lastUpdateBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.secondary + '20',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    marginTop: SPACING.base,
    gap: 6,
  },
  lastUpdateText: {
    color: COLORS.secondary,
    fontSize: TYPOGRAPHY.fontSize.sm,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  
  // Order Card
  orderCard: {
    backgroundColor: COLORS.surface,
    padding: SPACING.base,
    borderRadius: RADIUS.lg,
    marginBottom: SPACING.base,
    flexDirection: 'row',
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  orderIcon: {
    width: 44,
    height: 44,
    borderRadius: RADIUS.md,
    backgroundColor: COLORS.primary + '15',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  orderInfo: {
    flex: 1,
  },
  orderLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  orderNumber: {
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
    marginTop: 2,
  },
  
  // Controls Card
  controlsCard: {
    backgroundColor: COLORS.surface,
    padding: SPACING.base,
    borderRadius: RADIUS.lg,
    marginBottom: SPACING.base,
    ...SHADOWS.sm,
  },
  controlRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  controlLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
    gap: SPACING.md,
  },
  controlIcon: {
    width: 40,
    height: 40,
    borderRadius: RADIUS.md,
    justifyContent: 'center',
    alignItems: 'center',
  },
  controlLabel: {
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  controlSubLabel: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textMuted,
    marginTop: 2,
  },
  
  // Buttons
  startButton: {
    backgroundColor: COLORS.primary,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING.lg,
    borderRadius: RADIUS.lg,
    gap: SPACING.sm,
    marginBottom: SPACING.base,
    ...SHADOWS.primary,
  },
  buttonDisabled: {
    opacity: 0.7,
  },
  startButtonText: {
    color: COLORS.secondary,
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  stopButton: {
    backgroundColor: COLORS.success,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: SPACING.lg,
    borderRadius: RADIUS.lg,
    gap: SPACING.sm,
    marginBottom: SPACING.base,
    ...SHADOWS.md,
  },
  stopButtonText: {
    color: COLORS.white,
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Tips Card
  tipsCard: {
    backgroundColor: COLORS.primaryMuted,
    padding: SPACING.base,
    borderRadius: RADIUS.lg,
    marginBottom: SPACING.base,
  },
  tipsHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    marginBottom: SPACING.md,
  },
  tipsTitle: {
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  tipsList: {
    gap: SPACING.sm,
  },
  tipItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  tipText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textSecondary,
    flex: 1,
  },
  
  // Error
  errorCard: {
    backgroundColor: COLORS.accentMuted,
    padding: SPACING.md,
    borderRadius: RADIUS.md,
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  errorText: {
    color: COLORS.accent,
    fontSize: TYPOGRAPHY.fontSize.base,
    flex: 1,
  },
});
