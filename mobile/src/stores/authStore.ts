import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { authService } from '../services/auth';
import { Courier } from '../types';

interface AuthState {
    isAuthenticated: boolean;
    isLoading: boolean;
    token: string | null;
    courier: Courier | null;
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    checkAuth: () => Promise<void>;
    updateCourier: (courier: Partial<Courier>) => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
    isAuthenticated: false,
    isLoading: true,
    token: null,
    courier: null,

    login: async (email: string, password: string) => {
        const response = await authService.login(email, password);
        await SecureStore.setItemAsync('token', response.token);
        set({
            isAuthenticated: true,
            token: response.token,
            courier: response.courier,
        });
    },

    logout: async () => {
        await SecureStore.deleteItemAsync('token');
        set({
            isAuthenticated: false,
            token: null,
            courier: null,
        });
    },

    checkAuth: async () => {
        try {
            const token = await SecureStore.getItemAsync('token');
            if (token) {
                const courier = await authService.getProfile(token);
                set({
                    isAuthenticated: true,
                    token,
                    courier,
                    isLoading: false,
                });
            } else {
                set({ isLoading: false });
            }
        } catch (error) {
            await SecureStore.deleteItemAsync('token');
            set({ isLoading: false });
        }
    },

    updateCourier: (courierData: Partial<Courier>) => {
        const current = get().courier;
        if (current) {
            set({ courier: { ...current, ...courierData } });
        }
    },
}));
