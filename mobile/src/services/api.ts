import axios from 'axios';
import * as SecureStore from 'expo-secure-store';
import { API_URL } from '../constants/config';

const api = axios.create({
    baseURL: API_URL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

// Add auth token to requests
api.interceptors.request.use(async (config) => {
    const token = await SecureStore.getItemAsync('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Handle response errors
api.interceptors.response.use(
    (response) => response.data,
    (error) => {
        const message = error.response?.data?.error?.message || error.message || 'An error occurred';
        return Promise.reject(new Error(message));
    }
);

export const orderService = {
    list: (filters: any) => api.get('/orders', { params: filters }),
    getById: (id: string) => api.get(`/orders/${id}`),
    create: (data: any) => api.post('/orders', data),
    updateStatus: (id: string, status: string) => api.put(`/orders/${id}/status`, { status }),
    accept: (id: string) => api.put(`/orders/${id}/accept`),
    decline: (id: string) => api.put(`/orders/${id}/decline`),
};

export const pricingService = {
    getEstimate: (data: any) => api.post('/pricing/estimate', data),
};

export const courierService = {
    getProfile: () => api.get('/couriers/profile'),
    updateProfile: (data: any) => api.put('/couriers/profile', data),
    getDashboard: () => api.get('/couriers/dashboard'),
};

export default api;
