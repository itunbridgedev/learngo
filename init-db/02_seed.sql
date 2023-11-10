-- 02_seed.sql
INSERT INTO products (name, price) VALUES
('Product 1', 9.99),
('Product 2', 19.99),
('Product 3', 29.99);
INSERT INTO users (username, email, password_hash) VALUES 
('johndoe', 'john.doe@gmail.com', 'b7aa19d33add937b884a01b4567a7e2038d46f179ec3d317a66c2c5f927cd06c'),
('sallysue', 'sally.sue@gmail.com', 'b20200c33a3277c7473d6b7b039daf3ff1a181e59634dee943e3b566e1e867f0'),
('tomsawyer', 'tom.sawyer@gmail.com', '658454ee30a99986b819d8241ea77544ed44c17030868e8ffad19bafcfd44a81');
INSERT INTO orders (customer_id, total_price, status) VALUES 
(1, 49.99, 'pending'),
(1, 29.99, 'pending'),
(2, 19.99, 'pending'),
(3, 59.99, 'pending');