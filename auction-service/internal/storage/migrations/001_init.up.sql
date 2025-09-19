-- Таблица lots
CREATE TABLE IF NOT EXISTS lots (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    start_price DOUBLE PRECISION,
    current_price DOUBLE PRECISION,
    current_winner VARCHAR(255),
    status VARCHAR(50),
    end_time_unix BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица bids
CREATE TABLE IF NOT EXISTS bids (
    id VARCHAR(255) PRIMARY KEY,
    lot_id VARCHAR(255),
    user_id VARCHAR(255),
    amount DOUBLE PRECISION,
    timestamp_unix BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_lots_status ON lots(status);
CREATE INDEX idx_lots_end_time ON lots(end_time_unix);
CREATE INDEX idx_bids_lot_id ON bids(lot_id);
CREATE INDEX idx_bids_user_id ON bids(user_id);
CREATE INDEX idx_bids_timestamp ON bids(timestamp_unix);