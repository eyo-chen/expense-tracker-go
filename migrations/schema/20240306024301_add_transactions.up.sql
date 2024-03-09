CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    main_category_id INT NOT NULL,
    sub_category_id INT NOT NULL,
    price DECIMAL(12, 2) NOT NULL,
    note VARCHAR(255),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (main_category_id) REFERENCES main_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (sub_category_id) REFERENCES sub_categories(id) ON DELETE CASCADE,
    INDEX idx_id (id),
    INDEX idx_user_id (user_id),
    INDEX idx_main_category_id (main_category_id),
    INDEX idx_sub_category_id (sub_category_id)
);