CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    neighbor_id INT NOT NULL,
    neighbors_friend_id INT NOT NULL,
    friended_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_neighbors
        FOREIGN KEY (neighbor_id)
            REFERENCES neighbors(id),
    CONSTRAINT fk_friends
        FOREIGN KEY (neighbors_friend_id)
            REFERENCES neighbors(id)
);