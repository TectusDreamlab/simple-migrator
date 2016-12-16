# Simple Migrator

A small library that is used to handle migrations created by [goose](https://bitbucket.org/liamstask/goose) that are embedded in to the application using [go-bindata](https://github.com/jteeuwen/go-bindata).

The purpose is to make sure that whenever we start the application, we make sure that there're no pending migrations to be apply to the database, so that the source code and the db version will not get out of sync accidently.

## Features
- Can read SQL migrations embedded into application.
- Can check whether there're pending migrations and apply.


## Usage
```
import migrator "github.com/WUMUXIAN/simple-migrator"

// Handle DB migrations.
conf := migrator.NewDBConf("mysql", "root:mx@tcp(localhost:3306)/dreamlab?charset=utf8 ", MustAsset, AssetDir, AssetNames, "db/migrations")
version, _ := migrator.GetMostRecentDBVersion(conf)
err := migrator.RunMigrations(conf, version)
if err != nil {
	panic(err)
} else {
	fmt.Println("DB Migration Completed")
}
```