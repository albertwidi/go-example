// pghelper is a package to provide a test struct/object for postgres database. The package is design to match the function signature
// of intrernal postgres package in this project.

package pghelper

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/albertwidi/pkg/postgres"
	testingpkg "github.com/albertwidi/pkg/testing"
	"github.com/albertwidi/pkg/testing/pgtest"

	"github.com/albertwidi/go-example/internal/env"
)

// PGQuery interface defines the type that use this package should implements PGQuery.
type PGQuery interface {
	// Postgres returns the postgres connection object from github.com/albertwidi/pkg/postgres. We need the connection
	// because we will create the query object from the connection.
	Postgres() *postgres.Postgres
}

type Helper[T PGQuery] struct {
	dbName       string
	conn         *postgres.Postgres
	pgtestHelper *pgtest.PGTest

	queriesFn   func(*postgres.Postgres) T
	testQueries T

	mu sync.Mutex
	// forks is the list of forked helper throughout the test. We need to track the lis of forked helper as we want
	// to track the resource of helper and close them properly.
	forks []*Helper[T]
	// fork is a mark that the test helper had been forked, thus several expections should be made when
	// doing several operation like closing connections.
	fork   bool
	closed bool
}

func New[T PGQuery](ctx context.Context, dbName string, fn func(*postgres.Postgres) T) (*Helper[T], error) {
	if !testing.Testing() {
		return nil, errors.New("can only be used in test")
	}
	th := &Helper[T]{
		dbName:       dbName,
		pgtestHelper: pgtest.New(),
	}
	pg, err := prepareTest(ctx, dbName)
	if err != nil {
		return nil, err
	}
	th.conn = pg
	th.testQueries = fn(pg)
	th.queriesFn = fn
	return th, nil
}

func (th *Helper[T]) Queries() T {
	return th.testQueries
}

// Close closes all connections from the test helper.
func (th *Helper[T]) Close() error {
	th.mu.Lock()
	defer th.mu.Unlock()
	if th.closed {
		return nil
	}

	var err error
	if th.conn != nil {
		errClose := th.conn.Close()
		if errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	// If not a fork, then we should close all the connections in the test helper as it will closes all connections
	// to the forked schemas. But in fork, we should avoid this as we don't want to control this from forked test helper.
	if !th.fork {
		errClose := th.pgtestHelper.Close()
		if errClose != nil {
			errors.Join(err, errClose)
		}
		// Closes all the forked helper, this closes the postgres connection in each helper.
		for _, forkedHelper := range th.forks {
			if err := forkedHelper.Close(); err != nil {
				return err
			}
		}
		// Drop the database after test so we will always have a fresh database when we start the test.
		config := th.conn.Config()
		config.DBName = ""
		pg, err := postgres.Connect(context.Background(), config)
		if err != nil {
			return err
		}
		defer pg.Close()
		err = pgtest.DropDatabase(context.Background(), pg, th.dbName)
	}
	if err == nil {
		th.closed = true
	}
	return err
}

// ForkPostgresSchema forks the sourceSchema with the underlying connection inside the Queries. The function will return a new connection
// with default search_path into the new schema. The schema name currently is random and cannot be defined by the user.
func (th *Helper[T]) ForkPostgresSchema(ctx context.Context, q T) (*Helper[T], error) {
	if th.fork {
		return nil, errors.New("cannot fork the schema from a forked test helper, please use the original test helper")
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	if th.closed {
		return nil, errors.New("cannot create a fork from closed test helper")
	}
	pg, err := th.pgtestHelper.ForkSchema(ctx, q.Postgres(), q.Postgres().DefaultSearchPath())
	if err != nil {
		return nil, err
	}
	newTH := &Helper[T]{
		dbName:       th.dbName,
		conn:         pg,
		testQueries:  th.queriesFn(pg),
		pgtestHelper: th.pgtestHelper,
		fork:         true,
	}
	// Append the forks to the origin
	th.forks = append(th.forks, newTH)
	return newTH, nil
}

// DefaultSearchPath returns the default PostgreSQL search path. This helper function invoke the pg.DefaultSearchPath
// to do this. This function added to avoid the user/client to go deeper to the postgres object to invoke this function.
func (th *Helper[T]) DefaultSearchPath() string {
	return th.conn.DefaultSearchPath()
}

// CloseFunc is a helper function to close the test helper via testing.T.Cleanup.
func (th *Helper[T]) CloseFunc(t *testing.T) func() {
	return func() {
		if err := th.Close(); err != nil {
			t.Log(err)
		}
	}
}

// prepareTest prepares the designated postgres database by creating the database and applying the schema. The function returns a postgres connection
// to the database that can be used for testing purposes.
func prepareTest(ctx context.Context, dbName string) (*postgres.Postgres, error) {
	// TEST_PG_DSN can be used to set different DSN for flexible test setup.
	pgDSN := env.GetEnvOrDefault("TEST_PG_DSN", "postgres://postgres:postgres@localhost:5432/")
	if err := pgtest.CreateDatabase(ctx, pgDSN, dbName, false); err != nil {
		return nil, err
	}

	// Create a new connection with the correct database name.
	config, err := postgres.NewConfigFromDSN(pgDSN)
	if err != nil {
		return nil, err
	}
	config.DBName = dbName
	// Connect to the PostgreSQL with the configuration.
	testConn, err := postgres.Connect(context.Background(), config)
	if err != nil {
		return nil, err
	}
	// Read the schema and apply the schema.
	repoRoot, err := testingpkg.RepositoryRoot()
	if err != nil {
		return nil, err
	}
	out, err := os.ReadFile(filepath.Join(repoRoot, "database/ledger/schema.sql"))
	if err != nil {
		return nil, err
	}
	_, err = testConn.Exec(context.Background(), string(out))
	if err != nil {
		return nil, err
	}
	return testConn, nil
}
