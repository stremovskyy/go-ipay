module github.com/stremovskyy/go-ipay/examples/status_recorder

go 1.25.0

require (
	github.com/stremovskyy/go-ipay v1.17.1
	github.com/stremovskyy/recorder v1.3.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.34 // indirect
	golang.org/x/text v0.41.0 // indirect
)

replace github.com/stremovskyy/go-ipay => ../..
