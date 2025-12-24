import { useState, useEffect, useRef, useCallback } from 'react';
import * as Location from 'expo-location';
import { API_URL, WS_URL } from '../constants/config';
import { useAuthStore } from '../stores/authStore';

interface TrackingData {
    orderId: string;
    latitude: number;
    longitude: number;
    speed: number;
    heading: number;
    distanceRemaining: number;
    etaMinutes: number;
    estimatedArrival: string;
    driverName: string;
    driverPhone: string;
    vehicleType: string;
    isActive: boolean;
}

interface LocationUpdateResult {
    success: boolean;
    distanceRemaining?: number;
    etaMinutes?: number;
}

// Hook for DRIVER to send location updates
export function useDriverTracking(orderId: string) {
    const [isTracking, setIsTracking] = useState(false);
    const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
    const [error, setError] = useState<string | null>(null);
    const { token } = useAuthStore();
    const locationSubscription = useRef<Location.LocationSubscription | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    // Start tracking and sending locations
    const startTracking = useCallback(async (driverInfo: {
        name: string;
        phone: string;
        vehicleType: string;
        vehiclePlate?: string;
    }) => {
        try {
            // Request location permissions
            const { status } = await Location.requestForegroundPermissionsAsync();
            if (status !== 'granted') {
                throw new Error('Location permission denied');
            }

            // Start tracking on server
            const response = await fetch(`${API_URL}/tracking/${orderId}/start`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(driverInfo),
            });

            if (!response.ok) {
                throw new Error('Failed to start tracking');
            }

            // Connect WebSocket for real-time updates
            wsRef.current = new WebSocket(`${WS_URL}?token=${token}`);

            wsRef.current.onopen = () => {
                console.log('Tracking WebSocket connected');
            };

            // Start watching location
            locationSubscription.current = await Location.watchPositionAsync(
                {
                    accuracy: Location.Accuracy.High,
                    timeInterval: 5000,        // Update every 5 seconds
                    distanceInterval: 10,      // Or every 10 meters
                },
                async (location) => {
                    const { latitude, longitude, speed, heading, altitude, accuracy } = location.coords;

                    // Send via WebSocket for faster delivery
                    if (wsRef.current?.readyState === WebSocket.OPEN) {
                        wsRef.current.send(JSON.stringify({
                            type: 'location_update',
                            payload: {
                                orderId,
                                latitude,
                                longitude,
                                speed: speed ? speed * 3.6 : 0, // m/s to km/h
                                heading: heading || 0,
                                altitude: altitude || 0,
                                accuracy: accuracy || 0,
                            },
                        }));
                    }

                    // Also send via HTTP as backup
                    try {
                        await fetch(`${API_URL}/tracking/${orderId}/location`, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                                Authorization: `Bearer ${token}`,
                            },
                            body: JSON.stringify({
                                latitude,
                                longitude,
                                speed: speed ? speed * 3.6 : 0,
                                heading: heading || 0,
                                altitude: altitude || 0,
                                accuracy: accuracy || 0,
                            }),
                        });
                    } catch (e) {
                        console.warn('HTTP location update failed, WebSocket should handle it');
                    }

                    setLastUpdate(new Date());
                }
            );

            setIsTracking(true);
            setError(null);
        } catch (e: any) {
            setError(e.message);
            throw e;
        }
    }, [orderId, token]);

    // Stop tracking
    const stopTracking = useCallback(async (reason: string = 'completed') => {
        if (locationSubscription.current) {
            locationSubscription.current.remove();
            locationSubscription.current = null;
        }

        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }

        try {
            await fetch(`${API_URL}/tracking/${orderId}/stop`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({ reason }),
            });
        } catch (e) {
            console.error('Failed to stop tracking on server:', e);
        }

        setIsTracking(false);
    }, [orderId, token]);

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            if (locationSubscription.current) {
                locationSubscription.current.remove();
            }
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, []);

    return {
        isTracking,
        lastUpdate,
        error,
        startTracking,
        stopTracking,
    };
}

