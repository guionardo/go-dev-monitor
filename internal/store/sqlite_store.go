package store

import (
	"context"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"

	"github.com/guionardo/go-dev-monitor/internal/logging"
	"github.com/guionardo/go-dev-monitor/internal/repository"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
	"zombiezen.com/go/sqlite/sqlitex"
)

type SqliteStore struct {
	dbFile string
	pool   *sqlitex.Pool
	lock   sync.RWMutex
}

func NewSqliteStore(storeFolder string) (*SqliteStore, error) {
	dataFile := path.Join(storeFolder, "db.sqlite")
	logging.Debug("sqlite_store", slog.String("data_file", dataFile))
	if err := applyMigration(dataFile); err != nil {
		return nil, err
	}
	pool, err := sqlitex.NewPool(dataFile, sqlitex.PoolOptions{
		PoolSize: 10,
		Flags:    sqlite.OpenReadWrite | sqlite.OpenWAL,
	})
	if err != nil {
		return nil, err
	}

	return &SqliteStore{
		dbFile: dataFile,
		pool:   pool,
	}, nil
}

func (s *SqliteStore) BeginPosts(hostName string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := s.pool.Take(ctx)
	if err != nil {
		return err
	}
	defer s.pool.Put(conn)
	sql := "DELETE FROM repos WHERE hostname=?1"
	return sqlitex.Execute(conn, sql, &sqlitex.ExecOptions{
		Args: []any{hostName},
	})
}

func (s *SqliteStore) Post(hostName string, repository *repository.Local) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := s.pool.Take(ctx)
	if err != nil {
		return err
	}
	defer s.pool.Put(conn)

	sql := "INSERT OR REPLACE INTO repos (origin, updated_at, hostname, data) VALUES (?1,?2,?3,?4)"
	data, _ := repository.MarshalJSON()
	err = sqlitex.Execute(conn, sql, &sqlitex.ExecOptions{
		Args: []any{
			repository.Origin,
			time.Now(),
			hostName,
			string(data),
		},
	})
	return err
}

func (s *SqliteStore) GetSummary() (map[string][]*repository.Local, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	conn, err := sqlite.OpenConn(s.dbFile, sqlite.OpenReadWrite|sqlite.OpenWAL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	summary := make(map[string][]*repository.Local)
	stmt := conn.Prep("SELECT origin,updated_at,hostname,data FROM repos ORDER BY origin")
	for {
		if hasRow, err := stmt.Step(); err != nil {
			// ... handle error
		} else if !hasRow {
			break
		}
		origin := stmt.GetText("origin")
		if repos := summary[origin]; len(repos) == 0 {
			summary[origin] = make([]*repository.Local, 0, 1)
		}

		data := stmt.GetText("data")
		repo := &repository.Local{}
		if err = repo.UnmarshalJSON([]byte(data)); err != nil {
			continue
		}

		summary[origin] = append(summary[origin], repo)
	}

	return summary, nil
}

func applyMigration(dbFile string) error {
	schema := sqlitemigration.Schema{
		// Each element of the Migrations slice is applied in sequence. When you
		// want to change the schema, add a new SQL script to this list.
		//
		// Existing databases will pick up at the same position in the Migrations
		// slice as they last left off.
		Migrations: []string{
			`CREATE TABLE "repos" (
	"id"	INTEGER NOT NULL,
	"origin"	TEXT NOT NULL,
	"updated_at"	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"hostname"	TEXT NOT NULL,
	"data"	TEXT NOT NULL,
	PRIMARY KEY("id")
);`,

			`CREATE UNIQUE INDEX "origin_hostname_idx" ON "repos" (
	"origin"	ASC,
	"hostname"	ASC
);`,
		},

		// // The RepeatableMigration is run after all other Migrations if any
		// // migration was run. It is useful for creating triggers and views.
		// RepeatableMigration: "DROP VIEW IF EXISTS bar;\n" +
		// 	"CREATE VIEW bar ( id, name ) AS SELECT id, name FROM foo;\n",
	}

	// Set up a temporary directory to store the database.
	dir, err := os.MkdirTemp("", "sqlitemigration")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	// Open a pool. This does not block, and will start running any migrations
	// asynchronously.
	pool := sqlitemigration.NewPool(dbFile, schema, sqlitemigration.Options{
		Flags: sqlite.OpenReadWrite | sqlite.OpenCreate,
		PrepareConn: func(conn *sqlite.Conn) error {
			// Enable foreign keys. See https://sqlite.org/foreignkeys.html
			return sqlitex.ExecuteTransient(conn, "PRAGMA foreign_keys = ON;", nil)
		},
		OnError: func(e error) {
			logging.Error("migration", e)
		},
	})
	defer func() {
		_ = pool.Close()
	}()
	// Get a connection. This blocks until the migration completes.
	conn, err := pool.Get(context.TODO())
	if err != nil {
		// handle error
		return err
	}
	defer pool.Put(conn)

	// Print the list of schema objects created.
	const listSchemaQuery = `SELECT "type", "name" FROM sqlite_master ORDER BY 1, 2;`

	err = sqlitex.ExecuteTransient(conn, listSchemaQuery, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			logging.Debug("store",
				slog.String("type", stmt.ColumnText(0)),
				slog.String("name", stmt.ColumnText(1)))
			return nil
		},
	})
	return err
}
