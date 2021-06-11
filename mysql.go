package dbb

import (
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	DefMysqlCreateTableVersion string = `CREATE TABLE __version (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		ctime TIMESTAMP NOT NULL DEFAULT current_timestamp,
		mtime TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
		ver VARCHAR(15) NOT NULL,
	
		UNIQUE ver (ver)
	) ENGINE InnoDB COLLATE 'utf8mb4_unicode_ci' CHARACTER SET 'utf8mb4';`

	DefMysqlInsertVersion string = `INSERT INTO __version (ver) VALUES (?);`
)

//================================================================
// mysql
//================================================================
func (s *Schema) mysqlBuild(db *sqlx.DB, pw string, vhs *VerHandlers) error {
	if err := s.mysqlCreateDBIfNotExists(db); err != nil {
		return err
	}

	if err := s.mysqlCreateUserIfNotExists(db, pw); err != nil {
		return err
	}

	ver := ""
	if err := db.Get(&ver, `SELECT ver FROM __version WHERE id = (SELECT MAX(id) FROM __version)`); err != nil {
		return err
	}

	for vh := vhs.getVerHandler(ver); vh != nil; vh = vhs.getVerHandler(ver) {
		if handlerFunc, ok := vh.HandlerFunc.(func(*sqlx.DB) error); !ok {
			return ErrHandlerFuncAssertFail
		} else if err := handlerFunc(db); err != nil {
			return err
		}

		if _, err := db.Exec(DefMysqlInsertVersion, vh.ToVer); err != nil {
			return err
		}

		log.Println("Database schema from [v" + ver + "] to [v" + vh.ToVer + "] upgraded.")
		ver = vh.ToVer
	}

	log.Println("Database Schema upgrade completed.")
	return nil
}

func (s *Schema) mysqlCreateDBIfNotExists(db *sqlx.DB) error {
	exists := false

	if err := db.Get(&exists, `SELECT IF(COUNT(1), true, false) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?`, s.DBName); err != nil {
		return err
	}

	if !exists {
		if _, err := db.Exec(`CREATE DATABASE ` + s.DBName + ` COLLATE 'utf8mb4_unicode_ci' CHARACTER SET 'utf8mb4';`); err != nil {
			return err
		}

		if _, err := db.Exec(`use ` + s.DBName + `;`); err != nil {
			return err
		}

		if _, err := db.Exec(DefMysqlCreateTableVersion); err != nil {
			return err
		}

		if _, err := db.Exec(DefMysqlInsertVersion, DefInitVer); err != nil {
			return err
		}
	} else {
		if _, err := db.Exec(`use ` + s.DBName + `;`); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) mysqlCreateUserIfNotExists(db *sqlx.DB, pw string) error {
	if _, err := db.Exec(`CREATE USER IF NOT EXISTS '` + s.AppUser + `'@'` + s.AppUserHost + `' IDENTIFIED BY '` + pw + `';`); err != nil {
		return err
	}

	if _, err := db.Exec(`GRANT ALL PRIVILEGES ON ` + s.DBName + `.* TO '` + s.AppUser + `'@'` + s.AppUserHost + `';`); err != nil {
		return err
	}

	if _, err := db.Exec(`FLUSH PRIVILEGES;`); err != nil {
		return err
	}

	return nil
}
