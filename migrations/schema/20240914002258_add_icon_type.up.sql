ALTER TABLE main_categories
ADD COLUMN icon_type ENUM('1', '2') NOT NULL AFTER user_id; -- 1 for 'default', 2 for 'custom'