CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    number_of_items INT NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS shipping_packs (
    id SERIAL PRIMARY KEY,
    quantity INT NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS order_shipping (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    pack_size INT NOT NULL,
    shipping_pack_quantity INT NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

INSERT INTO shipping_packs (id, quantity) VALUES (DEFAULT, 5000);
INSERT INTO shipping_packs (id, quantity) VALUES (DEFAULT, 2000);
INSERT INTO shipping_packs (id, quantity) VALUES (DEFAULT, 1000);
INSERT INTO shipping_packs (id, quantity) VALUES (DEFAULT, 500);
INSERT INTO shipping_packs (id, quantity) VALUES (DEFAULT, 250);
