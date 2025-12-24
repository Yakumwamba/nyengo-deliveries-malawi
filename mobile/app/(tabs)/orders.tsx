import React, { useState } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, TextInput } from 'react-native';
import { router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS, STATUS_COLORS } from '../../src/constants/theme';

const mockOrders = [
  { id: '1', orderNumber: 'NYG-20231215-ABC123', status: 'pending', customer: 'John Doe', amount: 45, pickup: 'Cairo Road', delivery: 'Northmead', time: '10 min ago' },
  { id: '2', orderNumber: 'NYG-20231215-DEF456', status: 'accepted', customer: 'Jane Smith', amount: 32, pickup: 'Manda Hill', delivery: 'Kabulonga', time: '25 min ago' },
  { id: '3', orderNumber: 'NYG-20231215-GHI789', status: 'in_transit', customer: 'Bob Wilson', amount: 58, pickup: 'East Park Mall', delivery: 'Roma', time: '45 min ago' },
  { id: '4', orderNumber: 'NYG-20231214-JKL012', status: 'delivered', customer: 'Alice Brown', amount: 75, pickup: 'Levy Mall', delivery: 'Ibex Hill', time: '2 hours ago' },
  { id: '5', orderNumber: 'NYG-20231214-MNO345', status: 'cancelled', customer: 'Charlie Davis', amount: 28, pickup: 'Arcades', delivery: 'Chelstone', time: '3 hours ago' },
];

const filterOptions = [
  { key: 'all', label: 'All', icon: 'apps' },
  { key: 'pending', label: 'Pending', icon: 'time' },
  { key: 'in_transit', label: 'Active', icon: 'bicycle' },
  { key: 'delivered', label: 'Done', icon: 'checkmark-circle' },
];

export default function OrdersScreen() {
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('all');
  const [isSearchFocused, setIsSearchFocused] = useState(false);

  const filteredOrders = mockOrders.filter((order) => {
    if (filter !== 'all' && order.status !== filter) return false;
    if (search && !order.orderNumber.toLowerCase().includes(search.toLowerCase()) && !order.customer.toLowerCase().includes(search.toLowerCase())) return false;
    return true;
  });

  const getOrderCount = (status: string) => {
    if (status === 'all') return mockOrders.length;
    if (status === 'in_transit') return mockOrders.filter(o => o.status === 'in_transit' || o.status === 'accepted').length;
    return mockOrders.filter(o => o.status === status).length;
  };

  return (
    <View style={styles.container}>
      {/* Search Bar */}
      <View style={styles.searchSection}>
        <View style={[styles.searchContainer, isSearchFocused && styles.searchContainerFocused]}>
          <Ionicons name="search" size={20} color={isSearchFocused ? COLORS.primary : COLORS.textMuted} />
          <TextInput 
            style={styles.searchInput} 
            placeholder="Search orders, customers..." 
            value={search} 
            onChangeText={setSearch}
            onFocus={() => setIsSearchFocused(true)}
            onBlur={() => setIsSearchFocused(false)}
            placeholderTextColor={COLORS.textMuted}
          />
          {search.length > 0 && (
            <TouchableOpacity onPress={() => setSearch('')}>
              <Ionicons name="close-circle" size={20} color={COLORS.textMuted} />
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Filter Tabs */}
      <View style={styles.filtersContainer}>
        <View style={styles.filters}>
          {filterOptions.map((f) => (
            <TouchableOpacity 
              key={f.key} 
              style={[styles.filterButton, filter === f.key && styles.filterActive]} 
              onPress={() => setFilter(f.key)}
            >
              <Ionicons 
                name={f.icon as any} 
                size={16} 
                color={filter === f.key ? COLORS.secondary : COLORS.textSecondary} 
              />
              <Text style={[styles.filterText, filter === f.key && styles.filterTextActive]}>
                {f.label}
              </Text>
              <View style={[styles.filterBadge, filter === f.key && styles.filterBadgeActive]}>
                <Text style={[styles.filterBadgeText, filter === f.key && styles.filterBadgeTextActive]}>
                  {getOrderCount(f.key)}
                </Text>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Orders List */}
      <FlatList
        data={filteredOrders}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <OrderCard order={item} />}
        contentContainerStyle={styles.list}
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIcon}>
              <Ionicons name="receipt-outline" size={48} color={COLORS.textMuted} />
            </View>
            <Text style={styles.emptyTitle}>No orders found</Text>
            <Text style={styles.emptyText}>Try adjusting your search or filters</Text>
          </View>
        }
      />
    </View>
  );
}

