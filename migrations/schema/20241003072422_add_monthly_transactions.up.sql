CREATE TABLE IF NOT EXISTS `monthly_transactions` (
  id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id int NOT NULL,
  month_date DATE NOT NULL,
  total_expense decimal(15,2) NOT NULL,
  total_income decimal(15,2) NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_user_id (user_id),
  INDEX idx_month_date (month_date),
  UNIQUE INDEX idx_user_month_date (user_id, month_date)
);