CREATE TABLE IF NOT EXISTS accounts (
  id VARCHAR(20) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100),
  phone VARCHAR(30),
  industry VARCHAR(50),
  sap_customer_id VARCHAR(20),
  status VARCHAR(20) DEFAULT 'ACTIVE',
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS opportunities (
  id VARCHAR(20) PRIMARY KEY,
  account_id VARCHAR(20) REFERENCES accounts(id),
  name VARCHAR(100),
  value NUMERIC(15,2),
  currency VARCHAR(5) DEFAULT 'USD',
  stage VARCHAR(30) DEFAULT 'PROSPECTING',
  sap_order_id VARCHAR(20),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS leads (
  id VARCHAR(20) PRIMARY KEY,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  email VARCHAR(100),
  company VARCHAR(100),
  status VARCHAR(20) DEFAULT 'NEW',
  source VARCHAR(50),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contacts (
  id VARCHAR(20) PRIMARY KEY,
  account_id VARCHAR(20) REFERENCES accounts(id),
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  email VARCHAR(100),
  phone VARCHAR(30),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Seed accounts
INSERT INTO accounts (id, name, email, phone, industry, sap_customer_id)
VALUES
  ('ACC-001', 'PT Maju Bersama', 'crm@majubersama.co.id', '+6221-555-1001', 'Manufacturing', 'C-001'),
  ('ACC-002', 'CV Teknologi Nusantara', 'crm@teknusa.co.id', '+6221-555-1002', 'Technology', 'C-002'),
  ('ACC-003', 'PT Sumber Rezeki', 'crm@sumberrezeki.com', '+6221-555-1003', 'Trading', 'C-003')
ON CONFLICT DO NOTHING;

-- Seed contacts
INSERT INTO contacts (id, account_id, first_name, last_name, email, phone)
VALUES
  ('CON-001', 'ACC-001', 'Budi', 'Santoso', 'budi@majubersama.co.id', '+62811-001-001'),
  ('CON-002', 'ACC-002', 'Dewi', 'Rahayu', 'dewi@teknusa.co.id', '+62811-002-002')
ON CONFLICT DO NOTHING;

-- Seed leads
INSERT INTO leads (id, first_name, last_name, email, company, status, source)
VALUES
  ('LEAD-001', 'Andi', 'Wijaya', 'andi@prospect.co.id', 'PT Prospect Baru', 'NEW', 'REFERRAL'),
  ('LEAD-002', 'Siti', 'Nurhaliza', 'siti@newcorp.co.id', 'CV New Corp', 'CONTACTED', 'WEBSITE')
ON CONFLICT DO NOTHING;