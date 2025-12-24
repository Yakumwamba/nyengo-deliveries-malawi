import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Dimensions } from 'react-native';
import { router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

const { width, height } = Dimensions.get('window');

export default function OnboardingScreen() {
  return (
    <View style={styles.container}>
      {/* Background Decoration */}
      <View style={styles.backgroundDecor}>
        <View style={styles.decorCircle1} />
        <View style={styles.decorCircle2} />
        <View style={styles.decorCircle3} />
      </View>

      {/* Main Content */}
      <View style={styles.content}>
        {/* Logo Section */}
        <View style={styles.logoSection}>
          <View style={styles.logoContainer}>
            <View style={styles.logoOuter}>
              <View style={styles.logoInner}>
                <Ionicons name="bicycle" size={48} color={COLORS.secondary} />
              </View>
            </View>
          </View>
          
          <Text style={styles.brandName}>Nyengo</Text>
          <Text style={styles.brandTagline}>Deliveries</Text>
        </View>

        {/* Hero Text */}
        <View style={styles.heroSection}>
          <Text style={styles.heroTitle}>
            Earn More,{'\n'}
            Deliver <Text style={styles.heroHighlight}>Faster</Text>
          </Text>
          <Text style={styles.heroDescription}>
            Join thousands of couriers making a difference. Track orders, connect with customers, and maximize your earnings.
          </Text>
        </View>

        {/* Features */}
        <View style={styles.features}>
          <FeatureItem 
            icon="trending-up" 
            title="Track Earnings" 
            description="Real-time insights"
            color={COLORS.primary}
          />
          <FeatureItem 
            icon="chatbubbles" 
            title="Live Chat" 
            description="Direct messaging"
            color={COLORS.accent}
          />
          <FeatureItem 
            icon="navigate" 
            title="Smart Routes" 
            description="Optimized paths"
            color={COLORS.success}
          />
        </View>
      </View>

      {/* Bottom Section */}
      <View style={styles.bottomSection}>
        <TouchableOpacity 
          style={styles.getStartedButton} 
          onPress={() => router.replace('/(auth)/login')}
          activeOpacity={0.9}
        >
          <Text style={styles.getStartedText}>Get Started</Text>
          <View style={styles.buttonArrow}>
            <Ionicons name="arrow-forward" size={20} color={COLORS.secondary} />
          </View>
        </TouchableOpacity>

        <TouchableOpacity 
          style={styles.loginLink}
          onPress={() => router.replace('/(auth)/login')}
        >
          <Text style={styles.loginLinkText}>
            Already have an account? <Text style={styles.loginLinkBold}>Login</Text>
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}

function FeatureItem({ 
  icon, 
  title, 
  description,
  color,
}: { 
  icon: string; 
  title: string;
  description: string;
  color: string;
}) {
  return (
    <View style={styles.featureItem}>
      <View style={[styles.featureIcon, { backgroundColor: color }]}>
        <Ionicons name={icon as any} size={22} color={COLORS.white} />
      </View>
      <Text style={styles.featureTitle}>{title}</Text>
      <Text style={styles.featureDescription}>{description}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.secondary,
  },
  
  // Background Decoration
  backgroundDecor: {
    position: 'absolute',
    width: '100%',
    height: '100%',
  },
  decorCircle1: {
    position: 'absolute',
    width: 300,
    height: 300,
    borderRadius: 150,
    backgroundColor: COLORS.primary + '08',
    top: -100,
    right: -100,
  },
  decorCircle2: {
    position: 'absolute',
    width: 200,
    height: 200,
    borderRadius: 100,
    backgroundColor: COLORS.primary + '05',
    top: 150,
    left: -50,
  },
  decorCircle3: {
    position: 'absolute',
    width: 150,
    height: 150,
    borderRadius: 75,
    backgroundColor: COLORS.accent + '08',
    bottom: 200,
    right: -30,
  },
  
  // Content
  content: { 
    flex: 1,
    paddingHorizontal: SPACING.xl,
    paddingTop: 80,
  },
  
  // Logo
  logoSection: {
    alignItems: 'center',
    marginBottom: SPACING['3xl'],
  },
  logoContainer: {
    marginBottom: SPACING.lg,
  },
  logoOuter: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.primary,
    justifyContent: 'center',
    alignItems: 'center',
    ...SHADOWS.primary,
  },
  logoInner: {
    width: 70,
    height: 70,
    borderRadius: 35,
    backgroundColor: COLORS.white,
    justifyContent: 'center',
    alignItems: 'center',
  },
  brandName: {
    fontSize: TYPOGRAPHY.fontSize['4xl'],
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.primary,
    letterSpacing: -1,
  },
  brandTagline: {
    fontSize: TYPOGRAPHY.fontSize.xl,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
    color: COLORS.textMuted,
    marginTop: -4,
  },
  
  // Hero
  heroSection: {
    marginBottom: SPACING['2xl'],
  },
  heroTitle: { 
    fontSize: TYPOGRAPHY.fontSize['5xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.white,
    lineHeight: 48,
    letterSpacing: -1,
  },
  heroHighlight: {
    color: COLORS.primary,
  },
  heroDescription: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    color: COLORS.textMuted, 
    marginTop: SPACING.md,
    lineHeight: 24,
  },
  
  // Features
  features: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
  },
  featureItem: { 
    alignItems: 'center',
    flex: 1,
  },
  featureIcon: { 
    width: 52, 
    height: 52, 
    borderRadius: RADIUS.lg, 
    justifyContent: 'center', 
    alignItems: 'center', 
    marginBottom: SPACING.sm,
  },
  featureTitle: { 
    fontSize: TYPOGRAPHY.fontSize.sm, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.white,
    textAlign: 'center',
    marginBottom: 2,
  },
  featureDescription: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    textAlign: 'center',
  },
  
  // Bottom Section
  bottomSection: {
    paddingHorizontal: SPACING.xl,
    paddingBottom: SPACING['3xl'],
  },
  getStartedButton: { 
    backgroundColor: COLORS.primary, 
    borderRadius: RADIUS.xl, 
    height: 60, 
    flexDirection: 'row', 
    justifyContent: 'center', 
    alignItems: 'center',
    marginBottom: SPACING.lg,
    ...SHADOWS.primary,
  },
  getStartedText: { 
    color: COLORS.secondary, 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    marginRight: SPACING.md,
  },
  buttonArrow: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: COLORS.secondary + '20',
    justifyContent: 'center',
    alignItems: 'center',
  },
  loginLink: {
    alignItems: 'center',
  },
  loginLinkText: {
    color: COLORS.textMuted,
    fontSize: TYPOGRAPHY.fontSize.base,
  },
  loginLinkBold: {
    color: COLORS.primary,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
});
