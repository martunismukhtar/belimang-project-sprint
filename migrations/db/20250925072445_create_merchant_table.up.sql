CREATE TABLE merchants (
    id uuid PRIMARY KEY,
    name VARCHAR(30) NOT NULL CHECK (char_length(name) >= 2),
    merchant_category VARCHAR(30) NOT NULL CHECK (
        merchant_category IN (
            'SmallRestaurant',
            'MediumRestaurant',
            'LargeRestaurant',
            'MerchandiseRestaurant',
            'BoothKiosk',
            'ConvenienceStore'
        )
    ),
    image_url TEXT NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM clock_timestamp()) * 1e9)::BIGINT
);

CREATE INDEX idx_merchants_name ON merchants(name);
CREATE INDEX idx_merchants_category ON merchants(merchant_category);
CREATE INDEX idx_merchants_location ON merchants(lat, long);