ALTER TABLE main_categories
ADD CONSTRAINT unique_name_user_type UNIQUE (name, user_id, type);