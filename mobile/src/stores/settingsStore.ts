import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';

interface SettingsState {
    currency: string;
    currencySymbol: string;
    notifications: boolean;
    darkMode: boolean;
    language: string;
    setCurrency: (currency: string) => void;
    setNotifications: (enabled: boolean) => void;
    setDarkMode: (enabled: boolean) => void;
    setLanguage: (language: string) => void;
    loadSettings: () => Promise<void>;
}

export const useSettingsStore = create<SettingsState>((set) => ({
    currency: 'ZMW',
    currencySymbol: 'K',
    notifications: true,
    darkMode: false,
    language: 'en',

    setCurrency: (currency: string) => {
        const symbols: Record<string, string> = { ZMW: 'K', USD: '$', ZAR: 'R', KES: 'KSh' };
        set({ currency, currencySymbol: symbols[currency] || currency });
        SecureStore.setItemAsync('currency', currency);
    },

    setNotifications: (enabled: boolean) => {
        set({ notifications: enabled });
        SecureStore.setItemAsync('notifications', String(enabled));
    },

    setDarkMode: (enabled: boolean) => {
        set({ darkMode: enabled });
        SecureStore.setItemAsync('darkMode', String(enabled));
    },

    setLanguage: (language: string) => {
        set({ language });
        SecureStore.setItemAsync('language', language);
    },

    loadSettings: async () => {
        const currency = await SecureStore.getItemAsync('currency');
        const notifications = await SecureStore.getItemAsync('notifications');
        const darkMode = await SecureStore.getItemAsync('darkMode');
        const language = await SecureStore.getItemAsync('language');

        set({
            currency: currency || 'ZMW',
            currencySymbol: currency === 'USD' ? '$' : currency === 'ZAR' ? 'R' : 'K',
            notifications: notifications !== 'false',
            darkMode: darkMode === 'true',
            language: language || 'en',
        });
    },
}));
