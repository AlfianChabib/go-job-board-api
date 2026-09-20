CREATE TABLE companies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recruiter_id UUID NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  logo_url VARCHAR(255),
  banner_url VARCHAR(255),
  website VARCHAR(255),
  industry VARCHAR(255) NOT NULL,
  employee_size VARCHAR(50) NOT NULL,
  description TEXT,
  location VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_companies_recruiter FOREIGN KEY (recruiter_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Indexing for performance optimization
CREATE INDEX idx_companies_recruiter_id ON companies (recruiter_id);
CREATE INDEX idx_companies_name ON companies (name);
CREATE INDEX idx_companies_industry ON companies (industry);
