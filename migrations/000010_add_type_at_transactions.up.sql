ALTER TABLE transactions
ADD COLUMN type ENUM('1', '2') NOT NULL AFTER user_id; -- 1 for 'income', 2 for 'expense'


