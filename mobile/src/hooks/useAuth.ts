import { useCallback } from 'react';
import { useAuthStore } from '../stores/authStore';
import { authService } from '../services/auth';

export function useAuth() {
    const { isAuthenticated, isLoading, courier, token, login, logout, updateCourier } = useAuthStore();

    const refreshProfile = useCallback(async () => {
        if (!token) return;
        try {
            const profile = await authService.getProfile(token);
            updateCourier(profile);
        } catch (error) {
            console.error('Failed to refresh profile:', error);
        }
    }, [token, updateCourier]);

    return {
        isAuthenticated,
        isLoading,
        courier,
        token,
        login,
        logout,
        refreshProfile,
    };
}
