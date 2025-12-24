import { useState, useCallback } from 'react';
import { pricingService } from '../services/api';
import { PriceEstimate } from '../types';

interface PricingRequest {
    pickupLatitude: number;
    pickupLongitude: number;
    deliveryLatitude: number;
    deliveryLongitude: number;
    packageSize?: string;
    packageWeight?: number;
    isFragile?: boolean;
    isExpress?: boolean;
}

export function usePricing() {
    const [estimate, setEstimate] = useState<PriceEstimate | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const getEstimate = useCallback(async (request: PricingRequest) => {
        setIsLoading(true);
        setError(null);
        try {
            const result = await pricingService.getEstimate(request);
            setEstimate(result.data);
            return result.data;
        } catch (err: any) {
            setError(err.message || 'Failed to get price estimate');
            throw err;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const clearEstimate = useCallback(() => {
        setEstimate(null);
        setError(null);
    }, []);

    return {
        estimate,
        isLoading,
        error,
        getEstimate,
        clearEstimate,
    };
}
