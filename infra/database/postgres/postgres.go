package database

import (
	"fmt"
	"os"
	"time"

	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	slog "golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateConnectionPool(db *gorm.DB) error {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	conn, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	conn.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	conn.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	conn.SetConnMaxLifetime(time.Hour)

	log.Info("database connection pool created")

	return nil
}

func RunMigrations(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	return err
}

func GeneratePostgresURI(env *config.Config) string {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var (
		dbUrl    = env.DATABASE_URL
		host     = env.PG_HOST
		port     = env.PG_PORT
		dbname   = env.PG_NAME
		user     = env.PG_USER
		password = env.PG_PASS
		sslmode  = env.PG_SSLMODE
	)
	if env.ENV == config.Development {
		dbUrl = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	}

	log.Info("postgres uri generated successfully", "url", dbUrl)

	return dbUrl
}

func NewPostgres(config *config.Config) (*gorm.DB, error) {

	var (
		db  *gorm.DB
		err error
	)

	db, err = gorm.Open(postgres.Open(GeneratePostgresURI(config)), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db = db.Debug()

	return db, nil
}
