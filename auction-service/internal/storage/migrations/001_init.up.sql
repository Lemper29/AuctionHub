CREATE TABLE lots (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_price DECIMAL(15, 2) NOT NULL,
    current_price DECIMAL(15, 2) NOT NULL,
    current_winner VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    end_time_unix BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids (
    id VARCHAR(255) PRIMARY KEY,
    lot_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    timestamp_unix BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE CASCADE
);

CREATE INDEX idx_lots_status ON lots(status);
CREATE INDEX idx_lots_end_time ON lots(end_time_unix);
CREATE INDEX idx_bids_lot_id ON bids(lot_id);
CREATE INDEX idx_bids_user_id ON bids(user_id);
CREATE INDEX idx_bids_timestamp ON bids(timestamp_unix);