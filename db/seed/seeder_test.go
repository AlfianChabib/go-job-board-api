package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seedExpectation struct {
	match        string
	columns      []string
	rows         [][]driver.Value
	rowsAffected int64
	err          error
}

type seedMockDriver struct {
	mu     sync.Mutex
	queues map[string][]seedExpectation
}

var (
	seedDrvOnce sync.Once
	seedDrvInst = &seedMockDriver{queues: make(map[string][]seedExpectation)}
)

func newSeedMockDB(t *testing.T) (*gorm.DB, *seedQueue) {
	t.Helper()
	seedDrvOnce.Do(func() {
		sql.Register("seed_unit_mock", seedDrvInst)
	})

	dsn := t.Name() + "_" + uuid.NewString()
	q := &seedQueue{dsn: dsn, drv: seedDrvInst}

	sqlDB, err := sql.Open("seed_unit_mock", dsn)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return gormDB, q
}

type seedQueue struct {
	dsn string
	drv *seedMockDriver
}

func (q *seedQueue) Expect(exp seedExpectation) {
	q.drv.mu.Lock()
	defer q.drv.mu.Unlock()
	q.drv.queues[q.dsn] = append(q.drv.queues[q.dsn], exp)
}

func (d *seedMockDriver) pop(dsn, query string) (seedExpectation, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	q := d.queues[dsn]
	if len(q) == 0 {
		return seedExpectation{}, fmt.Errorf("unexpected SQL query: %s", query)
	}
	exp := q[0]
	d.queues[dsn] = q[1:]
	if exp.match != "" && !strings.Contains(strings.ToLower(query), strings.ToLower(exp.match)) {
		return seedExpectation{}, fmt.Errorf("query %q did not match expected substring %q", query, exp.match)
	}
	return exp, nil
}

func (d *seedMockDriver) Open(name string) (driver.Conn, error) {
	return &seedConn{dsn: name, drv: d}, nil
}

type seedConn struct {
	dsn string
	drv *seedMockDriver
}

func (c *seedConn) Prepare(query string) (driver.Stmt, error) {
	return &seedStmt{conn: c, query: query}, nil
}
func (c *seedConn) Close() error              { return nil }
func (c *seedConn) Begin() (driver.Tx, error) { return &seedTx{}, nil }

func (c *seedConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	exp, err := c.drv.pop(c.dsn, query)
	if err != nil {
		return nil, err
	}
	if exp.err != nil {
		return nil, exp.err
	}
	cols := exp.columns
	if len(cols) == 0 {
		cols = []string{"id"}
	}
	return &seedRows{columns: cols, rows: exp.rows}, nil
}

func (c *seedConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	exp, err := c.drv.pop(c.dsn, query)
	if err != nil {
		return nil, err
	}
	if exp.err != nil {
		return nil, exp.err
	}
	return driver.RowsAffected(exp.rowsAffected), nil
}

type seedTx struct{}

func (t *seedTx) Commit() error   { return nil }
func (t *seedTx) Rollback() error { return nil }

type seedStmt struct {
	conn  *seedConn
	query string
}

func (s *seedStmt) Close() error  { return nil }
func (s *seedStmt) NumInput() int { return -1 }
func (s *seedStmt) Exec(_ []driver.Value) (driver.Result, error) {
	return s.conn.ExecContext(context.Background(), s.query, nil)
}
func (s *seedStmt) Query(_ []driver.Value) (driver.Rows, error) {
	return s.conn.QueryContext(context.Background(), s.query, nil)
}

type seedRows struct {
	columns []string
	rows    [][]driver.Value
	idx     int
}

func (r *seedRows) Columns() []string { return r.columns }
func (r *seedRows) Close() error      { return nil }
func (r *seedRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.idx]
	r.idx++
	for i := range dest {
		if i < len(row) {
			dest[i] = row[i]
		}
	}
	return nil
}

func setupTempSeedDir(t *testing.T) string {
	t.Helper()
	origDir, err := os.Getwd()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "db", "seed"), 0755))
	require.NoError(t, os.Chdir(tmpDir))

	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})
	return filepath.Join(tmpDir, "db", "seed")
}

