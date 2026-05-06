-- +goose Up
ALTER TABLE import_items ADD COLUMN product_id UUID;
ALTER TABLE import_items ADD COLUMN processed_at TIMESTAMP;
ALTER TABLE import_items ADD COLUMN metadata JSONB;
ALTER TABLE import_items ADD COLUMN display_order INT DEFAULT 0;
ALTER TABLE import_items ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT now();

-- Create index on product_id for faster lookups
CREATE INDEX import_items_product_id_idx ON import_items(product_id) WHERE product_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS import_items_product_id_idx;
ALTER TABLE import_items DROP COLUMN IF EXISTS product_id;
ALTER TABLE import_items DROP COLUMN IF EXISTS processed_at;
ALTER TABLE import_items DROP COLUMN IF EXISTS metadata;
ALTER TABLE import_items DROP COLUMN IF EXISTS display_order;
ALTER TABLE import_items DROP COLUMN IF EXISTS created_at;