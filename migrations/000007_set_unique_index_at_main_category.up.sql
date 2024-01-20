ALTER TABLE main_categories
ADD CONSTRAINT unique_name_user_type UNIQUE (name, user_id, type);

ALTER TABLE main_categories
ADD CONSTRAINT unique_icon_user UNIQUE (icon_id, user_id);