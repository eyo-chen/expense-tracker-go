ALTER TABLE main_categories
ADD CONSTRAINT unique_icon_user UNIQUE (icon_id, user_id);