package dbb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type DBI struct {
	*sqlx.DB
	Type          string
	Host          string
	Port          string
	AdminUser     string
	AdminPassword string
}

func AdminConn(typ, host, port, adminUser, adminPw string) (*DBI, error) {
	dbi := &DBI{
		Type:          typ,
		Host:          host,
		Port:          port,
		AdminUser:     adminUser,
		AdminPassword: adminPw,
	}

	return dbi, dbi.open()
}

func (dbi *DBI) Build(dbName, appUser, appUserHost, pw string, vhs *VerHandlers) error {
	s := &Schema{
		DBName:      dbName,
		AppUser:     appUser,
		AppUserHost: appUserHost,
	}

	switch dbi.Type {
	case "mysql":
		return s.mysqlBuild(dbi.DB, pw, vhs)
	default:
	}

	return ErrNothingChanged
}

func (dbi *DBI) CloseDB() {
	if dbi.DB != nil {
		dbi.DB.Close()
	}
}

//================================================================
//
//================================================================
func (dbi *DBI) open() error {
	if db, err := sqlx.Open(dbi.Type, dbi.protocol()); err != nil {
		return err
	} else {
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(time.Duration(60) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(60) * time.Second)
		dbi.DB = db
		return nil
	}
}

func (dbi *DBI) protocol() string {
	protocol := ""
	switch dbi.Type {
	case "mysql":
		protocol = fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			dbi.AdminUser,
			dbi.AdminPassword,
			dbi.Host,
			dbi.Port,
		)
	default:
	}

	return protocol
}
