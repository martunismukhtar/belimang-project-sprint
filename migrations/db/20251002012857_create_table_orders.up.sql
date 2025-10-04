CREATE TABLE orders (
    id UUID PRIMARY KEY,
    
    user_id UUID NULL,
    -- status order_status NOT NULL,
    total_price DECIMAL NOT NULL DEFAULT 0,

    created_at BIGINT DEFAULT extract(epoch from now()),
    

    
    -- CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE
);

