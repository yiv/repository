package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DatabaseAdminInfo struct {
	DBUser string
	DBPwd  string
	DBName string
	DBAddr string
	DBPort string
}

type DbRepo struct {
	conn   *sqlx.DB
	logger log.Logger
	info   *DatabaseAdminInfo
}

func NewDBConn(info *DatabaseAdminInfo) (dbConn *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", info.DBUser, info.DBPwd, info.DBAddr, info.DBPort, info.DBName)
	if dbConn, err = sqlx.Open("mysql", dsn); err != nil {
		return nil, err
	}
	if err = dbConn.Ping(); err != nil {
		return nil, err
	}
	return dbConn, nil
}

func CreateTableIfNotExist(dbConn *sqlx.DB, dbUser, dbName string, tableName string, template string) (err error) {
	var (
		res sql.Result
		aff int64
	)
	if res, err = dbConn.Exec(fmt.Sprintf("create table IF NOT EXISTS %s.%s (%s)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;", dbName, tableName, template)); err != nil {
		return
	}
	aff, err = res.RowsAffected()
	if aff > 0 {
		sqlstmt := fmt.Sprintf("grant all privileges on  %s.%s to %s", dbName, tableName, dbUser)
		if _, err = dbConn.Exec(sqlstmt); err != nil {
			return
		}
	}
	return
}

func UseDatabase(dbConn *sqlx.DB, dbName string) (err error) {
	_, err = dbConn.Exec(fmt.Sprintf("use database %s;", dbName))
	return
}