// Hook for CUSTOMER/STORE to watch tracking updates
export function useTrackingWatch(orderId: string) {
    const [trackingData, setTrackingData] = useState<TrackingData | null>(null);
    const [locationHistory, setLocationHistory] = useState<Array<{ latitude: number; longitude: number; timestamp: string }>>([]);
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const { token } = useAuthStore();

    // Fetch initial tracking data
    const fetchTrackingData = useCallback(async () => {
        try {
            const response = await fetch(`${API_URL}/tracking/${orderId}`);
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    setTrackingData({
                        orderId,
                        latitude: data.data.currentLocation?.latitude || 0,
                        longitude: data.data.currentLocation?.longitude || 0,
                        speed: data.data.currentLocation?.speed || 0,
                        heading: data.data.currentLocation?.heading || 0,
                        distanceRemaining: data.data.distanceRemaining || 0,
                        etaMinutes: data.data.etaMinutes || 0,
                        estimatedArrival: data.data.estimatedArrival || '',
                        driverName: data.data.driverName || '',
                        driverPhone: data.data.driverPhone || '',
                        vehicleType: data.data.vehicleType || '',
                        isActive: data.data.isActive || false,
                    });
                }
            }
        } catch (e: any) {
            setError(e.message);
        }
    }, [orderId]);

    // Connect to WebSocket for live updates
    useEffect(() => {
        if (!orderId) return;

        // Fetch initial data
        fetchTrackingData();

        // Connect WebSocket
        const wsUrl = token ? `${WS_URL}?token=${token}` : WS_URL;
        wsRef.current = new WebSocket(wsUrl);

        wsRef.current.onopen = () => {
            setIsConnected(true);

            // Subscribe to this order's updates
            wsRef.current?.send(JSON.stringify({
                type: 'subscribe',
                payload: {
                    orderId,
                    action: 'subscribe',
                },
            }));
        };

        wsRef.current.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);

                if (message.type === 'location_update') {
                    const update = message.payload;

                    if (update.orderId === orderId) {
                        setTrackingData((prev) => ({
                            ...prev!,
                            latitude: update.latitude,
                            longitude: update.longitude,
                            speed: update.speed || 0,
                            heading: update.heading || 0,
                            distanceRemaining: update.distanceRemaining || prev?.distanceRemaining || 0,
                            etaMinutes: update.etaMinutes || prev?.etaMinutes || 0,
                        }));

                        // Add to history
                        setLocationHistory((prev) => [
                            ...prev.slice(-99), // Keep last 100 points
                            {
                                latitude: update.latitude,
                                longitude: update.longitude,
                                timestamp: new Date().toISOString(),
                            },
                        ]);
                    }
                }

                if (message.type === 'tracking_stopped') {
                    setTrackingData((prev) => prev ? { ...prev, isActive: false } : null);
                }
            } catch (e) {
                console.error('Failed to parse WebSocket message:', e);
            }
        };

        wsRef.current.onclose = () => {
            setIsConnected(false);
        };

        wsRef.current.onerror = () => {
            setError('WebSocket connection error');
            setIsConnected(false);
        };

        // Ping every 30 seconds to keep connection alive
        const pingInterval = setInterval(() => {
            if (wsRef.current?.readyState === WebSocket.OPEN) {
                wsRef.current.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000);

        return () => {
            clearInterval(pingInterval);
            if (wsRef.current) {
                wsRef.current.send(JSON.stringify({
                    type: 'subscribe',
                    payload: { orderId, action: 'unsubscribe' },
                }));
                wsRef.current.close();
            }
        };
    }, [orderId, token, fetchTrackingData]);

    return {
        trackingData,
        locationHistory,
        isConnected,
        error,
        refetch: fetchTrackingData,
    };
}

// Fetch location history for an order
export async function fetchLocationHistory(orderId: string, limit: number = 100) {
    const response = await fetch(`${API_URL}/tracking/${orderId}/history?limit=${limit}`);
    if (!response.ok) {
        throw new Error('Failed to fetch location history');
    }
    const data = await response.json();
    return data.data?.points || [];
}
