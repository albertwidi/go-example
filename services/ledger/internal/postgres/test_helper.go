// Code is generated by the helper script. DO NOT EDIT.

package postgres

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/albertwidi/pkg/postgres"
	testingpkg "github.com/albertwidi/pkg/testing"
	"github.com/albertwidi/pkg/testing/pgtest"

	"github.com/albertwidi/go-example/internal/env"
)

type TestHelper struct {
	dbName string
	testQueries *Queries
	conn *postgres.Postgres
	pgtestHelper *pgtest.PGTest
	mu sync.Mutex
	// forks is the list of forked helper throughout the test. We need to track the lis of forked helper as we want
	// to track the resource of helper and close them properly.
	forks []*TestHelper
	// fork is a mark that the test helper had been forked, thus several expections should be made when
	// doing several operation like closing connections.
	fork bool
	closed bool
}

func NewTestHelper(ctx context.Context) (*TestHelper, error) {
	th :=&TestHelper{
		dbName: "go_example",
		pgtestHelper: pgtest.New(),
	}
	q, err := th.prepareTest(ctx)
	if err != nil {
		return nil, err
	}
	th.testQueries = q
	return th, nil
}

func (th *TestHelper) Queries() *Queries{
	return th.testQueries
}

// prepareTest prepares the designated postgres database by creating the database and applying the schema. The function returns a postgres connection
// to the database that can be used for testing purposes.
func (th *TestHelper) prepareTest(ctx context.Context) (*Queries, error) {
	pgDSN := env.GetEnvOrDefault("TEST_PG_DSN", "postgres://postgres:postgres@localhost:5432/")
	if err := pgtest.CreateDatabase(ctx, pgDSN, th.dbName, false); err != nil {
		return nil, err
	}

	// Create a new connection with the correct database name.
	config, err := postgres.NewConfigFromDSN(pgDSN)
	if err != nil {
		return nil, err
	}
	config.DBName = th.dbName
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
	// Assgign the connection for the test helper.
	th.conn = testConn
	return New(testConn), nil
}

// Close closes all connections from the test helper.
func (th *TestHelper) Close() error {
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
			errors.Join(err ,errClose)
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
		return pgtest.DropDatabase(context.Background(), pg, th.dbName)
	}
	if err == nil {
		th.closed = true
	}
	return err
}

// ForkPostgresSchema forks the sourceSchema with the underlying connection inside the Queries. The function will return a new connection
// with default search_path into the new schema. The schema name currently is random and cannot be defined by the user.
func (th *TestHelper) ForkPostgresSchema(ctx context.Context, q *Queries, sourceSchema string) (*TestHelper, error) {
	th.mu.Lock()
	defer th.mu.Unlock()
	if th.fork {
		return nil, errors.New("cannot fork the schema from a forked test helper, please use the original test helper")
	}
	pg , err:= th.pgtestHelper.ForkSchema(ctx, q.db, sourceSchema)
	if err != nil {
		return nil, err
	}
	newTH := &TestHelper{
		dbName: th.dbName,
		conn: pg,
		testQueries: New(pg),
		pgtestHelper: th.pgtestHelper,
		fork: true,
	}
	// Append the forks to the origin
	th.forks = append(th.forks, newTH)
	return newTH, nil
}
