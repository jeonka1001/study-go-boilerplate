package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"boilerplate/app/config"
	"boilerplate/app/core/helper/logger"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQL opens a GORM connection to MySQL/MariaDB and wires it into the fx lifecycle
// so the underlying *sql.DB is closed cleanly on shutdown.
func NewMySQL(lc fx.Lifecycle, c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.DB.Username, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	if c.DB.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(c.DB.MaxOpenConns)
	}
	if c.DB.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(c.DB.MaxIdleConns)
	}
	if c.DB.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(c.DB.ConnMaxLifetime) * time.Minute)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Zap.Infof("Connecting to MySQL at %s:%s/%s", c.DB.Host, c.DB.Port, c.DB.Name)
			return pingWithContext(ctx, sqlDB)
		},
		OnStop: func(ctx context.Context) error {
			logger.Zap.Infof("Closing MySQL connection")
			return sqlDB.Close()
		},
	})

	return db, nil
}

func pingWithContext(ctx context.Context, sqlDB *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}
