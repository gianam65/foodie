-- Remove columns added in this migration
ALTER TABLE orders 
    DROP COLUMN IF EXISTS restaurant_id,
    DROP COLUMN IF EXISTS items,
    DROP COLUMN IF EXISTS payment_method,
    DROP COLUMN IF EXISTS delivery_address,
    DROP COLUMN IF EXISTS updated_at;

