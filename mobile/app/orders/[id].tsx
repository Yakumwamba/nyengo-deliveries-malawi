import React, { useState } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, Linking, Alert } from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS, STATUS_COLORS } from '../../src/constants/theme';

export default function OrderDetailsScreen() {
  const { id } = useLocalSearchParams();
  const [isAccepting, setIsAccepting] = useState(false);

  const order = {
    id,
    orderNumber: 'NYG-20231215-ABC123',
    status: 'in_transit',
    customer: { name: 'John Doe', phone: '+260 97 123 4567' },
    pickup: { address: 'Cairo Road, Shop 45', lat: -15.4167, lng: 28.2833 },
    delivery: { address: '23 Independence Ave, Northmead', lat: -15.4101, lng: 28.3122 },
    package: { description: 'Electronics - Laptop', size: 'medium', weight: 3.5 },
    pricing: { base: 15, distance: 25, total: 45, earnings: 40.5 },
    distance: 5.2,
    estimatedTime: '15 min',
    createdAt: '2023-12-15T10:30:00Z',
  };

  const statusConfig = STATUS_COLORS[order.status as keyof typeof STATUS_COLORS] || STATUS_COLORS.pending;

  const handleCall = () => {
    Linking.openURL(`tel:${order.customer.phone}`);
  };

  const handleOpenMaps = () => {
    const url = `https://www.google.com/maps/dir/?api=1&destination=${order.delivery.lat},${order.delivery.lng}`;
    Linking.openURL(url);
  };

  const handleAccept = () => {
    setIsAccepting(true);
    setTimeout(() => {
      setIsAccepting(false);
      Alert.alert('Success', 'Order accepted successfully!');
    }, 1500);
  };

  return (
    <View style={styles.container}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Status Header */}
        <View style={styles.statusHeader}>
          <View style={[styles.statusBadge, { backgroundColor: statusConfig.bg }]}>
            <View style={[styles.statusDot, { backgroundColor: statusConfig.icon }]} />
            <Text style={[styles.statusText, { color: statusConfig.text }]}>
              {order.status.replace('_', ' ')}
            </Text>
          </View>
          <Text style={styles.orderNumber}>{order.orderNumber}</Text>
          <View style={styles.headerStats}>
            <View style={styles.headerStat}>
              <Ionicons name="map" size={18} color={COLORS.primary} />
              <Text style={styles.headerStatValue}>{order.distance} km</Text>
            </View>
            <View style={styles.headerStatDivider} />
            <View style={styles.headerStat}>
              <Ionicons name="time" size={18} color={COLORS.primary} />
              <Text style={styles.headerStatValue}>{order.estimatedTime}</Text>
            </View>
          </View>
        </View>

        {/* Customer Info */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Customer</Text>
          <View style={styles.card}>
            <View style={styles.customerRow}>
              <View style={styles.customerAvatar}>
                <Text style={styles.customerAvatarText}>
                  {order.customer.name.charAt(0)}
                </Text>
              </View>
              <View style={styles.customerInfo}>
                <Text style={styles.customerName}>{order.customer.name}</Text>
                <Text style={styles.customerPhone}>{order.customer.phone}</Text>
              </View>
              <View style={styles.customerActions}>
                <TouchableOpacity 
                  style={[styles.customerActionBtn, { backgroundColor: COLORS.success + '15' }]}
                  onPress={handleCall}
                >
                  <Ionicons name="call" size={18} color={COLORS.success} />
                </TouchableOpacity>
                <TouchableOpacity 
                  style={[styles.customerActionBtn, { backgroundColor: COLORS.info + '15' }]}
                >
                  <Ionicons name="chatbubble" size={18} color={COLORS.info} />
                </TouchableOpacity>
              </View>
            </View>
          </View>
        </View>

        {/* Route */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Route</Text>
          <View style={styles.card}>
            <View style={styles.routeContainer}>
              <View style={styles.routeTimeline}>
                <View style={[styles.routeDot, { backgroundColor: COLORS.success }]} />
                <View style={styles.routeLine} />
                <View style={[styles.routeDot, { backgroundColor: COLORS.accent }]} />
              </View>
              <View style={styles.routeDetails}>
                <View style={styles.routePoint}>
                  <Text style={styles.routeLabel}>PICKUP</Text>
                  <Text style={styles.routeAddress}>{order.pickup.address}</Text>
                </View>
                <View style={styles.routePointSpacing} />
                <View style={styles.routePoint}>
                  <Text style={styles.routeLabel}>DELIVERY</Text>
                  <Text style={styles.routeAddress}>{order.delivery.address}</Text>
                </View>
              </View>
            </View>
            <TouchableOpacity style={styles.navigateButton} onPress={handleOpenMaps}>
              <Ionicons name="navigate" size={20} color={COLORS.secondary} />
              <Text style={styles.navigateText}>Open in Maps</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Package Info */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Package Details</Text>
          <View style={styles.card}>
            <View style={styles.packageHeader}>
              <View style={styles.packageIcon}>
                <Ionicons name="cube" size={24} color={COLORS.primary} />
              </View>
              <Text style={styles.packageDesc}>{order.package.description}</Text>
            </View>
            <View style={styles.packageDetails}>
              <View style={styles.packageItem}>
                <View style={[styles.packageItemIcon, { backgroundColor: COLORS.primary + '15' }]}>
                  <Ionicons name="expand-outline" size={16} color={COLORS.primary} />
                </View>
                <View>
                  <Text style={styles.packageItemLabel}>Size</Text>
                  <Text style={styles.packageItemValue}>{order.package.size}</Text>
                </View>
              </View>
              <View style={styles.packageItem}>
                <View style={[styles.packageItemIcon, { backgroundColor: COLORS.accent + '15' }]}>
                  <Ionicons name="scale-outline" size={16} color={COLORS.accent} />
                </View>
                <View>
                  <Text style={styles.packageItemLabel}>Weight</Text>
                  <Text style={styles.packageItemValue}>{order.package.weight} kg</Text>
                </View>
              </View>
              <View style={styles.packageItem}>
                <View style={[styles.packageItemIcon, { backgroundColor: COLORS.success + '15' }]}>
                  <Ionicons name="speedometer-outline" size={16} color={COLORS.success} />
                </View>
                <View>
                  <Text style={styles.packageItemLabel}>Distance</Text>
                  <Text style={styles.packageItemValue}>{order.distance} km</Text>
                </View>
              </View>
            </View>
          </View>
        </View>

        {/* Pricing */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Pricing</Text>
          <View style={styles.card}>
            <View style={styles.pricingRow}>
              <Text style={styles.pricingLabel}>Base fare</Text>
              <Text style={styles.pricingValue}>K {order.pricing.base}</Text>
            </View>
            <View style={styles.pricingRow}>
              <Text style={styles.pricingLabel}>Distance fare</Text>
              <Text style={styles.pricingValue}>K {order.pricing.distance}</Text>
            </View>
            <View style={styles.divider} />
            <View style={styles.pricingRow}>
              <Text style={styles.pricingLabelBold}>Total</Text>
              <Text style={styles.pricingValueBold}>K {order.pricing.total}</Text>
            </View>
            <View style={styles.earningsContainer}>
              <View style={styles.earningsIcon}>
                <Ionicons name="wallet" size={20} color={COLORS.success} />
              </View>
              <View style={styles.earningsInfo}>
                <Text style={styles.earningsLabel}>Your Earnings</Text>
                <Text style={styles.earningsValue}>K {order.pricing.earnings}</Text>
              </View>
            </View>
          </View>
        </View>

        <View style={{ height: 100 }} />
      </ScrollView>

      {/* Action Buttons - Fixed at bottom */}
      <View style={styles.actionsContainer}>
        <TouchableOpacity style={styles.declineButton}>
          <Ionicons name="close-circle" size={22} color={COLORS.accent} />
          <Text style={styles.declineText}>Decline</Text>
        </TouchableOpacity>
        <TouchableOpacity 
          style={styles.acceptButton}
          onPress={handleAccept}
          disabled={isAccepting}
        >
          {isAccepting ? (
            <Text style={styles.acceptText}>Accepting...</Text>
          ) : (
            <>
              <Text style={styles.acceptText}>Accept Order</Text>
              <Ionicons name="checkmark-circle" size={22} color={COLORS.secondary} />
            </>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.background,
  },
  
  // Status Header
  statusHeader: { 
    backgroundColor: COLORS.secondary,
    padding: SPACING.lg,
    alignItems: 'center',
    borderBottomLeftRadius: RADIUS['2xl'],
    borderBottomRightRadius: RADIUS['2xl'],
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    gap: 6,
    marginBottom: SPACING.sm,
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
  },
  statusText: { 
    fontSize: TYPOGRAPHY.fontSize.sm,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    textTransform: 'capitalize',
  },
  orderNumber: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize['2xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    marginBottom: SPACING.md,
  },
  headerStats: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.secondaryLight,
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    gap: SPACING.md,
  },
  headerStat: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  headerStatValue: {
    color: COLORS.white,
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  headerStatDivider: {
    width: 1,
    height: 16,
    backgroundColor: COLORS.secondaryMuted,
  },
  
  // Sections
  section: { 
    padding: SPACING.base, 
    paddingBottom: 0,
  },
  sectionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    color: COLORS.textMuted,
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginBottom: SPACING.sm,
    marginLeft: SPACING.xs,
  },
  card: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl, 
    padding: SPACING.base,
    ...SHADOWS.sm,
  },
  
  // Customer
  customerRow: { 
    flexDirection: 'row', 
    alignItems: 'center',
  },
  customerAvatar: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  customerAvatarText: {
    fontSize: TYPOGRAPHY.fontSize.xl,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.secondary,
  },
  customerInfo: {
    flex: 1,
  },
  customerName: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  customerPhone: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textMuted,
    marginTop: 2,
  },
  customerActions: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  customerActionBtn: { 
    width: 42, 
    height: 42, 
    borderRadius: RADIUS.md, 
    justifyContent: 'center', 
    alignItems: 'center',
  },
  
  // Route
  routeContainer: {
    flexDirection: 'row',
    marginBottom: SPACING.base,
  },
  routeTimeline: {
    alignItems: 'center',
    marginRight: SPACING.md,
    paddingTop: 4,
  },
  routeDot: {
    width: 14,
    height: 14,
    borderRadius: 7,
  },
  routeLine: {
    width: 2,
    flex: 1,
    backgroundColor: COLORS.border,
    marginVertical: 4,
  },
  routeDetails: {
    flex: 1,
  },
  routePoint: {
    marginBottom: SPACING.sm,
  },
  routePointSpacing: {
    height: SPACING.lg,
  },
  routeLabel: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    letterSpacing: 0.5,
    marginBottom: 4,
  },
  routeAddress: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.text,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  navigateButton: { 
    flexDirection: 'row', 
    backgroundColor: COLORS.primary, 
    borderRadius: RADIUS.base, 
    padding: SPACING.md, 
    justifyContent: 'center', 
    alignItems: 'center',
    gap: SPACING.sm,
    ...SHADOWS.primary,
  },
  navigateText: { 
    color: COLORS.secondary, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    fontSize: TYPOGRAPHY.fontSize.base,
  },
  
  // Package
  packageHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.base,
    gap: SPACING.md,
  },
  packageIcon: {
    width: 44,
    height: 44,
    borderRadius: RADIUS.md,
    backgroundColor: COLORS.primary + '15',
    justifyContent: 'center',
    alignItems: 'center',
  },
  packageDesc: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    color: COLORS.text,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    flex: 1,
  },
  packageDetails: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
  },
  packageItem: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    gap: SPACING.sm,
  },
  packageItemIcon: {
    width: 32,
    height: 32,
    borderRadius: RADIUS.sm,
    justifyContent: 'center',
    alignItems: 'center',
  },
  packageItemLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
  },
  packageItemValue: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.text,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    textTransform: 'capitalize',
  },
  
  // Pricing
  pricingRow: { 
    flexDirection: 'row', 
    justifyContent: 'space-between', 
    marginBottom: SPACING.sm,
  },
  pricingLabel: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.textSecondary,
  },
  pricingValue: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.text,
  },
  pricingLabelBold: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  pricingValueBold: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  divider: { 
    height: 1, 
    backgroundColor: COLORS.border, 
    marginVertical: SPACING.sm,
  },
  earningsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.success + '10',
    padding: SPACING.md,
    borderRadius: RADIUS.md,
    marginTop: SPACING.sm,
    gap: SPACING.md,
  },
  earningsIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: COLORS.success + '20',
    justifyContent: 'center',
    alignItems: 'center',
  },
  earningsInfo: {},
  earningsLabel: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.success,
  },
  earningsValue: { 
    fontSize: TYPOGRAPHY.fontSize.xl, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.success,
  },
  
  // Actions
  actionsContainer: { 
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    flexDirection: 'row', 
    padding: SPACING.base,
    paddingBottom: SPACING.xl,
    gap: SPACING.md,
    backgroundColor: COLORS.surface,
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
    ...SHADOWS.lg,
  },
  declineButton: { 
    flex: 0.35,
    flexDirection: 'row',
    backgroundColor: COLORS.accent + '10', 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base, 
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 1, 
    borderColor: COLORS.accent + '30',
    gap: SPACING.sm,
  },
  declineText: { 
    color: COLORS.accent, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    fontSize: TYPOGRAPHY.fontSize.base,
  },
  acceptButton: { 
    flex: 0.65,
    flexDirection: 'row',
    backgroundColor: COLORS.primary, 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base, 
    alignItems: 'center',
    justifyContent: 'center',
    gap: SPACING.sm,
    ...SHADOWS.primary,
  },
  acceptText: { 
    color: COLORS.secondary, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    fontSize: TYPOGRAPHY.fontSize.lg,
  },
});
