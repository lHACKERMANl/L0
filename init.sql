CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50),
    entry VARCHAR(50),
    date_created TIMESTAMP,
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    oof_shard VARCHAR(10)
);

CREATE TABLE delivery (
    order_uid VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(200),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE payment (
    transaction VARCHAR(50) PRIMARY KEY,
    request_id VARCHAR(50),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount INT,
    payment_dt TIMESTAMP,
    bank VARCHAR(50),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE items (
    chrt_id INT PRIMARY KEY,
    order_uid VARCHAR(50),
    track_number VARCHAR(50),
    price INT,
    rid VARCHAR(50),
    name VARCHAR(100),
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id INT,
    brand VARCHAR(100),
    status INT
);

CREATE TABLE IF NOT EXISTS cache_table (
    key TEXT PRIMARY KEY,
    value TEXT
);

CREATE INDEX IF NOT EXISTS idx_cache_key ON cache_table (key);

INSERT INTO orders (order_uid, track_number, entry, date_created, customer_id, delivery_service, shardkey, sm_id, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', '2021-11-26T06:22:19Z', 'test', 'meest', '9', 99, '1');

INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ('b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com');

INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', '', 'USD', 'wbpay', 1817, '2021-11-26T06:22:19Z', 'alpha', 1500, 317, 0);

INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (9934930, 'b563feb7b2b84b6test', 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202);


