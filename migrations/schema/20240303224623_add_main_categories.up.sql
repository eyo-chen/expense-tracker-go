CREATE TABLE IF NOT EXISTS main_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type ENUM('1', '2') NOT NULL, -- 1 for 'income', 2 for 'expense'
    user_id INT NOT NULL,
    icon_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (icon_id) REFERENCES icons(id) ON DELETE CASCADE,
    INDEX idx_id (id),
    INDEX idx_user_id (user_id),
    INDEX idx_icon_id (icon_id)
);