func TestSeedSkills(t *testing.T) {
	t.Run("skips when skills-dataset.csv does not exist", func(t *testing.T) {
		_ = setupTempSeedDir(t)
		db, _ := newSeedMockDB(t)
		assert.NotPanics(t, func() {
			SeedSkills(db)
		})
	})

	t.Run("parses CSV, deduplicates skills, and inserts in batches", func(t *testing.T) {
		seedDir := setupTempSeedDir(t)
		csvContent := strings.Join([]string{
			"id,name,label,abbreviation",
			"1,golang,,",
			"2,golang,Go Language,GO",
			"3,short_row",
			"4,   ,Empty Name,EN",
			"5,docker,Docker,DKR",
		}, "\n")
		require.NoError(t, os.WriteFile(filepath.Join(seedDir, "skills-dataset.csv"), []byte(csvContent), 0644))

		db, mock := newSeedMockDB(t)
		mock.Expect(seedExpectation{
			match:        "insert into \"skills\"",
			rowsAffected: 2,
		})

		SeedSkills(db)
	})
}

func TestSeedRecruiters(t *testing.T) {
	seedDir := setupTempSeedDir(t)
	longPass := strings.Repeat("p", 80) // > 72 bytes causes bcrypt error
	recruitersJSON := fmt.Sprintf(`[
		{
			"name": "Recruiter One",
			"email": "r1@example.com",
			"password": "Password123!",
			"company": {
				"name": "Acme Corp",
				"industry": "Tech",
				"employee_size": "51-200",
				"location": "Jakarta"
			}
		},
		{
			"name": "Recruiter Fail Create",
			"email": "fail@example.com",
			"password": "Password123!",
			"company": {
				"name": "Fail Corp",
				"industry": "Tech",
				"employee_size": "1-10",
				"location": "Bandung"
			}
		},
		{
			"name": "Recruiter Long Pass",
			"email": "long@example.com",
			"password": %q,
			"company": {
				"name": "Long Corp",
				"industry": "Tech",
				"employee_size": "1-10",
				"location": "Surabaya"
			}
		}
	]`, longPass)
	require.NoError(t, os.WriteFile(filepath.Join(seedDir, "mock_recruiters.json"), []byte(recruitersJSON), 0644))

	db, mock := newSeedMockDB(t)

	// Recruiter 1: user not found -> create user -> create auth -> company not found -> create company
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"users\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})
	mock.Expect(seedExpectation{match: "insert into \"auth\"", rowsAffected: 1})
	mock.Expect(seedExpectation{match: "select * from \"companies\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"companies\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})

	// Recruiter 2: user not found -> create user fails
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{match: "insert into \"users\"", err: errors.New("insert user failed")})

	// Recruiter 3: user not found -> create user succeeds -> password hash fails (>72 bytes)
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"users\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})

	SeedRecruiters(db)
}

func TestSeedJobs(t *testing.T) {
	seedDir := setupTempSeedDir(t)
	jobsJSON := `[
		{
			"company_name": "Missing Corp",
			"title": "Ghost Job",
			"description": "Desc",
			"employment_type": "Full-time",
			"work_mode": "Remote",
			"location": "Jakarta",
			"currency": "",
			"is_salary_negotiable": true,
			"status": "Open"
		},
		{
			"company_name": "Acme Corp",
			"title": "Backend Engineer",
			"description": "Build APIs",
			"employment_type": "Full-time",
			"work_mode": "Remote",
			"location": "Jakarta",
			"currency": "",
			"is_salary_negotiable": true,
			"status": "Open"
		},
		{
			"company_name": "Acme Corp",
			"title": "Failing Job",
			"description": "Fail",
			"employment_type": "Full-time",
			"work_mode": "Onsite",
			"location": "Jakarta",
			"currency": "USD",
			"is_salary_negotiable": false,
			"status": "Open"
		}
	]`
	require.NoError(t, os.WriteFile(filepath.Join(seedDir, "mock_jobs.json"), []byte(jobsJSON), 0644))

	db, mock := newSeedMockDB(t)
	companyID := uuid.New()

	// Job 1: Company not found
	mock.Expect(seedExpectation{match: "select * from \"companies\"", rows: nil})

	// Job 2: Company found -> existingJob not found -> create job succeeds
	mock.Expect(seedExpectation{
		match:   "select * from \"companies\"",
		columns: []string{"id", "name"},
		rows:    [][]driver.Value{{companyID.String(), "Acme Corp"}},
	})
	mock.Expect(seedExpectation{match: "select * from \"jobs\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"jobs\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})

	// Job 3: Company found -> existingJob not found -> create job fails
	mock.Expect(seedExpectation{
		match:   "select * from \"companies\"",
		columns: []string{"id", "name"},
		rows:    [][]driver.Value{{companyID.String(), "Acme Corp"}},
	})
	mock.Expect(seedExpectation{match: "select * from \"jobs\"", rows: nil})
	mock.Expect(seedExpectation{match: "insert into \"jobs\"", err: errors.New("create job error")})

	SeedJobs(db)
}

