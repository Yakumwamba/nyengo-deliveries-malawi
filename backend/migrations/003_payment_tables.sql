-- Nyengo Deliveries - Payment & Payout Tables Migration
-- Run this migration to create the required tables for the payment system

-- ============================================================
-- PAYOUTS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    courier_id UUID NOT NULL REFERENCES couriers(id) ON DELETE CASCADE,
    order_ids JSONB NOT NULL DEFAULT '[]',
    total_amount DECIMAL(12, 2) NOT NULL,
    platform_fee DECIMAL(12, 2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'ZMW',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payout_method VARCHAR(30) NOT NULL,
    payout_details TEXT,
    transaction_ref VARCHAR(100),
    failure_reason TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for payouts
CREATE INDEX IF NOT EXISTS idx_payouts_courier_id ON payouts(courier_id);
CREATE INDEX IF NOT EXISTS idx_payouts_status ON payouts(status);
CREATE INDEX IF NOT EXISTS idx_payouts_created_at ON payouts(created_at DESC);

-- ============================================================
-- PAYOUT_ORDERS TABLE (links orders to payouts)
-- ============================================================
CREATE TABLE IF NOT EXISTS payout_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payout_id UUID NOT NULL REFERENCES payouts(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(payout_id, order_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_payout_orders_payout_id ON payout_orders(payout_id);
CREATE INDEX IF NOT EXISTS idx_payout_orders_order_id ON payout_orders(order_id);

-- ============================================================
-- COURIER_WALLETS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS courier_wallets (
    courier_id UUID PRIMARY KEY REFERENCES couriers(id) ON DELETE CASCADE,
    available_balance DECIMAL(12, 2) NOT NULL DEFAULT 0,
    pending_balance DECIMAL(12, 2) NOT NULL DEFAULT 0,
    total_earnings DECIMAL(12, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'ZMW',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- WALLET_TRANSACTIONS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    courier_id UUID NOT NULL REFERENCES couriers(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    payout_id UUID REFERENCES payouts(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL, -- 'earning', 'payout', 'adjustment', 'refund'
    amount DECIMAL(12, 2) NOT NULL,
    balance_before DECIMAL(12, 2) NOT NULL,
    balance_after DECIMAL(12, 2) NOT NULL,
    description TEXT,
    reference VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_courier_id ON wallet_transactions(courier_id);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_created_at ON wallet_transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_type ON wallet_transactions(type);

-- ============================================================
-- STORE_PAYMENT_CONFIGS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS store_payment_configs (
    store_id UUID PRIMARY KEY,
    store_name VARCHAR(200) NOT NULL,
    payment_api_url TEXT NOT NULL,
    api_key TEXT NOT NULL,
    webhook_secret TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- ADD PAYMENT COLUMNS TO ORDERS TABLE (if not exists)
-- ============================================================
DO $$
BEGIN
    -- Add payment_reference if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'orders' AND column_name = 'payment_reference') THEN
        ALTER TABLE orders ADD COLUMN payment_reference VARCHAR(100);
    END IF;
    
    -- Add store_id if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'orders' AND column_name = 'store_id') THEN
        ALTER TABLE orders ADD COLUMN store_id UUID;
    END IF;
    
    -- Add external_order_id if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'orders' AND column_name = 'external_order_id') THEN
        ALTER TABLE orders ADD COLUMN external_order_id VARCHAR(200);
    END IF;
END
$$;

-- Index for store orders lookup
CREATE INDEX IF NOT EXISTS idx_orders_store_id ON orders(store_id);
CREATE INDEX IF NOT EXISTS idx_orders_external_order_id ON orders(external_order_id);

-- ============================================================
-- PAYOUT STATUS HISTORY (OPTIONAL - for audit trail)
-- ============================================================
CREATE TABLE IF NOT EXISTS payout_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payout_id UUID NOT NULL REFERENCES payouts(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    changed_by VARCHAR(100), -- 'system', 'admin:<id>', etc.
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payout_status_history_payout_id ON payout_status_history(payout_id);

-- ============================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================

-- Function to update timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for payouts
DROP TRIGGER IF EXISTS update_payouts_updated_at ON payouts;
CREATE TRIGGER update_payouts_updated_at
    BEFORE UPDATE ON payouts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for courier_wallets
DROP TRIGGER IF EXISTS update_courier_wallets_updated_at ON courier_wallets;
CREATE TRIGGER update_courier_wallets_updated_at
    BEFORE UPDATE ON courier_wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for store_payment_configs
DROP TRIGGER IF EXISTS update_store_payment_configs_updated_at ON store_payment_configs;
CREATE TRIGGER update_store_payment_configs_updated_at
    BEFORE UPDATE ON store_payment_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- FUNCTION: Record payout status changes (for audit)
-- ============================================================
CREATE OR REPLACE FUNCTION record_payout_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO payout_status_history (payout_id, old_status, new_status, changed_by)
        VALUES (NEW.id, OLD.status, NEW.status, 'system');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS payout_status_change_trigger ON payouts;
CREATE TRIGGER payout_status_change_trigger
    AFTER UPDATE ON payouts
    FOR EACH ROW
    EXECUTE FUNCTION record_payout_status_change();

-- ============================================================
-- SAMPLE DATA (for testing - remove in production)
-- ============================================================
-- Uncomment to insert sample store payment config
/*
INSERT INTO store_payment_configs (store_id, store_name, payment_api_url, api_key, webhook_secret)
VALUES (
    'a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6',
    'Sample Store',
    'https://store-api.example.com/api',
    'sk_test_sample_api_key_12345',
    'whsec_sample_webhook_secret_67890'
);
*/

-- ============================================================
-- GRANT PERMISSIONS (adjust as needed for your setup)
-- ============================================================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON payouts TO nyengo_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON payout_orders TO nyengo_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON courier_wallets TO nyengo_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON wallet_transactions TO nyengo_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON store_payment_configs TO nyengo_app;
-- GRANT SELECT, INSERT ON payout_status_history TO nyengo_app;

COMMENT ON TABLE payouts IS 'Courier payout requests for completed deliveries';
COMMENT ON TABLE payout_orders IS 'Links orders to payout requests';
COMMENT ON TABLE courier_wallets IS 'Courier in-app wallet balances';
COMMENT ON TABLE wallet_transactions IS 'Transaction history for courier wallets';
COMMENT ON TABLE store_payment_configs IS 'Payment API configuration for integrated stores';
COMMENT ON TABLE payout_status_history IS 'Audit trail for payout status changes';
