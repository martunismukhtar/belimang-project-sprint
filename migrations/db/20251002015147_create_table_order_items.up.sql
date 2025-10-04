CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    merchant_id UUID NOT NULL,    
    item_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now()),

    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_item_id FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
    CONSTRAINT fk_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants (id) ON DELETE CASCADE
);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_item_id ON order_items (item_id);

CREATE INDEX idx_orders_merchant_id ON orders (merchant_id);
-- CREATE INDEX idx_merchants_name ON merchants(name);