func TestSeedCandidates(t *testing.T) {
	seedDir := setupTempSeedDir(t)
	longPass := strings.Repeat("p", 80)
	candidatesJSON := fmt.Sprintf(`[
		{
			"name": "Candidate One",
			"email": "c1@example.com",
			"password": "Password123!",
			"profile": {
				"headline": "Backend Dev"
			},
			"skills": ["golang"],
			"experiences": [
				{
					"company_name": "Tech Co",
					"position": "Developer",
					"start_date": "2022-01-01",
					"end_date": "2023-01-01",
					"is_current": false
				}
			]
		},
		{
			"name": "Candidate Fail User",
			"email": "failuser@example.com",
			"password": "Password123!",
			"profile": {},
			"skills": [],
			"experiences": []
		},
		{
			"name": "Candidate Long Pass",
			"email": "longpass@example.com",
			"password": %q,
			"profile": {},
			"skills": [],
			"experiences": []
		},
		{
			"name": "Candidate Fail Profile",
			"email": "failprof@example.com",
			"password": "Password123!",
			"profile": {},
			"skills": [],
			"experiences": []
		}
	]`, longPass)
	require.NoError(t, os.WriteFile(filepath.Join(seedDir, "mock_candidates.json"), []byte(candidatesJSON), 0644))

	db, mock := newSeedMockDB(t)

	// Candidate 1: create user -> create auth -> create profile -> find skill -> create profile_skill -> find experience -> create experience
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"users\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})
	mock.Expect(seedExpectation{match: "insert into \"auth\"", rowsAffected: 1})
	mock.Expect(seedExpectation{match: "select * from \"profiles\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"profiles\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})
	mock.Expect(seedExpectation{
		match:   "select * from \"skills\"",
		columns: []string{"id", "name"},
		rows:    [][]driver.Value{{uuid.NewString(), "golang"}},
	})
	mock.Expect(seedExpectation{match: "insert into \"profile_skills\"", rowsAffected: 1})
	mock.Expect(seedExpectation{match: "select * from \"experiences\"", rows: nil})
	mock.Expect(seedExpectation{match: "insert into \"experiences\"", rowsAffected: 1})

	// Candidate 2: user not found -> create user fails
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{match: "insert into \"users\"", err: errors.New("create user error")})

	// Candidate 3: user not found -> create user succeeds -> hash fails
	mock.Expect(seedExpectation{match: "select * from \"users\"", rows: nil})
	mock.Expect(seedExpectation{
		match:   "insert into \"users\"",
		columns: []string{"id"},
		rows:    [][]driver.Value{{uuid.NewString()}},
	})

	// Candidate 4: user found -> profile not found -> create profile fails
	mock.Expect(seedExpectation{
		match:   "select * from \"users\"",
		columns: []string{"id", "email"},
		rows:    [][]driver.Value{{uuid.NewString(), "failprof@example.com"}},
	})
	mock.Expect(seedExpectation{match: "select * from \"profiles\"", rows: nil})
	mock.Expect(seedExpectation{match: "insert into \"profiles\"", err: errors.New("create profile error")})

	SeedCandidates(db)
}
