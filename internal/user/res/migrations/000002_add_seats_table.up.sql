CREATE TABLE IF NOT EXISTS seats (
    tenant_id VARCHAR(36) PRIMARY KEY,
    seats_used INTEGER DEFAULT 0,
    max_seats INTEGER DEFAULT 0
);