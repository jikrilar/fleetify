CREATE DATABASE IF NOT EXISTS fleetify_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE fleetify_db;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(100) NOT NULL,
  role ENUM('SA', 'APPROVAL') NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS vehicles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  license_plate VARCHAR(30) NOT NULL,
  model VARCHAR(100) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_vehicles_license_plate (license_plate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS master_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  item_name VARCHAR(150) NOT NULL,
  type ENUM('PART', 'SERVICE') NOT NULL,
  price DECIMAL(15,2) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS maintenance_reports (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  vehicle_id BIGINT UNSIGNED NOT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  odometer INT UNSIGNED NOT NULL,
  complaint TEXT NOT NULL,
  status ENUM('PENDING_APPROVAL', 'APPROVED', 'COMPLETED') NOT NULL,
  initial_photo VARCHAR(255) NULL,
  proof_photo VARCHAR(255) NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_reports_status (status),
  KEY idx_reports_created_by (created_by),
  KEY idx_reports_vehicle_id (vehicle_id),
  CONSTRAINT fk_reports_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
  CONSTRAINT fk_reports_user FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS report_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  report_id BIGINT UNSIGNED NOT NULL,
  item_id BIGINT UNSIGNED NOT NULL,
  quantity INT UNSIGNED NOT NULL,
  price_snapshot DECIMAL(15,2) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_report_items_report_id (report_id),
  KEY idx_report_items_item_id (item_id),
  CONSTRAINT fk_report_items_report FOREIGN KEY (report_id) REFERENCES maintenance_reports(id),
  CONSTRAINT fk_report_items_item FOREIGN KEY (item_id) REFERENCES master_items(id),
  CONSTRAINT chk_report_items_quantity CHECK (quantity >= 1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
