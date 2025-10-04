CREATE TABLE delivery_estimate (
    id UUID PRIMARY KEY,
    user_id UUID NULL,
    orders JSONB NOT NULL, 
    total_price DECIMAL(10, 2) NOT NULL, -- harga saat estimasi (preview)
    estimated_delivery_time_minutes FLOAT NOT NULL,    
    created_at TIMESTAMP DEFAULT NOW()

    -- constraint fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
