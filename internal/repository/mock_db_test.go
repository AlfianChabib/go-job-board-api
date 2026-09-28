package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mockExpectation struct {
	match        string
	columns      []string
	rows         [][]driver.Value
	rowsAffected int64
	err          error
}

type mockDBState struct {
	mu           sync.Mutex
	expectations []*mockExpectation
}

func (s *mockDBState) Expect(exp mockExpectation) *mockDBState {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expectations = append(s.expectations, &exp)
	return s
}

func (s *mockDBState) pop(query string) (*mockExpectation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queryLower := strings.ToLower(query)
	for i, exp := range s.expectations {
		if exp.match == "" || strings.Contains(queryLower, strings.ToLower(exp.match)) {
			s.expectations = append(s.expectations[:i], s.expectations[i+1:]...)
			return exp, nil
		}
	}
	return nil, fmt.Errorf("unexpected SQL query: %s", query)
}

var (
	driverOnce sync.Once
	dsnCounter uint64
	stateMap   sync.Map
)

type customMockDriver struct{}

func (d *customMockDriver) Open(dsn string) (driver.Conn, error) {
	val, ok := stateMap.Load(dsn)
	if !ok {
		return nil, fmt.Errorf("mock state not found for dsn %s", dsn)
	}
	return &customMockConn{state: val.(*mockDBState)}, nil
}

type customMockConn struct {
	state *mockDBState
}

func (c *customMockConn) Prepare(query string) (driver.Stmt, error) {
	return &customMockStmt{conn: c, query: query}, nil
}

func (c *customMockConn) Close() error { return nil }

func (c *customMockConn) Begin() (driver.Tx, error) {
	return &customMockTx{}, nil
}

func (c *customMockConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return &customMockTx{}, nil
}

func (c *customMockConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	exp, err := c.state.pop(query)
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
	return &customMockRows{columns: cols, rows: exp.rows}, nil
}

func (c *customMockConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	exp, err := c.state.pop(query)
	if err != nil {
		return nil, err
	}
	if exp.err != nil {
		return nil, exp.err
	}
	return driver.RowsAffected(exp.rowsAffected), nil
}

type customMockTx struct{}

func (t *customMockTx) Commit() error   { return nil }
func (t *customMockTx) Rollback() error { return nil }

type customMockStmt struct {
	conn  *customMockConn
	query string
}

func (s *customMockStmt) Close() error  { return nil }
func (s *customMockStmt) NumInput() int { return -1 }

func (s *customMockStmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.conn.ExecContext(context.Background(), s.query, nil)
}

func (s *customMockStmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.conn.QueryContext(context.Background(), s.query, nil)
}

type customMockRows struct {
	columns []string
	rows    [][]driver.Value
	idx     int
}

func (r *customMockRows) Columns() []string { return r.columns }
func (r *customMockRows) Close() error      { return nil }

func (r *customMockRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.idx]
	r.idx++
	for i := range dest {
		if i < len(row) {
			dest[i] = row[i]
		} else {
			dest[i] = nil
		}
	}
	return nil
}

func newMockGormDB(t *testing.T) (*gorm.DB, *mockDBState) {
	driverOnce.Do(func() {
		sql.Register("gorm_unit_mock", &customMockDriver{})
	})

	id := atomic.AddUint64(&dsnCounter, 1)
	dsn := fmt.Sprintf("mock_dsn_%d", id)
	state := &mockDBState{}
	stateMap.Store(dsn, state)
	t.Cleanup(func() {
		stateMap.Delete(dsn)
	})

	sqlDB, err := sql.Open("gorm_unit_mock", dsn)
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

	return gormDB, state
}
