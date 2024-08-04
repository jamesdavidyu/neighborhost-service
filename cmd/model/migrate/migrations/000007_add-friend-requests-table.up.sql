CREATE TABLE IF NOT EXISTS friend_requests (
    id SERIAL PRIMARY KEY,
    neighbor_id INT NOT NULL,
    requesting_friend_id INT NOT NULL,
    status VARCHAR(10) NOT NULL,
    friend_requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_neighbors
        FOREIGN KEY(neighbor_id)
            REFERENCES neighbors(id),
    CONSTRAINT fk_friends
        FOREIGN KEY(requesting_friend_id)
            REFERENCES neighbors(id)
);