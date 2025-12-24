export interface Courier {
    id: string;
    email: string;
    companyName: string;
    ownerName: string;
    phone: string;
    alternatePhone?: string;
    whatsapp?: string;
    address: string;
    city: string;
    country: string;
    logoUrl?: string;
    description?: string;
    serviceAreas: string[];
    vehicleTypes: string[];
    maxWeight: number;
    baseRatePerKm: number;
    minimumFare: number;
    rating: number;
    totalReviews: number;
    totalDeliveries: number;
    successRate: number;
    isVerified: boolean;
    isActive: boolean;
    walletBalance: number;
    createdAt: string;
    updatedAt: string;
}

export interface LoginResponse {
    token: string;
    courier: Courier;
}

export interface RegisterRequest {
    email: string;
    password: string;
    companyName: string;
    ownerName: string;
    phone: string;
    address: string;
    city: string;
    country: string;
    serviceAreas: string[];
    vehicleTypes: string[];
}

export interface Order {
    id: string;
    orderNumber: string;
    courierId: string;
    storeId?: string;
    customerName: string;
    customerPhone: string;
    customerEmail?: string;
    pickupAddress: string;
    pickupLatitude: number;
    pickupLongitude: number;
    pickupNotes?: string;
    deliveryAddress: string;
    deliveryLatitude: number;
    deliveryLongitude: number;
    deliveryNotes?: string;
    packageDescription: string;
    packageSize: string;
    packageWeight: number;
    isFragile: boolean;
    distance: number;
    baseFare: number;
    distanceFare: number;
    surgeFare: number;
    totalFare: number;
    platformFee: number;
    courierEarnings: number;
    paymentMethod: string;
    paymentStatus: string;
    status: string;
    createdAt: string;
    updatedAt: string;
}

export interface OrderFilters {
    status?: string[];
    dateFrom?: string;
    dateTo?: string;
    search?: string;
    page?: number;
    pageSize?: number;
}

export interface PriceEstimate {
    currency: string;
    currencySymbol: string;
    distance: number;
    duration: number;
    baseFare: number;
    distanceFare: number;
    totalFare: number;
    formattedTotal: string;
}

export interface DashboardStats {
    today: {
        totalOrders: number;
        completed: number;
        pending: number;
        totalRevenue: number;
        totalEarnings: number;
    };
    thisMonth: {
        totalOrders: number;
        completed: number;
        totalRevenue: number;
        totalEarnings: number;
        avgRating: number;
    };
}
