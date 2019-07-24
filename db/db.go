package db

type DB struct{}

var db *DB

func Init() {

}

func GetDB() *DB {
	return db
}
