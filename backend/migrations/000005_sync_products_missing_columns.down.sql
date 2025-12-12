-- Rollback: Remove synced columns from products table

ALTER TABLE products DROP COLUMN IF EXISTS product_image;
