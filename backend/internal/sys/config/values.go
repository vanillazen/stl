package config

var Key = newCfgKeyReg()

func newCfgKeyReg() *cfgKeyReg {
	return &cfgKeyReg{
		// API Server

		APIServerHost:     "http.api.server.host",
		APIServerPort:     "http.api.server.port",
		APIServerTimeout:  "http.api.server.shutdown.timeout.secs",
		APIErrorExposeInt: "api.errors.expose.internal",

		// Postgres

		PgUser:   "db.pg.user",
		PgPass:   "db.pg.pass",
		PgDB:     "db.pg.database",
		PgHost:   "db.pg.host",
		PgPort:   "db.pg.port",
		PgSchema: "db.pg.schema",
		PgSSL:    "db.pg.sslmode",

		// SQLite

		SQLiteFilePath: "db.sqlite.filepath",
		SQLiteUser:     "db.sqlite.user",
		SQLitePass:     "db.sqlite.pass",
		SQLiteDB:       "db.sqlite.database",
		SQLiteHost:     "db.sqlite.host",
		SQLitePort:     "db.sqlite.port",
		SQLiteSchema:   "db.sqlite.schema",
		SQLiteSSL:      "db.sqlite.sslmode",
	}
}

type cfgKeyReg struct {
	APIServerHost     string
	APIServerPort     string
	APIServerTimeout  string
	APIErrorExposeInt string

	// Postgres

	PgUser   string
	PgPass   string
	PgDB     string
	PgHost   string
	PgPort   string
	PgSchema string
	PgSSL    string

	// SQLite

	SQLiteFilePath string
	SQLiteUser     string
	SQLitePass     string
	SQLiteDB       string
	SQLiteHost     string
	SQLitePort     string
	SQLiteSchema   string
	SQLiteSSL      string
}
