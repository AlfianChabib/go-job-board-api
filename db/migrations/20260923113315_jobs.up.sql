CREATE TABLE jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  requirements TEXT,
  employment_type VARCHAR(50) NOT NULL,
  work_mode VARCHAR(50) NOT NULL,
  location VARCHAR(255) NOT NULL,
  min_salary BIGINT,
  max_salary BIGINT,
  currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
  is_salary_negotiable BOOLEAN NOT NULL DEFAULT false,
  status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_jobs_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Indexing for performance optimization
CREATE INDEX idx_jobs_company_id ON jobs (company_id);
CREATE INDEX idx_jobs_status ON jobs (status);
CREATE INDEX idx_jobs_employment_type ON jobs (employment_type);
CREATE INDEX idx_jobs_work_mode ON jobs (work_mode);
CREATE INDEX idx_jobs_location ON jobs (location);
CREATE INDEX idx_jobs_created_at ON jobs (created_at DESC);

