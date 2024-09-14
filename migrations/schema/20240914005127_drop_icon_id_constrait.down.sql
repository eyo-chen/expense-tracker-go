ALTER TABLE main_categories
ADD CONSTRAINT main_categories_ibfk_2 FOREIGN KEY (icon_id) REFERENCES icons(id);
