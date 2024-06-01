ALTER TABLE users
ADD COLUMN is_set_init_category BOOLEAN NOT NULL DEFAULT FALSE AFTER password_hash;