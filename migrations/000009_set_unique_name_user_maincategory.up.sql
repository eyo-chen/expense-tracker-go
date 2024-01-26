ALTER TABLE sub_categories
ADD CONSTRAINT unique_name_user_maincategory UNIQUE (name, user_id, main_category_id);