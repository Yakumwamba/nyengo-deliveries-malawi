import React from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, Alert, Image } from 'react-native';
import { router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useAuthStore } from '../../src/stores/authStore';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

export default function ProfileScreen() {
  const { courier, logout } = useAuthStore();

  const handleLogout = () => {
    Alert.alert(
      'Logout',
      'Are you sure you want to logout?',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Logout', 
          style: 'destructive', 
          onPress: () => { 
            logout(); 
            router.replace('/(auth)/login'); 
          } 
        },
      ]
    );
  };

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      {/* Profile Header */}
      <View style={styles.header}>
        <View style={styles.headerTop}>
          <TouchableOpacity style={styles.settingsButton}>
            <Ionicons name="settings-outline" size={24} color={COLORS.white} />
          </TouchableOpacity>
        </View>
        
        <View style={styles.profileSection}>
          <View style={styles.avatarContainer}>
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>
                {courier?.companyName?.charAt(0) || 'N'}
              </Text>
            </View>
            <TouchableOpacity style={styles.editAvatarButton}>
              <Ionicons name="camera" size={14} color={COLORS.white} />
            </TouchableOpacity>
          </View>
          
          <Text style={styles.companyName}>
            {courier?.companyName || 'Nyengo Courier'}
          </Text>
          <Text style={styles.email}>
            {courier?.email || 'courier@nyengo.com'}
          </Text>
          
          <View style={styles.verifiedBadge}>
            <Ionicons name="shield-checkmark" size={14} color={COLORS.success} />
            <Text style={styles.verifiedText}>Verified Courier</Text>
          </View>
        </View>

        {/* Rating */}
        <View style={styles.ratingCard}>
          <View style={styles.ratingStars}>
            {[1, 2, 3, 4, 5].map((star) => (
              <Ionicons
                key={star}
                name={star <= 4 ? 'star' : 'star-half'}
                size={18}
                color={COLORS.primary}
              />
            ))}
          </View>
          <Text style={styles.ratingValue}>
            {courier?.rating?.toFixed(1) || '4.8'}
          </Text>
          <Text style={styles.ratingReviews}>
            ({courier?.totalReviews || 48} reviews)
          </Text>
        </View>
      </View>

      {/* Stats Cards */}
      <View style={styles.statsContainer}>
        <View style={styles.statsCard}>
          <View style={styles.stat}>
            <View style={[styles.statIcon, { backgroundColor: COLORS.primary + '15' }]}>
              <Ionicons name="cube" size={20} color={COLORS.primary} />
            </View>
            <Text style={styles.statValue}>{courier?.totalDeliveries || 156}</Text>
            <Text style={styles.statLabel}>Deliveries</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.stat}>
            <View style={[styles.statIcon, { backgroundColor: COLORS.success + '15' }]}>
              <Ionicons name="checkmark-circle" size={20} color={COLORS.success} />
            </View>
            <Text style={styles.statValue}>
              {((courier?.successRate || 98.5))?.toFixed(0)}%
            </Text>
            <Text style={styles.statLabel}>Success Rate</Text>
          </View>
          <View style={styles.statDivider} />
          <View style={styles.stat}>
            <View style={[styles.statIcon, { backgroundColor: COLORS.accent + '15' }]}>
              <Ionicons name="wallet" size={20} color={COLORS.accent} />
            </View>
            <Text style={styles.statValue}>
              K {courier?.walletBalance?.toFixed(0) || '2,450'}
            </Text>
            <Text style={styles.statLabel}>Balance</Text>
          </View>
        </View>
      </View>

      {/* Menu Sections */}
      <View style={styles.menuSection}>
        <Text style={styles.menuSectionTitle}>Account</Text>
        <View style={styles.menuCard}>
          <MenuItem 
            icon="person-outline" 
            title="Edit Profile" 
            subtitle="Update your personal info"
            iconColor={COLORS.primary}
          />
          <MenuItem 
            icon="card-outline" 
            title="Payment Methods" 
            subtitle="Manage your payout options"
            iconColor={COLORS.success}
          />
          <MenuItem 
            icon="document-text-outline" 
            title="Documents" 
            subtitle="ID, license & vehicle docs"
            iconColor={COLORS.info}
          />
        </View>
      </View>

      <View style={styles.menuSection}>
        <Text style={styles.menuSectionTitle}>Settings</Text>
        <View style={styles.menuCard}>
          <MenuItem 
            icon="notifications-outline" 
            title="Notifications" 
            subtitle="Push & email preferences"
            iconColor={COLORS.warning}
            showBadge
          />
          <MenuItem 
            icon="location-outline" 
            title="Service Areas" 
            subtitle="Manage delivery zones"
            iconColor={COLORS.accent}
          />
          <MenuItem 
            icon="shield-outline" 
            title="Privacy & Security" 
            subtitle="Password & data settings"
            iconColor={COLORS.secondary}
          />
        </View>
      </View>

      <View style={styles.menuSection}>
        <Text style={styles.menuSectionTitle}>Support</Text>
        <View style={styles.menuCard}>
          <MenuItem 
            icon="help-circle-outline" 
            title="Help Center" 
            subtitle="FAQs & guides"
            iconColor={COLORS.info}
          />
          <MenuItem 
            icon="chatbubble-outline" 
            title="Contact Support" 
            subtitle="Get help from our team"
            iconColor={COLORS.success}
          />
          <MenuItem 
            icon="information-circle-outline" 
            title="About Nyengo" 
            subtitle="Terms, privacy & more"
            iconColor={COLORS.textSecondary}
          />
        </View>
      </View>

      {/* Logout Button */}
      <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
        <Ionicons name="log-out-outline" size={22} color={COLORS.accent} />
        <Text style={styles.logoutText}>Logout</Text>
      </TouchableOpacity>

      <Text style={styles.version}>Version 1.0.0</Text>
      
      <View style={{ height: 32 }} />
    </ScrollView>
  );
}

