CREATE TABLE applications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_id UUID NOT NULL,
  candidate_id UUID NOT NULL,
  resume_url VARCHAR(255) NOT NULL,
  resume_filename VARCHAR(255),
  cover_letter TEXT,
  expected_salary BIGINT,
  status VARCHAR(50) NOT NULL DEFAULT 'APPLIED',
  recruiter_notes TEXT,
  rejection_reason VARCHAR(255),
  withdrawn_reason VARCHAR(255),
  applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,

  CONSTRAINT fk_applications_job FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
  CONSTRAINT fk_applications_candidate FOREIGN KEY (candidate_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT uq_applications_job_candidate UNIQUE (job_id, candidate_id)
);

-- Indexing for performance optimization
CREATE INDEX idx_applications_job_id ON applications (job_id);
CREATE INDEX idx_applications_candidate_id ON applications (candidate_id);
CREATE INDEX idx_applications_status ON applications (status);
CREATE INDEX idx_applications_job_status ON applications (job_id, status);
CREATE INDEX idx_applications_created_at ON applications (created_at DESC);
CREATE INDEX idx_applications_deleted_at ON applications (deleted_at);
