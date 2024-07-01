CREATE TABLE IF NOT EXISTS neighbors (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(30) NOT NULL UNIQUE,
    zipcode VARCHAR(5) NOT NULL,
    password VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT 'false' NOT NULL,
    neighborhood_id INT DEFAULT 1 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_neighborhoods
        FOREIGN KEY(neighborhood_id)
            REFERENCES neighborhoods(id)
);