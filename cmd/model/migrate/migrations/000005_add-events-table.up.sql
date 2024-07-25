CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL, 
    description VARCHAR(1000) NOT NULL,
    start TIMESTAMP NOT NULL,
    "end" TIMESTAMP NOT NULL,
    reoccurrence VARCHAR(14),
    for_unloggedins BOOLEAN DEFAULT 'false' NOT NULL,
    for_unverifieds BOOLEAN DEFAULT 'false' NOT NULL,
    invite_only BOOLEAN DEFAULT 'false' NOT NULL,
    host_id INT NOT NULL,
    address_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_neighbors
        FOREIGN KEY(host_id)
            REFERENCES neighbors(id),
    CONSTRAINT fk_addresses
        FOREIGN KEY(address_id)
            REFERENCES addresses(id)
);