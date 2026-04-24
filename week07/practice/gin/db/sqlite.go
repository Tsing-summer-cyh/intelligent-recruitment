package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var SQLite *sql.DB

func InitSQLite() {
	var err error
	SQLite, err = sql.Open("sqlite", "./students.db")
	if err != nil {
		log.Fatalf("[DB] SQLite 数据库连接失败: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		grade TEXT NOT NULL
	);`
	if _, err = SQLite.Exec(createTableSQL); err != nil {
		log.Fatalf("[DB] SQLite 创建表失败: %v", err)
	}
	log.Println("✅ SQLite 初始化成功")
}
