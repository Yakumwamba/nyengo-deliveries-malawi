import { useCallback, useEffect } from 'react';
import { useOrderStore } from '../stores/orderStore';

export function useOrders() {
    const {
        orders,
        currentOrder,
        isLoading,
        filters,
        pagination,
        fetchOrders,
        fetchOrderById,
        updateOrderStatus,
        acceptOrder,
        declineOrder,
        setFilters,
    } = useOrderStore();

    const refresh = useCallback(() => {
        fetchOrders(filters);
    }, [fetchOrders, filters]);

    const loadMore = useCallback(() => {
        if (pagination.page < pagination.totalPages) {
            fetchOrders({ ...filters, page: pagination.page + 1 });
        }
    }, [fetchOrders, filters, pagination]);

    useEffect(() => {
        fetchOrders();
    }, []);

    return {
        orders,
        currentOrder,
        isLoading,
        filters,
        pagination,
        refresh,
        loadMore,
        fetchOrderById,
        updateOrderStatus,
        acceptOrder,
        declineOrder,
        setFilters,
    };
}
