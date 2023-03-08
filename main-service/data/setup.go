package data

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Log *zap.Logger

/* make database connection set up */
func makeDbConnection() *gorm.DB{

    time.Sleep(time.Duration(5000 * time.Millisecond))

	dsn := "docker:docker@tcp(mysql_host:3306)/test_database?multiStatements=True&charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "bulkmsg:xCu34REN8hqB6Ydh.@tcp(10.120.12.80:3306)/bulk_message2"


	sqlDB, error := sql.Open("mysql", dsn)
    if error != nil {
		Log.Fatal(error.Error())
	}

    db, err := gorm.Open(mysql.New(mysql.Config{
        Conn: sqlDB,
    }), &gorm.Config{})

	if err != nil {
			Log.Warn("database not yet ready!")
		} else {
			Log.Info("connected to database!")
		}

	if db == nil {
		Log.Panic("can't connect to database!")
	}

  db.AutoMigrate(&ImportData{},&ImportHistory{},&MessageTemplate{},&BulkMessage{},&BulkMessagesBox{})
  return db
}

/* execute sqlite file */
func executeSqliteFile(filePath string) ([]string, *sql.Rows) {

	sqlitedb, err := sql.Open("sqlite3", filePath)

    if err != nil {
        Log.Fatal(err.Error())
    }

	rows, err := sqlitedb.Query("SELECT * FROM sample1")

	if err != nil {
        Log.Fatal(err.Error())
    }

	col, err := rows.Columns()

	if err != nil {
        Log.Fatal(err.Error())
    }

	return col,rows
}

func MakeLogger() {

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "./tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "Message",
	    "levelKey": "Level",
	    "levelEncoder": "capital",
		"timeKey": "Time",
		"timeEncoder": "iso8601"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	Log, _ = cfg.Build()
	defer Log.Sync()
	Log.Info("logger construction succeeded")
}

