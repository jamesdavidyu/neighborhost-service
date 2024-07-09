CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY, 
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    address VARCHAR(50) NOT NULL,
    city VARCHAR(45) NOT NULL,
    state VARCHAR(35) NOT NULL,
    zipcode VARCHAR(5) NOT NULL,
    neighbor_id INT NOT NULL UNIQUE,
    neighborhood_id INT NOT NULL, 
    recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_neighbors
        FOREIGN KEY(neighbor_id)
            REFERENCES neighbors(id),
    CONSTRAINT fk_neighborhoods
        FOREIGN KEY(neighborhood_id)
            REFERENCES neighborhoods(id)
);