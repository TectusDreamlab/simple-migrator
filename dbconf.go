package goose

import (
	"database/sql"
	"net/url"
	"strings"
)

// DBDriver encapsulates the info needed to work with
// a specific database driver
type DBDriver struct {
	Name    string
	DSN     string
	Import  string
	Dialect SqlDialect
}

type DBConf struct {
	MigrationsDir string
	AssetMust     func(path string) []byte
	AssetNames    func() []string
	AssetDir      func(path string) ([]string, error)
	Driver        DBDriver
}

func NewDBConf(driverName, driverDSN string, assetMust func(path string) []byte, assetDir func(path string) ([]string, error), assetNames func() []string, migrationsDir string) *DBConf {
	d := newDBDriver(driverName, driverDSN)
	return &DBConf{
		MigrationsDir: migrationsDir,
		Driver:        d,
		AssetNames:    assetNames,
		AssetDir:      assetDir,
		AssetMust:     assetMust,
	}
}

// Create a new DBDriver and populate driver specific
// fields for drivers that we know about.
// Further customization may be done in NewDBConf
func newDBDriver(name, open string) DBDriver {
	d := DBDriver{
		Name: name,
		DSN:  open,
	}

	switch strings.ToLower(name) {
	case "postgres":
		d.Name = "postgres"
		d.Import = "github.com/lib/pq"
		d.Dialect = &PostgresDialect{}

	case "redshift":
		d.Name = "postgres"
		d.Import = "github.com/lib/pq"
		d.Dialect = &RedshiftDialect{}

	case "mymysql":
		d.Import = "github.com/ziutek/mymysql/godrv"
		d.Dialect = &MySqlDialect{}

	case "mysql":
		d.Import = "github.com/go-sql-driver/mysql"
		d.Dialect = &MySqlDialect{}

	case "sqlite3":
		d.Name = "sqlite3"
		d.Import = "github.com/mattn/go-sqlite3"
		d.Dialect = &Sqlite3Dialect{}
	}

	return d
}

// ensure we have enough info about this driver
func (drv *DBDriver) IsValid() bool {
	return len(drv.Import) > 0 && drv.Dialect != nil
}

// OpenDBFromDBConf wraps database/sql.DB.Open() and configures
// the newly opened DB based on the given DBConf.
//
// Callers must Close() the returned DB.
func openDBFromDBConf(conf *DBConf) (*sql.DB, error) {
	// we depend on time parsing, so make sure it's enabled with the mysql driver
	if conf.Driver.Name == "mysql" {
		i := strings.Index(conf.Driver.DSN, "?")
		if i == -1 {
			i = len(conf.Driver.DSN)
			conf.Driver.DSN = conf.Driver.DSN + "?"
		}
		i++

		q, err := url.ParseQuery(conf.Driver.DSN[i:])
		if err != nil {
			return nil, err
		}
		q.Set("parseTime", "true")

		conf.Driver.DSN = conf.Driver.DSN[:i] + q.Encode()
	}

	return sql.Open(conf.Driver.Name, conf.Driver.DSN)
}
