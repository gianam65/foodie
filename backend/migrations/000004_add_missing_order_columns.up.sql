-- Add missing columns to orders table (for existing databases)
-- This migration is safe to run even if columns already exist

-- Add restaurant_id if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='orders' AND column_name='restaurant_id') THEN
        ALTER TABLE orders ADD COLUMN restaurant_id VARCHAR(36) NOT NULL DEFAULT '';
    END IF;
END $$;

-- Add items if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='orders' AND column_name='items') THEN
        ALTER TABLE orders ADD COLUMN items JSONB NOT NULL DEFAULT '[]'::jsonb;
    END IF;
END $$;

-- Add payment_method if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='orders' AND column_name='payment_method') THEN
        ALTER TABLE orders ADD COLUMN payment_method VARCHAR(50) NOT NULL DEFAULT '';
    END IF;
END $$;

-- Add delivery_address if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='orders' AND column_name='delivery_address') THEN
        ALTER TABLE orders ADD COLUMN delivery_address TEXT NOT NULL DEFAULT '';
    END IF;
END $$;

-- Add updated_at if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='orders' AND column_name='updated_at') THEN
        ALTER TABLE orders ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();
    END IF;
END $$;

