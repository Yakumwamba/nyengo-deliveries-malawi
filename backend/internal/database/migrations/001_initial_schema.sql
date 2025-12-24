-- 001_initial_schema.sql
-- Nyengo Deliveries Initial Database Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Couriers table
CREATE TABLE IF NOT EXISTS couriers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    owner_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    alternate_phone VARCHAR(50),
    whatsapp VARCHAR(50),
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'Zambia',
    logo_url TEXT,
    description TEXT,
    service_areas TEXT[] DEFAULT '{}',
    vehicle_types TEXT[] DEFAULT '{}',
    max_weight DECIMAL(10, 2) DEFAULT 50.0,
    operating_hours JSONB,
    base_rate_per_km DECIMAL(10, 2) DEFAULT 5.0,
    minimum_fare DECIMAL(10, 2) DEFAULT 20.0,
    custom_pricing BOOLEAN DEFAULT FALSE,
    rating DECIMAL(3, 2) DEFAULT 0.0,
    total_reviews INTEGER DEFAULT 0,
    total_deliveries INTEGER DEFAULT 0,
    success_rate DECIMAL(5, 2) DEFAULT 0.0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    verification_docs TEXT[],
    wallet_balance DECIMAL(12, 2) DEFAULT 0.0,
    bank_details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_active_at TIMESTAMP WITH TIME ZONE
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    courier_id UUID NOT NULL REFERENCES couriers(id),
    store_id UUID,
    external_order_id VARCHAR(255),
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_email VARCHAR(255),
    pickup_address TEXT NOT NULL,
    pickup_latitude DECIMAL(10, 8) NOT NULL,
    pickup_longitude DECIMAL(11, 8) NOT NULL,
    pickup_notes TEXT,
    pickup_contact_name VARCHAR(255),
    pickup_contact_phone VARCHAR(50),
    delivery_address TEXT NOT NULL,
    delivery_latitude DECIMAL(10, 8) NOT NULL,
    delivery_longitude DECIMAL(11, 8) NOT NULL,
    delivery_notes TEXT,
    package_description TEXT NOT NULL,
    package_size VARCHAR(20) NOT NULL,
    package_weight DECIMAL(10, 2) DEFAULT 0.0,
    is_fragile BOOLEAN DEFAULT FALSE,
    requires_signature BOOLEAN DEFAULT FALSE,
    distance DECIMAL(10, 2) NOT NULL,
    base_fare DECIMAL(10, 2) NOT NULL,
    distance_fare DECIMAL(10, 2) NOT NULL,
    surge_fare DECIMAL(10, 2) DEFAULT 0.0,
    total_fare DECIMAL(10, 2) NOT NULL,
    platform_fee DECIMAL(10, 2) NOT NULL,
    courier_earnings DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_status VARCHAR(50) DEFAULT 'pending',
    payment_reference VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    status_history JSONB DEFAULT '[]',
    scheduled_pickup TIMESTAMP WITH TIME ZONE,
    actual_pickup TIMESTAMP WITH TIME ZONE,
    estimated_delivery TIMESTAMP WITH TIME ZONE,
    actual_delivery TIMESTAMP WITH TIME ZONE,
    delivery_proof_url TEXT,
    recipient_name VARCHAR(255),
    signature_url TEXT,
    customer_rating INTEGER,
    customer_feedback TEXT,
    notes TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Delivery tracking table
CREATE TABLE IF NOT EXISTS delivery_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    courier_id UUID NOT NULL REFERENCES couriers(id),
    driver_name VARCHAR(255) NOT NULL,
    driver_phone VARCHAR(50) NOT NULL,
    vehicle_type VARCHAR(50),
    vehicle_plate VARCHAR(50),
    current_latitude DECIMAL(10, 8),
    current_longitude DECIMAL(11, 8),
    last_location_at TIMESTAMP WITH TIME ZONE,
    route_polyline TEXT,
    estimated_arrival TIMESTAMP WITH TIME ZONE,
    distance_remaining DECIMAL(10, 2),
    duration_remaining INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Location history table
CREATE TABLE IF NOT EXISTS location_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tracking_id UUID NOT NULL REFERENCES delivery_tracking(id),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    speed DECIMAL(6, 2),
    heading DECIMAL(6, 2),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Chat conversations table
CREATE TABLE IF NOT EXISTS chat_conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    courier_id UUID NOT NULL REFERENCES couriers(id),
    customer_phone VARCHAR(50) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_message_at TIMESTAMP WITH TIME ZONE,
    unread_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES chat_conversations(id),
    sender_type VARCHAR(20) NOT NULL,
    sender_id VARCHAR(255) NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    content TEXT,
    media_url TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API Keys table (for store integration)
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_name VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    permissions TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_orders_courier_id ON orders(courier_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_delivery_tracking_order_id ON delivery_tracking(order_id);
CREATE INDEX idx_delivery_tracking_is_active ON delivery_tracking(is_active);
CREATE INDEX idx_location_history_tracking_id ON location_history(tracking_id);
CREATE INDEX idx_chat_conversations_order_id ON chat_conversations(order_id);
CREATE INDEX idx_chat_messages_conversation_id ON chat_messages(conversation_id);
CREATE INDEX idx_couriers_email ON couriers(email);
CREATE INDEX idx_couriers_is_active ON couriers(is_active);
