USE fleetify_db;

INSERT INTO users (id, username, role, created_at, updated_at) VALUES
  (1, 'sa_user', 'SA', NOW(), NOW()),
  (2, 'approval_user', 'APPROVAL', NOW(), NOW())
ON DUPLICATE KEY UPDATE username = VALUES(username), role = VALUES(role), updated_at = NOW();

INSERT INTO vehicles (id, license_plate, model, created_at, updated_at) VALUES
  (1, 'B 1234 FTY', 'Toyota Avanza', NOW(), NOW()),
  (2, 'B 5678 FTY', 'Daihatsu Gran Max', NOW(), NOW()),
  (3, 'B 9012 FTY', 'Mitsubishi L300', NOW(), NOW())
ON DUPLICATE KEY UPDATE model = VALUES(model), updated_at = NOW();

INSERT INTO master_items (id, item_name, type, price, created_at, updated_at) VALUES
  (1, 'Engine Oil', 'PART', 350000, NOW(), NOW()),
  (2, 'Oil Filter', 'PART', 85000, NOW(), NOW()),
  (3, 'Brake Pad', 'PART', 450000, NOW(), NOW()),
  (4, 'General Service', 'SERVICE', 250000, NOW(), NOW()),
  (5, 'Brake Inspection', 'SERVICE', 150000, NOW(), NOW())
ON DUPLICATE KEY UPDATE item_name = VALUES(item_name), type = VALUES(type), price = VALUES(price), updated_at = NOW();
