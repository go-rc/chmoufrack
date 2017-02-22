package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Exercise (
	ID int NOT NULL AUTO_INCREMENT,
	name varchar(255) NOT NULL,
    comment text,
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Warmup (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	exerciseID integer,
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Warmdown (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	exerciseID integer,
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Intervals (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	laps tinyint NOT NULL,
    length INTEGER NOT NULL,
	percentage tinyint NOT NULL,
	rest text,
	effort_type varchar(32) DEFAULT "distance",
	effort text, -- storing time in there
	repeatID integer,
	exerciseID integer,
    CHECK(repeatID is not NULL or exerciseID is not NULL),
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Repeats  (
	id int NOT NULL AUTO_INCREMENT,
	repeats tinyint,
	position tinyint DEFAULT 0,
	exerciseID integer,
    PRIMARY KEY(ID)
);
`

//TODO: remove
var SQLDropTable = `
SET FOREIGN_KEY_CHECKS = 0;
SET GROUP_CONCAT_MAX_LEN=32768;
SET @tables = NULL;
SELECT GROUP_CONCAT(table_name) INTO @tables
  FROM information_schema.tables
  WHERE table_schema = (SELECT DATABASE());
SELECT IFNULL(@tables,'dummy') INTO @tables;

SET @tables = CONCAT('DROP TABLE IF EXISTS ', @tables);
PREPARE stmt FROM @tables;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
SET FOREIGN_KEY_CHECKS = 1;`

var SQLresetDB = `
	DELETE FROM Exercise;
	DELETE FROM Warmup;
	DELETE FROM Warmdown;
	DELETE FROM Intervals;
	DELETE FROM Repeats;`

type ArgsMap map[string]interface{}

func createSampleExercise(
	exerciceName, warmupEffort, warmdownEffort string,
	length2 int) (e Exercise) {
	var steps Steps

	step1 := Step{
		Type:       "warmup",
		Effort:     warmupEffort,
		EffortType: "distance",
	}
	steps = append(steps, step1)

	step2 := Step{
		Laps:       3,
		Length:     length2,
		Percentage: 90,
		Type:       "interval",
		EffortType: "distance",
	}
	steps = append(steps, step2)

	step3 := Step{
		Effort:     warmdownEffort,
		Type:       "warmdown",
		EffortType: "distance",
	}
	steps = append(steps, step3)

	e = Exercise{
		Name:    exerciceName,
		Comment: "NoComment",
		Steps:   steps,
	}
	return
}

func SQLInsertOrUpdate(table string, id int, am ArgsMap) (lastid int, err error) {
	var c int
	var res sql.Result
	var begin, query string

	var keys []interface{} = make([]interface{}, 0)
	var values []interface{} = make([]interface{}, 0)
	for k, v := range am {
		keys = append(keys, k)
		values = append(values, v)
	}

	if id != 0 {
		begin = "REPLACE INTO "
	} else {
		begin = "INSERT INTO "
	}

	query = begin + table + "("
	c = 1
	for _, k := range keys {
		query += k.(string)
		if c != len(am) {
			query += ","
		}
		c += 1
	}
	query += ") VALUES ("
	c = 1
	for range keys {
		query += `?`
		if c != len(am) {
			query += ","
		}
		c += 1
	}
	query += ");"

	res, err = sqlTX(query, values...)
	if err != nil {
		return
	}

	n, _ := res.LastInsertId()
	lastid = int(n)
	return
}

func sqlTX(query string, args ...interface{}) (res sql.Result, err error) {
	tx, err := DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return

	}

	defer stmt.Close()
	res, err = stmt.Exec(args...)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func DBConnect(dbconnection string, reset bool) (err error) {
	DB, err = sql.Open("mysql", dbconnection+"?multiStatements=true")

	if err != nil {
		return
	}

	if reset {
		_, err = DB.Exec(SQLDropTable)
		if err != nil {
			return
		}
	}

	_, err = DB.Exec(sqlTable)
	return
}

func InitFixturesDB() (err error) {
	_, err = DB.Exec(SQLresetDB)

	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	var repeatSteps Steps
	repeatStep := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep)

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)

	_, err = AddExercise(e)
	if err != nil {
		return
	}
	return
}
