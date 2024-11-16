CREATE TABLE IF NOT EXISTS `monthly_transactions` (
  id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id int NOT NULL,
  month_date DATE NOT NULL,
  total_expense DECIMAL(15, 2) NOT NULL,
  total_income DECIMAL(15, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE INDEX unique_user_month_date (user_id, month_date)
);