function OrderCard({ order }: { order: typeof mockOrders[0] }) {
  const statusConfig = STATUS_COLORS[order.status as keyof typeof STATUS_COLORS] || STATUS_COLORS.pending;
  
  return (
    <TouchableOpacity 
      style={styles.orderCard} 
      onPress={() => router.push(`/orders/${order.id}`)}
      activeOpacity={0.7}
    >
      {/* Header */}
      <View style={styles.orderHeader}>
        <View style={styles.orderHeaderLeft}>
          <Text style={styles.orderNumber}>{order.orderNumber}</Text>
          <Text style={styles.orderTime}>{order.time}</Text>
        </View>
        <View style={[styles.statusBadge, { backgroundColor: statusConfig.bg }]}>
          <View style={[styles.statusDot, { backgroundColor: statusConfig.icon }]} />
          <Text style={[styles.statusText, { color: statusConfig.text }]}>
            {order.status.replace('_', ' ')}
          </Text>
        </View>
      </View>

      {/* Customer */}
      <View style={styles.customerRow}>
        <View style={styles.customerAvatar}>
          <Text style={styles.customerAvatarText}>{order.customer.charAt(0)}</Text>
        </View>
        <View style={styles.customerInfo}>
          <Text style={styles.customerName}>{order.customer}</Text>
          <Text style={styles.customerLabel}>Customer</Text>
        </View>
      </View>

      {/* Route */}
      <View style={styles.routeContainer}>
        <View style={styles.routeTimeline}>
          <View style={[styles.routeDot, { backgroundColor: COLORS.success }]} />
          <View style={styles.routeLine} />
          <View style={[styles.routeDot, { backgroundColor: COLORS.accent }]} />
        </View>
        <View style={styles.routeDetails}>
          <View style={styles.routePoint}>
            <Text style={styles.routeLabel}>Pickup</Text>
            <Text style={styles.routeAddress}>{order.pickup}</Text>
          </View>
          <View style={styles.routePoint}>
            <Text style={styles.routeLabel}>Delivery</Text>
            <Text style={styles.routeAddress}>{order.delivery}</Text>
          </View>
        </View>
      </View>

      {/* Footer */}
      <View style={styles.orderFooter}>
        <View style={styles.amountContainer}>
          <Text style={styles.amountLabel}>Earnings</Text>
          <Text style={styles.amount}>K {order.amount}</Text>
        </View>
        <TouchableOpacity style={styles.viewButton}>
          <Text style={styles.viewButtonText}>View Details</Text>
          <Ionicons name="arrow-forward" size={16} color={COLORS.secondary} />
        </TouchableOpacity>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.background,
  },
  
  // Search
  searchSection: {
    paddingHorizontal: SPACING.base,
    paddingTop: SPACING.base,
    paddingBottom: SPACING.sm,
    backgroundColor: COLORS.secondary,
  },
  searchContainer: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    backgroundColor: COLORS.secondaryLight, 
    paddingHorizontal: SPACING.base,
    borderRadius: RADIUS.base,
    height: 48,
    gap: SPACING.sm,
  },
  searchContainerFocused: {
    backgroundColor: COLORS.secondaryMuted,
    borderWidth: 1,
    borderColor: COLORS.primary,
  },
  searchInput: { 
    flex: 1, 
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.white,
  },
  
  // Filters
  filtersContainer: {
    backgroundColor: COLORS.secondary,
    paddingBottom: SPACING.base,
    borderBottomLeftRadius: RADIUS.xl,
    borderBottomRightRadius: RADIUS.xl,
  },
  filters: { 
    flexDirection: 'row', 
    paddingHorizontal: SPACING.base,
    gap: SPACING.sm,
  },
  filterButton: { 
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.sm,
    paddingHorizontal: SPACING.md,
    borderRadius: RADIUS.base, 
    backgroundColor: COLORS.secondaryLight,
    gap: 4,
  },
  filterActive: { 
    backgroundColor: COLORS.primary,
  },
  filterText: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  filterTextActive: { 
    color: COLORS.secondary,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  filterBadge: {
    backgroundColor: COLORS.secondaryMuted,
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: RADIUS.full,
    marginLeft: 2,
  },
  filterBadgeActive: {
    backgroundColor: COLORS.secondary + '30',
  },
  filterBadgeText: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  filterBadgeTextActive: {
    color: COLORS.secondary,
  },
  
  // List
  list: { 
    padding: SPACING.base,
    paddingTop: SPACING.lg,
  },
  
  // Order Card
  orderCard: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl, 
    padding: SPACING.base,
    marginBottom: SPACING.md,
    ...SHADOWS.md,
  },
  orderHeader: { 
    flexDirection: 'row', 
    justifyContent: 'space-between', 
    alignItems: 'center', 
    marginBottom: SPACING.md,
  },
  orderHeaderLeft: {},
  orderNumber: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  orderTime: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    marginTop: 2,
  },
  statusBadge: { 
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: SPACING.md, 
    paddingVertical: SPACING.sm, 
    borderRadius: RADIUS.full,
    gap: 6,
  },
  statusDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
  },
  statusText: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    textTransform: 'capitalize',
  },
  
  // Customer
  customerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.md,
    paddingBottom: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  customerAvatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  customerAvatarText: {
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.secondary,
  },
  customerInfo: {},
  customerName: {
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  customerLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
  },
  
  // Route
  routeContainer: {
    flexDirection: 'row',
    marginBottom: SPACING.md,
  },
  routeTimeline: {
    alignItems: 'center',
    marginRight: SPACING.md,
    paddingVertical: 4,
  },
  routeDot: {
    width: 10,
    height: 10,
    borderRadius: 5,
  },
  routeLine: {
    width: 2,
    flex: 1,
    backgroundColor: COLORS.border,
    marginVertical: 4,
  },
  routeDetails: {
    flex: 1,
    justifyContent: 'space-between',
  },
  routePoint: {
    marginBottom: SPACING.sm,
  },
  routeLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginBottom: 2,
  },
  routeAddress: {
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.text,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  
  // Footer
  orderFooter: { 
    flexDirection: 'row', 
    justifyContent: 'space-between', 
    alignItems: 'center',
    paddingTop: SPACING.md,
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
  },
  amountContainer: {},
  amountLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
  },
  amount: { 
    fontSize: TYPOGRAPHY.fontSize.xl, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  viewButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.base,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.base,
    gap: 6,
  },
  viewButtonText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.secondary,
  },
  
  // Empty
  emptyContainer: {
    alignItems: 'center',
    paddingTop: SPACING['4xl'],
  },
  emptyIcon: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: COLORS.surface,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.base,
    ...SHADOWS.sm,
  },
  emptyTitle: {
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
    marginBottom: SPACING.sm,
  },
  emptyText: { 
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.textMuted,
    textAlign: 'center',
  },
});
