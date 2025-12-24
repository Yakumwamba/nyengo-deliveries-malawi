import React from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAuthStore } from '../../src/stores/authStore';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS, STATUS_COLORS } from '../../src/constants/theme';

export default function DashboardScreen() {
  const { courier } = useAuthStore();
  const [refreshing, setRefreshing] = React.useState(false);

  const onRefresh = React.useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1000);
  }, []);

  // Get time-based greeting
  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good Morning';
    if (hour < 17) return 'Good Afternoon';
    return 'Good Evening';
  };

  return (
    <ScrollView
      style={styles.container}
      refreshControl={
        <RefreshControl 
          refreshing={refreshing} 
          onRefresh={onRefresh} 
          colors={[COLORS.primary]}
          tintColor={COLORS.primary}
          progressBackgroundColor={COLORS.secondary}
        />
      }
    >
      {/* Header Section */}
      <View style={styles.header}>
        <View style={styles.headerContent}>
          <View style={styles.headerLeft}>
            <Text style={styles.greeting}>{getGreeting()} 👋</Text>
            <Text style={styles.companyName}>{courier?.companyName || 'Nyengo Courier'}</Text>
          </View>
          <TouchableOpacity style={styles.notificationBtn}>
            <Ionicons name="notifications-outline" size={24} color={COLORS.white} />
            <View style={styles.notificationDot} />
          </TouchableOpacity>
        </View>
        
        {/* Status Badge */}
        <TouchableOpacity style={styles.statusBadge}>
          <View style={styles.statusDot} />
          <Text style={styles.statusText}>Online</Text>
          <Ionicons name="chevron-down" size={16} color={COLORS.primary} />
        </TouchableOpacity>
      </View>

      {/* Stats Cards Grid */}
      <View style={styles.statsContainer}>
        <View style={styles.statsGrid}>
          <StatCard 
            icon="receipt" 
            title="Today's Orders" 
            value="12" 
            trend="+3" 
            color={COLORS.primary}
          />
          <StatCard 
            icon="checkmark-circle" 
            title="Completed" 
            value="8" 
            trend="+2" 
            color={COLORS.success}
          />
          <StatCard 
            icon="time" 
            title="Pending" 
            value="4" 
            color={COLORS.warning}
          />
          <StatCard 
            icon="wallet" 
            title="Earnings" 
            value="K 450" 
            trend="+12%" 
            color={COLORS.accent}
          />
        </View>
      </View>

      {/* Quick Actions */}
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Quick Actions</Text>
        </View>
        <View style={styles.actionsGrid}>
          <ActionButton 
            icon="add-circle" 
            title="New Order" 
            bgColor={COLORS.primary}
            iconColor={COLORS.secondary}
          />
          <ActionButton 
            icon="qr-code" 
            title="Scan QR" 
            bgColor={COLORS.secondary}
            iconColor={COLORS.primary}
          />
          <ActionButton 
            icon="navigate" 
            title="Navigate" 
            bgColor={COLORS.accent}
            iconColor={COLORS.white}
          />
          <ActionButton 
            icon="call" 
            title="Support" 
            bgColor={COLORS.secondaryLight}
            iconColor={COLORS.primary}
          />
        </View>
      </View>

      {/* Active Delivery Banner */}
      <View style={styles.section}>
        <TouchableOpacity style={styles.activeDeliveryBanner}>
          <View style={styles.activeDeliveryLeft}>
            <View style={styles.activeDeliveryIcon}>
              <Ionicons name="bicycle" size={28} color={COLORS.secondary} />
            </View>
            <View style={styles.activeDeliveryInfo}>
              <Text style={styles.activeDeliveryLabel}>Active Delivery</Text>
              <Text style={styles.activeDeliveryId}>NYG-20231215-ABC123</Text>
            </View>
          </View>
          <View style={styles.activeDeliveryRight}>
            <Text style={styles.activeDeliveryEta}>15 min</Text>
            <Ionicons name="arrow-forward-circle" size={32} color={COLORS.secondary} />
          </View>
        </TouchableOpacity>
      </View>

      {/* Recent Orders */}
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Recent Orders</Text>
          <TouchableOpacity>
            <Text style={styles.seeAllText}>See All</Text>
          </TouchableOpacity>
        </View>
        <OrderItem 
          orderNumber="NYG-20231215-ABC123" 
          status="in_transit" 
          customer="John Doe" 
          amount="K 45"
          pickup="Cairo Road"
          delivery="Northmead"
        />
        <OrderItem 
          orderNumber="NYG-20231215-DEF456" 
          status="pending" 
          customer="Jane Smith" 
          amount="K 32"
          pickup="Manda Hill"
          delivery="Kabulonga"
        />
        <OrderItem 
          orderNumber="NYG-20231215-GHI789" 
          status="delivered" 
          customer="Bob Wilson" 
          amount="K 58"
          pickup="East Park Mall"
          delivery="Roma"
        />
      </View>

      <View style={{ height: 24 }} />
    </ScrollView>
  );
}