function MenuItem({ 
  icon, 
  title, 
  subtitle,
  iconColor,
  showBadge,
}: { 
  icon: string; 
  title: string;
  subtitle: string;
  iconColor: string;
  showBadge?: boolean;
}) {
  return (
    <TouchableOpacity style={styles.menuItem}>
      <View style={[styles.menuItemIcon, { backgroundColor: iconColor + '15' }]}>
        <Ionicons name={icon as any} size={20} color={iconColor} />
      </View>
      <View style={styles.menuItemContent}>
        <Text style={styles.menuItemTitle}>{title}</Text>
        <Text style={styles.menuItemSubtitle}>{subtitle}</Text>
      </View>
      <View style={styles.menuItemRight}>
        {showBadge && (
          <View style={styles.menuItemBadge}>
            <Text style={styles.menuItemBadgeText}>3</Text>
          </View>
        )}
        <Ionicons name="chevron-forward" size={20} color={COLORS.textLight} />
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
    paddingTop: 50,
    paddingBottom: SPACING['3xl'],
    borderBottomLeftRadius: RADIUS['2xl'],
    borderBottomRightRadius: RADIUS['2xl'],
  },
  headerTop: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    marginBottom: SPACING.lg,
  },
  settingsButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: COLORS.secondaryLight,
    justifyContent: 'center',
    alignItems: 'center',
  },
  profileSection: {
    alignItems: 'center',
  },
  avatarContainer: {
    position: 'relative',
    marginBottom: SPACING.base,
  },
  avatar: { 
    width: 100, 
    height: 100, 
    borderRadius: 50, 
    backgroundColor: COLORS.primary, 
    justifyContent: 'center', 
    alignItems: 'center',
    borderWidth: 4,
    borderColor: COLORS.primaryDark,
  },
  avatarText: { 
    fontSize: TYPOGRAPHY.fontSize['4xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.secondary,
  },
  editAvatarButton: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: COLORS.accent,
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 3,
    borderColor: COLORS.secondary,
  },
  companyName: { 
    fontSize: TYPOGRAPHY.fontSize['2xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.white,
  },
  email: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.textMuted, 
    marginTop: 4,
  },
  verifiedBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.success + '20',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    marginTop: SPACING.md,
    gap: 6,
  },
  verifiedText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.success,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  ratingCard: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: SPACING.lg,
    gap: SPACING.sm,
  },
  ratingStars: {
    flexDirection: 'row',
    gap: 2,
  },
  ratingValue: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.white,
  },
  ratingReviews: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.textMuted,
  },
  
  // Stats
  statsContainer: {
    paddingHorizontal: SPACING.base,
    marginTop: -SPACING.xl,
  },
  statsCard: { 
    flexDirection: 'row', 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl, 
    padding: SPACING.lg,
    ...SHADOWS.lg,
  },
  stat: { 
    flex: 1, 
    alignItems: 'center',
  },
  statIcon: {
    width: 44,
    height: 44,
    borderRadius: RADIUS.md,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  statValue: { 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
  },
  statLabel: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted, 
    marginTop: 4,
  },
  statDivider: { 
    width: 1, 
    backgroundColor: COLORS.border,
    marginHorizontal: SPACING.sm,
  },
  
  // Menu Sections
  menuSection: {
    paddingHorizontal: SPACING.base,
    marginTop: SPACING.xl,
  },
  menuSectionTitle: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.textMuted,
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginBottom: SPACING.sm,
    marginLeft: SPACING.sm,
  },
  menuCard: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl,
    overflow: 'hidden',
    ...SHADOWS.sm,
  },
  menuItem: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    padding: SPACING.base,
    borderBottomWidth: 1, 
    borderBottomColor: COLORS.border,
  },
  menuItemIcon: {
    width: 40,
    height: 40,
    borderRadius: RADIUS.md,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  menuItemContent: {
    flex: 1,
  },
  menuItemTitle: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  menuItemSubtitle: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    color: COLORS.textMuted,
    marginTop: 2,
  },
  menuItemRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  menuItemBadge: {
    backgroundColor: COLORS.accent,
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: RADIUS.full,
  },
  menuItemBadgeText: {
    color: COLORS.white,
    fontSize: TYPOGRAPHY.fontSize.xs,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Logout
  logoutButton: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    justifyContent: 'center', 
    backgroundColor: COLORS.accent + '10', 
    marginHorizontal: SPACING.base,
    marginTop: SPACING.xl,
    padding: SPACING.base, 
    borderRadius: RADIUS.lg,
    borderWidth: 1,
    borderColor: COLORS.accent + '30',
    gap: SPACING.sm,
  },
  logoutText: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    color: COLORS.accent, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  version: { 
    textAlign: 'center', 
    color: COLORS.textMuted, 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    marginTop: SPACING.lg,
  },
});
