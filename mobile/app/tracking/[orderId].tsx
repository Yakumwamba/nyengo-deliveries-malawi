import React, { useEffect, useState, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Alert,
  Linking,
  Platform,
} from 'react-native';
import MapView, { Marker, Polyline, PROVIDER_GOOGLE } from 'react-native-maps';
import { useLocalSearchParams, router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useTrackingWatch } from '../../src/hooks/useTracking';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

export default function TrackingScreen() {
  const { orderId } = useLocalSearchParams<{ orderId: string }>();
  const { trackingData, locationHistory, isConnected, error } = useTrackingWatch(orderId);
  const mapRef = useRef<MapView>(null);
  const [region, setRegion] = useState({
    latitude: -15.4167,
    longitude: 28.2833,
    latitudeDelta: 0.05,
    longitudeDelta: 0.05,
  });

  // Update map when tracking data changes
  useEffect(() => {
    if (trackingData?.latitude && trackingData?.longitude) {
      const newRegion = {
        latitude: trackingData.latitude,
        longitude: trackingData.longitude,
        latitudeDelta: 0.01,
        longitudeDelta: 0.01,
      };
      setRegion(newRegion);
      mapRef.current?.animateToRegion(newRegion, 500);
    }
  }, [trackingData?.latitude, trackingData?.longitude]);

  const handleCallDriver = () => {
    if (trackingData?.driverPhone) {
      Linking.openURL(`tel:${trackingData.driverPhone}`);
    }
  };

  const handleWhatsApp = () => {
    if (trackingData?.driverPhone) {
      const phone = trackingData.driverPhone.replace(/\+/g, '');
      Linking.openURL(`https://wa.me/${phone}`);
    }
  };

  if (!trackingData) {
    return (
      <View style={styles.loadingContainer}>
        <View style={styles.loadingIcon}>
          <Ionicons name="location-outline" size={48} color={COLORS.primary} />
        </View>
        <Text style={styles.loadingTitle}>Loading tracking data...</Text>
        <Text style={styles.loadingText}>Please wait while we connect</Text>
        {error && <Text style={styles.errorText}>{error}</Text>}
      </View>
    );
  }

  if (!trackingData.isActive) {
    return (
      <View style={styles.loadingContainer}>
        <View style={styles.successIcon}>
          <Ionicons name="checkmark" size={48} color={COLORS.white} />
        </View>
        <Text style={styles.successTitle}>Delivery Completed!</Text>
        <Text style={styles.loadingText}>Your package has been delivered</Text>
        <TouchableOpacity style={styles.backButton} onPress={() => router.back()}>
          <Text style={styles.backButtonText}>Go Back</Text>
          <Ionicons name="arrow-forward" size={20} color={COLORS.secondary} />
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* Map View */}
      <MapView
        ref={mapRef}
        style={styles.map}
        provider={Platform.OS === 'android' ? PROVIDER_GOOGLE : undefined}
        region={region}
        showsUserLocation
        showsMyLocationButton
      >
        {/* Driver marker */}
        {trackingData.latitude !== 0 && (
          <Marker
            coordinate={{
              latitude: trackingData.latitude,
              longitude: trackingData.longitude,
            }}
            title={trackingData.driverName}
            description={`${trackingData.vehicleType} • ${trackingData.etaMinutes} min away`}
          >
            <View style={styles.driverMarker}>
              <Ionicons
                name={getVehicleIcon(trackingData.vehicleType)}
                size={22}
                color={COLORS.secondary}
              />
            </View>
          </Marker>
        )}

        {/* Route polyline from history */}
        {locationHistory.length > 1 && (
          <Polyline
            coordinates={locationHistory.map((p) => ({
              latitude: p.latitude,
              longitude: p.longitude,
            }))}
            strokeColor={COLORS.primary}
            strokeWidth={4}
          />
        )}
      </MapView>

      {/* Connection Status */}
      <View style={styles.connectionBadge}>
        <View style={[styles.connectionDot, { backgroundColor: isConnected ? COLORS.success : COLORS.accent }]} />
        <Text style={styles.connectionText}>
          {isConnected ? 'Live Tracking' : 'Connecting...'}
        </Text>
      </View>

      {/* ETA Card */}
      <View style={styles.etaCard}>
        <View style={styles.etaMain}>
          <Text style={styles.etaLabel}>ARRIVING IN</Text>
          <Text style={styles.etaTime}>{trackingData.etaMinutes}</Text>
          <Text style={styles.etaUnit}>minutes</Text>
        </View>
        <View style={styles.divider} />
        <View style={styles.etaDetails}>
          <View style={styles.etaDetailItem}>
            <Ionicons name="navigate" size={16} color={COLORS.primary} />
            <Text style={styles.etaDetailText}>
              {trackingData.distanceRemaining.toFixed(1)} km
            </Text>
          </View>
          <View style={styles.etaDetailItem}>
            <Ionicons name="speedometer" size={16} color={COLORS.accent} />
            <Text style={styles.etaDetailText}>
              {Math.round(trackingData.speed)} km/h
            </Text>
          </View>
        </View>
      </View>

      {/* Driver Info Card */}
      <View style={styles.driverCard}>
        <View style={styles.driverCardHeader}>
          <Text style={styles.driverCardTitle}>Your Driver</Text>
          <View style={styles.driverVehicleBadge}>
            <Ionicons name={getVehicleIcon(trackingData.vehicleType)} size={14} color={COLORS.secondary} />
            <Text style={styles.driverVehicleText}>{trackingData.vehicleType}</Text>
          </View>
        </View>
        
        <View style={styles.driverCardContent}>
          <View style={styles.driverInfo}>
            <View style={styles.driverAvatar}>
              <Text style={styles.driverAvatarText}>
                {trackingData.driverName.charAt(0)}
              </Text>
            </View>
            <View style={styles.driverDetails}>
              <Text style={styles.driverName}>{trackingData.driverName}</Text>
              <View style={styles.driverRating}>
                <Ionicons name="star" size={14} color={COLORS.primary} />
                <Text style={styles.driverRatingText}>4.9</Text>
              </View>
            </View>
          </View>

          <View style={styles.driverActions}>
            <TouchableOpacity 
              style={[styles.actionButton, { backgroundColor: COLORS.info + '15' }]} 
              onPress={handleCallDriver}
            >
              <Ionicons name="call" size={22} color={COLORS.info} />
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.actionButton, { backgroundColor: COLORS.success + '15' }]} 
              onPress={handleWhatsApp}
            >
              <Ionicons name="logo-whatsapp" size={22} color={COLORS.success} />
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </View>
  );
}

