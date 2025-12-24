import React, { useState } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, Dimensions } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

const { width } = Dimensions.get('window');

export default function EarningsScreen() {
  const [selectedPeriod, setSelectedPeriod] = useState('week');

  const periods = [
    { key: 'today', label: 'Today' },
    { key: 'week', label: 'Week' },
    { key: 'month', label: 'Month' },
  ];

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      {/* Balance Card */}
      <View style={styles.balanceSection}>
        <View style={styles.balanceCard}>
          <View style={styles.balanceHeader}>
            <View>
              <Text style={styles.balanceLabel}>Available Balance</Text>
              <Text style={styles.balanceAmount}>K 2,450.00</Text>
            </View>
            <View style={styles.balanceBadge}>
              <Ionicons name="trending-up" size={14} color={COLORS.success} />
              <Text style={styles.balanceBadgeText}>+12%</Text>
            </View>
          </View>
          
          <View style={styles.balanceActions}>
            <TouchableOpacity style={styles.balanceActionPrimary}>
              <Ionicons name="arrow-down-circle" size={22} color={COLORS.secondary} />
              <Text style={styles.balanceActionPrimaryText}>Withdraw</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.balanceActionSecondary}>
              <Ionicons name="stats-chart" size={20} color={COLORS.white} />
              <Text style={styles.balanceActionSecondaryText}>History</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Decorative Element */}
        <View style={styles.balanceDecor}>
          <View style={styles.decorCircle1} />
          <View style={styles.decorCircle2} />
        </View>
      </View>

      {/* Period Selector */}
      <View style={styles.periodSection}>
        <View style={styles.periodSelector}>
          {periods.map((period) => (
            <TouchableOpacity
              key={period.key}
              style={[
                styles.periodButton,
                selectedPeriod === period.key && styles.periodButtonActive,
              ]}
              onPress={() => setSelectedPeriod(period.key)}
            >
              <Text
                style={[
                  styles.periodButtonText,
                  selectedPeriod === period.key && styles.periodButtonTextActive,
                ]}
              >
                {period.label}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Stats Grid */}
      <View style={styles.statsSection}>
        <View style={styles.statsGrid}>
          <StatCard
            icon="cash"
            title="Total Earnings"
            value="K 1,890"
            subtitle="45 orders"
            color={COLORS.primary}
          />
          <StatCard
            icon="bicycle"
            title="Distance"
            value="128 km"
            subtitle="This week"
            color={COLORS.accent}
          />
          <StatCard
            icon="time"
            title="Hours Online"
            value="32h"
            subtitle="Average 6h/day"
            color={COLORS.info}
          />
          <StatCard
            icon="star"
            title="Rating"
            value="4.9"
            subtitle="48 reviews"
            color={COLORS.warning}
          />
        </View>
      </View>

      {/* Earnings Breakdown */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>December Breakdown</Text>
        <View style={styles.breakdownCard}>
          <BreakdownRow 
            label="Delivery Earnings" 
            value="K 8,540.00" 
            icon="cube"
            iconBg={COLORS.primary}
          />
          <BreakdownRow 
            label="Tips Received" 
            value="K 320.00" 
            icon="heart"
            iconBg={COLORS.accent}
          />
          <BreakdownRow 
            label="Bonus & Incentives" 
            value="K 150.00" 
            icon="gift"
            iconBg={COLORS.success}
            isPositive
          />
          <View style={styles.divider} />
          <BreakdownRow 
            label="Platform Fees (10%)" 
            value="-K 854.00" 
            icon="remove-circle"
            iconBg={COLORS.danger}
            isNegative
          />
          <View style={styles.divider} />
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Net Earnings</Text>
            <Text style={styles.totalValue}>K 8,156.00</Text>
          </View>
        </View>
      </View>

      {/* Recent Transactions */}
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Recent Transactions</Text>
          <TouchableOpacity>
            <Text style={styles.seeAllText}>See All</Text>
          </TouchableOpacity>
        </View>
        
        <TransactionItem 
          type="earning" 
          title="Order Completed" 
          subtitle="NYG-ABC123"
          amount="+K 45.00" 
          time="2 hours ago"
          icon="checkmark-circle"
        />
        <TransactionItem 
          type="earning" 
          title="Order Completed" 
          subtitle="NYG-DEF456"
          amount="+K 32.00" 
          time="4 hours ago"
          icon="checkmark-circle"
        />
        <TransactionItem 
          type="tip" 
          title="Tip Received" 
          subtitle="From John D."
          amount="+K 10.00" 
          time="5 hours ago"
          icon="heart"
        />
        <TransactionItem 
          type="withdrawal" 
          title="Withdrawal" 
          subtitle="To Mobile Money"
          amount="-K 500.00" 
          time="Yesterday"
          icon="arrow-up-circle"
        />
        <TransactionItem 
          type="earning" 
          title="Order Completed" 
          subtitle="NYG-GHI789"
          amount="+K 58.00" 
          time="Yesterday"
          icon="checkmark-circle"
        />
      </View>

      <View style={{ height: 32 }} />
    </ScrollView>
  );
}

function StatCard({ icon, title, value, subtitle, color }: { 
  icon: string; 
  title: string; 
  value: string;
  subtitle: string;
  color: string;
}) {
  return (
    <View style={styles.statCard}>
      <View style={[styles.statIcon, { backgroundColor: color + '15' }]}>
        <Ionicons name={icon as any} size={20} color={color} />
      </View>
      <Text style={styles.statValue}>{value}</Text>
      <Text style={styles.statTitle}>{title}</Text>
      <Text style={styles.statSubtitle}>{subtitle}</Text>
    </View>
  );
}

function BreakdownRow({ label, value, icon, iconBg, isPositive, isNegative }: { 
  label: string; 
  value: string;
  icon: string;
  iconBg: string;
  isPositive?: boolean;
  isNegative?: boolean;
}) {
  return (
    <View style={styles.breakdownRow}>
      <View style={styles.breakdownLeft}>
        <View style={[styles.breakdownIcon, { backgroundColor: iconBg + '15' }]}>
          <Ionicons name={icon as any} size={16} color={iconBg} />
        </View>
        <Text style={styles.breakdownLabel}>{label}</Text>
      </View>
      <Text style={[
        styles.breakdownValue,
        isPositive && { color: COLORS.success },
        isNegative && { color: COLORS.danger },
      ]}>
        {value}
      </Text>
    </View>
  );
}

function TransactionItem({ 
  type, 
  title, 
  subtitle,
  amount, 
  time,
  icon,
}: { 
  type: string; 
  title: string;
  subtitle: string;
  amount: string; 
  time: string;
  icon: string;
}) {
  const getColors = () => {
    switch (type) {
      case 'earning': return { bg: COLORS.success + '15', icon: COLORS.success };
      case 'tip': return { bg: COLORS.accent + '15', icon: COLORS.accent };
      case 'withdrawal': return { bg: COLORS.warning + '15', icon: COLORS.warning };
      default: return { bg: COLORS.primary + '15', icon: COLORS.primary };
    }
  };
  
  const colors = getColors();
  const isPositive = !amount.startsWith('-');
  
  return (
    <TouchableOpacity style={styles.transactionItem}>
      <View style={[styles.transactionIcon, { backgroundColor: colors.bg }]}>
        <Ionicons name={icon as any} size={20} color={colors.icon} />
      </View>
      <View style={styles.transactionInfo}>
        <Text style={styles.transactionTitle}>{title}</Text>
        <Text style={styles.transactionSubtitle}>{subtitle}</Text>
      </View>
      <View style={styles.transactionRight}>
        <Text style={[
          styles.transactionAmount, 
          { color: isPositive ? COLORS.success : COLORS.danger }
        ]}>
          {amount}
        </Text>
        <Text style={styles.transactionTime}>{time}</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.background,
  },
  
  // Balance Section
  balanceSection: {
    backgroundColor: COLORS.secondary,
    paddingTop: 60,
    paddingBottom: SPACING['3xl'],
    paddingHorizontal: SPACING.lg,
    borderBottomLeftRadius: RADIUS['2xl'],
    borderBottomRightRadius: RADIUS['2xl'],
    position: 'relative',
    overflow: 'hidden',
  },
  balanceCard: {},
  balanceHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: SPACING.xl,
  },
  balanceLabel: { 
    color: COLORS.textMuted, 
    fontSize: TYPOGRAPHY.fontSize.base,
    marginBottom: 4,
  },
  balanceAmount: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize['5xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    letterSpacing: -1,
  },
  balanceBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.success + '20',
    paddingHorizontal: 10,
    paddingVertical: 6,
    borderRadius: RADIUS.full,
    gap: 4,
  },
  balanceBadgeText: {
    color: COLORS.success,
    fontSize: TYPOGRAPHY.fontSize.sm,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  balanceActions: { 
    flexDirection: 'row', 
    gap: SPACING.md,
  },
  balanceActionPrimary: { 
    flex: 1,
    flexDirection: 'row', 
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.primary,
    paddingVertical: SPACING.md,
    borderRadius: RADIUS.base,
    gap: 8,
    ...SHADOWS.primary,
  },
  balanceActionPrimaryText: { 
    color: COLORS.secondary, 
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  balanceActionSecondary: { 
    flex: 1,
    flexDirection: 'row', 
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: COLORS.secondaryLight,
    paddingVertical: SPACING.md,
    borderRadius: RADIUS.base,
    gap: 8,
  },
  balanceActionSecondaryText: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  balanceDecor: {
    position: 'absolute',
    right: -20,
    top: 20,
  },
  decorCircle1: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.primary + '10',
    position: 'absolute',
    right: 0,
    top: 0,
  },
  decorCircle2: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: COLORS.primary + '15',
    position: 'absolute',
    right: 60,
    top: 60,
  },
  
  // Period Selector
  periodSection: {
    paddingHorizontal: SPACING.lg,
    marginTop: -SPACING.lg,
  },
  periodSelector: {
    flexDirection: 'row',
    backgroundColor: COLORS.surface,
    borderRadius: RADIUS.lg,
    padding: SPACING.xs,
    ...SHADOWS.lg,
  },
  periodButton: {
    flex: 1,
    paddingVertical: SPACING.md,
    alignItems: 'center',
    borderRadius: RADIUS.base,
  },
  periodButtonActive: {
    backgroundColor: COLORS.primary,
  },
  periodButtonText: {
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
    color: COLORS.textSecondary,
  },
  periodButtonTextActive: {
    color: COLORS.secondary,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Stats Section
  statsSection: {
    paddingHorizontal: SPACING.base,
    marginTop: SPACING.xl,
  },
  statsGrid: { 
    flexDirection: 'row', 
    flexWrap: 'wrap',
    gap: SPACING.md,
  },
  statCard: { 
    width: (width - SPACING.base * 2 - SPACING.md) / 2,
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base,
    ...SHADOWS.sm,
  },
  statIcon: { 
    width: 40, 
    height: 40, 
    borderRadius: RADIUS.md, 
    justifyContent: 'center', 
    alignItems: 'center', 
    marginBottom: SPACING.md,
  },
  statValue: { 
    fontSize: TYPOGRAPHY.fontSize.xl, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  statTitle: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.textSecondary,
    marginTop: 2,
  },
  statSubtitle: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted,
    marginTop: 2,
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
    marginBottom: SPACING.md,
  },
  sectionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
    marginBottom: SPACING.md,
  },
  seeAllText: {
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.primary,
    marginBottom: SPACING.md,
  },
  
  // Breakdown Card
  breakdownCard: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl, 
    padding: SPACING.base,
    ...SHADOWS.md,
  },
  breakdownRow: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: SPACING.md,
  },
  breakdownLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
  },
  breakdownIcon: {
    width: 32,
    height: 32,
    borderRadius: RADIUS.sm,
    justifyContent: 'center',
    alignItems: 'center',
  },
  breakdownLabel: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.text,
  },
  breakdownValue: { 
    fontSize: TYPOGRAPHY.fontSize.base,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  divider: { 
    height: 1, 
    backgroundColor: COLORS.border, 
  },
  totalRow: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: SPACING.md,
  },
  totalLabel: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.text,
  },
  totalValue: { 
    fontSize: TYPOGRAPHY.fontSize.xl,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.success,
  },
  
  // Transactions
  transactionItem: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.lg, 
    padding: SPACING.base, 
    marginBottom: SPACING.sm,
    flexDirection: 'row', 
    alignItems: 'center',
    ...SHADOWS.sm,
  },
  transactionIcon: { 
    width: 44, 
    height: 44, 
    borderRadius: RADIUS.md, 
    justifyContent: 'center', 
    alignItems: 'center', 
    marginRight: SPACING.md,
  },
  transactionInfo: { 
    flex: 1,
  },
  transactionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    color: COLORS.text,
  },
  transactionSubtitle: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.textMuted,
    marginTop: 2,
  },
  transactionRight: {
    alignItems: 'flex-end',
  },
  transactionAmount: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  transactionTime: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted,
    marginTop: 2,
  },
});
