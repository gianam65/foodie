-- Add total column to orders table (example migration)
ALTER TABLE orders 
ADD COLUMN total DECIMAL(10, 2) DEFAULT 0.00;

