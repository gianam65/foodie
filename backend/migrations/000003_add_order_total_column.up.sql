-- Add total column to orders table (example migration)
-- Use IF NOT EXISTS to handle case where column already exists (from initial table creation)
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS total DECIMAL(10, 2) DEFAULT 0.00;

