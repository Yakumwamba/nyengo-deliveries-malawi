import axios from 'axios';
import { API_URL } from '../constants/config';
import { Courier, LoginResponse, RegisterRequest } from '../types';

const authApi = axios.create({
    baseURL: API_URL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

export const authService = {
    login: async (email: string, password: string): Promise<LoginResponse> => {
        const response = await authApi.post('/couriers/login', { email, password });
        if (!response.data.success) {
            throw new Error(response.data.error?.message || 'Login failed');
        }
        return response.data.data;
    },

    register: async (data: RegisterRequest): Promise<Courier> => {
        const response = await authApi.post('/couriers/register', data);
        if (!response.data.success) {
            throw new Error(response.data.error?.message || 'Registration failed');
        }
        return response.data.data;
    },

    getProfile: async (token: string): Promise<Courier> => {
        const response = await authApi.get('/couriers/profile', {
            headers: { Authorization: `Bearer ${token}` },
        });
        if (!response.data.success) {
            throw new Error(response.data.error?.message || 'Failed to get profile');
        }
        return response.data.data;
    },
};
