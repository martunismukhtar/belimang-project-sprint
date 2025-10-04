CREATE TABLE items (
    id UUID PRIMARY KEY,
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,    
    product_category VARCHAR(30) NOT NULL CHECK (
        product_category IN (
            'Beverage',
            'Food',
            'Snack',
            'Condiments',
            'Additions'
        )
    ),
    price NUMERIC(15,2) NOT NULL,
    image_url TEXT,
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM clock_timestamp()) * 1e9)::BIGINT
);

CREATE INDEX idx_items_merchant_id ON items (merchant_id);