function StatCard({ icon, title, value, trend, color }: { 
  icon: string; 
  title: string; 
  value: string; 
  trend?: string;
  color: string;
}) {
  return (
    <View style={styles.statCard}>
      <View style={[styles.statIcon, { backgroundColor: color + '15' }]}>
        <Ionicons name={icon as any} size={22} color={color} />
      </View>
      <Text style={styles.statValue}>{value}</Text>
      <Text style={styles.statTitle}>{title}</Text>
      {trend && (
        <View style={[styles.trendBadge, { backgroundColor: COLORS.success + '15' }]}>
          <Ionicons name="trending-up" size={12} color={COLORS.success} />
          <Text style={[styles.trendText, { color: COLORS.success }]}>{trend}</Text>
        </View>
      )}
    </View>
  );
}

function ActionButton({ icon, title, bgColor, iconColor }: { 
  icon: string; 
  title: string; 
  bgColor: string;
  iconColor: string;
}) {
  return (
    <TouchableOpacity style={styles.actionButton}>
      <View style={[styles.actionIcon, { backgroundColor: bgColor }]}>
        <Ionicons name={icon as any} size={24} color={iconColor} />
      </View>
      <Text style={styles.actionTitle}>{title}</Text>
    </TouchableOpacity>
  );
}

function OrderItem({ 
  orderNumber, 
  status, 
  customer, 
  amount,
  pickup,
  delivery,
}: { 
  orderNumber: string; 
  status: string; 
  customer: string; 
  amount: string;
  pickup: string;
  delivery: string;
}) {
  const statusConfig = STATUS_COLORS[status as keyof typeof STATUS_COLORS] || STATUS_COLORS.pending;
  
  return (
    <TouchableOpacity style={styles.orderItem}>
      <View style={styles.orderLeft}>
        <View style={[styles.orderIcon, { backgroundColor: statusConfig.bg }]}>
          <Ionicons name="cube" size={20} color={statusConfig.icon} />
        </View>
      </View>
      <View style={styles.orderInfo}>
        <View style={styles.orderHeader}>
          <Text style={styles.orderNumber}>{orderNumber}</Text>
          <View style={[styles.orderStatus, { backgroundColor: statusConfig.bg }]}>
            <Text style={[styles.orderStatusText, { color: statusConfig.text }]}>
              {status.replace('_', ' ')}
            </Text>
          </View>
        </View>
        <Text style={styles.orderCustomer}>{customer}</Text>
        <View style={styles.orderRoute}>
          <Ionicons name="ellipse" size={8} color={COLORS.success} />
          <Text style={styles.orderRouteText}>{pickup}</Text>
          <Ionicons name="arrow-forward" size={12} color={COLORS.textMuted} />
          <Ionicons name="location" size={10} color={COLORS.accent} />
          <Text style={styles.orderRouteText}>{delivery}</Text>
        </View>
      </View>
      <View style={styles.orderRight}>
        <Text style={styles.orderAmount}>{amount}</Text>
        <Ionicons name="chevron-forward" size={20} color={COLORS.textMuted} />
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.background,
  },
  
  // Header
  header: { 
    backgroundColor: COLORS.secondary, 
    paddingHorizontal: SPACING.lg,
    paddingTop: 60,
    paddingBottom: SPACING['2xl'],
    borderBottomLeftRadius: RADIUS['2xl'],
    borderBottomRightRadius: RADIUS['2xl'],
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  headerLeft: {},
  greeting: { 
    color: COLORS.textMuted, 
    fontSize: TYPOGRAPHY.fontSize.base,
    marginBottom: 4,
  },
  companyName: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize['2xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    letterSpacing: -0.5,
  },
  notificationBtn: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: COLORS.secondaryLight,
    justifyContent: 'center',
    alignItems: 'center',
  },
  notificationDot: {
    position: 'absolute',
    top: 10,
    right: 10,
    width: 10,
    height: 10,
    borderRadius: 5,
    backgroundColor: COLORS.accent,
    borderWidth: 2,
    borderColor: COLORS.secondary,
  },
  statusBadge: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    alignSelf: 'flex-start',
    backgroundColor: COLORS.secondaryLight, 
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    marginTop: SPACING.base,
    gap: 6,
  },
  statusDot: { 
    width: 8, 
    height: 8, 
    borderRadius: 4, 
    backgroundColor: COLORS.success,
  },
  statusText: { 
    color: COLORS.primary, 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  
  // Stats
  statsContainer: {
    marginTop: -SPACING.xl,
    paddingHorizontal: SPACING.base,
  },
  statsGrid: { 
    flexDirection: 'row', 
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  statCard: { 
    width: '48%',
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base,
    ...SHADOWS.md,
  },
  statIcon: { 
    width: 44, 
    height: 44, 
    borderRadius: RADIUS.md, 
    justifyContent: 'center', 
    alignItems: 'center', 
    marginBottom: SPACING.md,
  },
  statValue: { 
    fontSize: TYPOGRAPHY.fontSize['2xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
    marginBottom: 2,
  },
  statTitle: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textSecondary,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  trendBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: RADIUS.sm,
    marginTop: SPACING.sm,
    gap: 2,
  },
  trendText: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  
  // Sections
  section: { 
    paddingHorizontal: SPACING.base,
    marginTop: SPACING.xl,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.base,
  },
  sectionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  seeAllText: {
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.primary,
  },
  
  // Actions
  actionsGrid: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
  },
  actionButton: { 
    alignItems: 'center',
    width: '22%',
  },
  actionIcon: { 
    width: 56, 
    height: 56, 
    borderRadius: RADIUS.lg, 
    justifyContent: 'center', 
    alignItems: 'center', 
    marginBottom: SPACING.sm,
    ...SHADOWS.md,
  },
  actionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textSecondary,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
    textAlign: 'center',
  },
  
  // Active Delivery Banner
  activeDeliveryBanner: {
    backgroundColor: COLORS.primary,
    borderRadius: RADIUS.lg,
    padding: SPACING.base,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    ...SHADOWS.primary,
  },
  activeDeliveryLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
  },
  activeDeliveryIcon: {
    width: 48,
    height: 48,
    borderRadius: RADIUS.md,
    backgroundColor: COLORS.white,
    justifyContent: 'center',
    alignItems: 'center',
  },
  activeDeliveryInfo: {},
  activeDeliveryLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.secondary,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    opacity: 0.8,
  },
  activeDeliveryId: {
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.secondary,
  },
  activeDeliveryRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  activeDeliveryEta: {
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.secondary,
  },
  
  // Order Items
  orderItem: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base, 
    marginBottom: SPACING.md,
    flexDirection: 'row', 
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  orderLeft: {
    marginRight: SPACING.md,
  },
  orderIcon: {
    width: 44,
    height: 44,
    borderRadius: RADIUS.md,
    justifyContent: 'center',
    alignItems: 'center',
  },
  orderInfo: { 
    flex: 1,
  },
  orderHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
    marginBottom: 4,
  },
  orderNumber: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    color: COLORS.text,
  },
  orderCustomer: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.textSecondary,
    marginBottom: 4,
  },
  orderRoute: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  orderRouteText: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
  },
  orderRight: { 
    alignItems: 'flex-end',
  },
  orderAmount: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
    marginBottom: 4,
  },
  orderStatus: { 
    paddingHorizontal: 8, 
    paddingVertical: 3, 
    borderRadius: RADIUS.sm,
  },
  orderStatusText: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    textTransform: 'capitalize',
  },
});
