import { create } from 'zustand';
import { Order, OrderFilters } from '../types';
import { orderService } from '../services/api';

interface OrderState {
    orders: Order[];
    currentOrder: Order | null;
    isLoading: boolean;
    filters: OrderFilters;
    pagination: {
        page: number;
        pageSize: number;
        totalCount: number;
        totalPages: number;
    };
    fetchOrders: (filters?: OrderFilters) => Promise<void>;
    fetchOrderById: (id: string) => Promise<void>;
    updateOrderStatus: (id: string, status: string) => Promise<void>;
    acceptOrder: (id: string) => Promise<void>;
    declineOrder: (id: string) => Promise<void>;
    setFilters: (filters: Partial<OrderFilters>) => void;
}

export const useOrderStore = create<OrderState>((set, get) => ({
    orders: [],
    currentOrder: null,
    isLoading: false,
    filters: { page: 1, pageSize: 20 },
    pagination: { page: 1, pageSize: 20, totalCount: 0, totalPages: 0 },

    fetchOrders: async (filters?: OrderFilters) => {
        set({ isLoading: true });
        try {
            const mergedFilters = { ...get().filters, ...filters };
            const response = await orderService.list(mergedFilters);
            set({
                orders: response.orders,
                pagination: {
                    page: response.page,
                    pageSize: response.pageSize,
                    totalCount: response.totalCount,
                    totalPages: response.totalPages,
                },
                filters: mergedFilters,
            });
        } finally {
            set({ isLoading: false });
        }
    },

    fetchOrderById: async (id: string) => {
        set({ isLoading: true });
        try {
            const order = await orderService.getById(id);
            set({ currentOrder: order });
        } finally {
            set({ isLoading: false });
        }
    },

    updateOrderStatus: async (id: string, status: string) => {
        await orderService.updateStatus(id, status);
        const orders = get().orders.map((o) =>
            o.id === id ? { ...o, status } : o
        );
        set({ orders });
    },

    acceptOrder: async (id: string) => {
        await orderService.accept(id);
        await get().fetchOrders();
    },

    declineOrder: async (id: string) => {
        await orderService.decline(id);
        await get().fetchOrders();
    },

    setFilters: (filters: Partial<OrderFilters>) => {
        set({ filters: { ...get().filters, ...filters } });
    },
}));