function getVehicleIcon(vehicleType: string): any {
  const icons: Record<string, string> = {
    bicycle: 'bicycle',
    motorcycle: 'speedometer',
    car: 'car',
    van: 'bus',
    truck: 'cube',
  };
  return icons[vehicleType.toLowerCase()] || 'navigate';
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
  },
  map: {
    flex: 1,
  },
  
  // Loading State
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.secondary,
    padding: SPACING.xl,
  },
  loadingIcon: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.primary + '20',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.xl,
  },
  loadingTitle: {
    fontSize: TYPOGRAPHY.fontSize.xl,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.white,
    marginBottom: SPACING.sm,
  },
  loadingText: {
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.textMuted,
  },
  errorText: {
    fontSize: TYPOGRAPHY.fontSize.base,
    color: COLORS.accent,
    marginTop: SPACING.md,
  },
  successIcon: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.success,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.xl,
  },
  successTitle: {
    fontSize: TYPOGRAPHY.fontSize['2xl'],
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.white,
    marginBottom: SPACING.sm,
  },
  backButton: {
    marginTop: SPACING.xl,
    backgroundColor: COLORS.primary,
    flexDirection: 'row',
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.md,
    borderRadius: RADIUS.lg,
    alignItems: 'center',
    gap: SPACING.sm,
    ...SHADOWS.primary,
  },
  backButtonText: {
    color: COLORS.secondary,
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Connection Badge
  connectionBadge: {
    position: 'absolute',
    top: 60,
    left: SPACING.base,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.secondary,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: RADIUS.full,
    ...SHADOWS.lg,
  },
  connectionDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 6,
  },
  connectionText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.white,
  },
  
  // Driver Marker
  driverMarker: {
    backgroundColor: COLORS.primary,
    padding: 12,
    borderRadius: RADIUS.full,
    borderWidth: 3,
    borderColor: COLORS.white,
    ...SHADOWS.lg,
  },
  
  // ETA Card
  etaCard: {
    position: 'absolute',
    top: 60,
    right: SPACING.base,
    backgroundColor: COLORS.surface,
    borderRadius: RADIUS.xl,
    padding: SPACING.base,
    minWidth: 120,
    ...SHADOWS.lg,
  },
  etaMain: {
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  etaLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    letterSpacing: 0.5,
  },
  etaTime: {
    fontSize: TYPOGRAPHY.fontSize['4xl'],
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.primary,
    lineHeight: 44,
  },
  etaUnit: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textSecondary,
    marginTop: -4,
  },
  divider: {
    height: 1,
    backgroundColor: COLORS.border,
    marginVertical: SPACING.sm,
  },
  etaDetails: {
    gap: SPACING.sm,
  },
  etaDetailItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  etaDetailText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textSecondary,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  
  // Driver Card
  driverCard: {
    position: 'absolute',
    bottom: SPACING.xl,
    left: SPACING.base,
    right: SPACING.base,
    backgroundColor: COLORS.surface,
    borderRadius: RADIUS.xl,
    padding: SPACING.base,
    ...SHADOWS.xl,
  },
  driverCardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  driverCardTitle: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  driverVehicleBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.xs,
    borderRadius: RADIUS.full,
    gap: 4,
  },
  driverVehicleText: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.secondary,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    textTransform: 'capitalize',
  },
  driverCardContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  driverInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  driverAvatar: {
    width: 52,
    height: 52,
    borderRadius: 26,
    backgroundColor: COLORS.secondary,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: SPACING.md,
  },
  driverAvatarText: {
    color: COLORS.primary,
    fontSize: TYPOGRAPHY.fontSize.xl,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  driverDetails: {},
  driverName: {
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
  },
  driverRating: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    marginTop: 2,
  },
  driverRatingText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textSecondary,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  driverActions: {
    flexDirection: 'row',
    gap: SPACING.sm,
  },
  actionButton: {
    width: 48,
    height: 48,
    borderRadius: RADIUS.lg,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
