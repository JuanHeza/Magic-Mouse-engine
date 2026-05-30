package data

import (
	"database/sql"
	"errors"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const databaseErrorOccurred = "A database error occurred %v"
var (
	//ErrDBConfig is returned when db configuration is invalid.
	ErrDBConfig = errors.New("invalid database configuration")
)


// Database Database holder
type Database struct {
	Db     *sql.DB
	Dbx *sqlx.DB
	logger *logrus.Entry
}

// NewDatabase returns new Database
func NewDatabase(log *logrus.Entry) *Database {
	return &Database{logger: log}
}

// genConnectionString replace host and port params
func (mdb *Database) genConnectionString(config *Config) string {
	var cs string
	cs = config.DB.ConnectionString

	//connectionString: "host=$ port=$ user=$ password=$ dbname=d_pos_gotemplate sslmode=disable search_path=s_pos_gotemplate"
	if config.DB.Server != "" {
		cs = strings.ReplaceAll(cs, "host=$", "host="+config.DB.Server)
	}

	if config.DB.Port != 0 {
		cs = strings.ReplaceAll(cs, "port=$", "port="+strconv.Itoa(config.DB.Port))
	}

	if config.DB.User != "" {
		cs = strings.ReplaceAll(cs, "user=$", "user="+config.DB.User)
	}

	if config.DB.Password != "" {
		cs = strings.ReplaceAll(cs, "password=$", "password="+config.DB.Password)
	}

	return cs
}

// Init init database
func (mdb *Database) Init(config *Config) error {
	var err error
	log := mdb.logger
	connLimitLog := config.DB.InitConnectLimitLog

	if config.DB.InitConnectRetries <= 0 {
		log.WithFields(logrus.Fields{"event": "DBConnectionError"}).
			Errorf("Failed to connect Database: invalid InitConnectRetries: %v", config.DB.InitConnectRetries)
		return ErrDBConfig
	}
	if (config.DB.InitConnectRetries < config.DB.InitConnectLimitLog) || (config.DB.InitConnectLimitLog == 0) {
		connLimitLog = 1
	}

	for connRetryCount := 0; connRetryCount < config.DB.InitConnectRetries; connRetryCount++ {
		err = mdb.Connect(config)
		if err == nil {
			log.WithFields(logrus.Fields{"event": "DBConnectionSuccess"}).
				Info("Database connected")
			return nil
		}

		isDBConnRetryErr := false
		var dbErrCode string
		switch dbErr := err.(type) {
		case *pgconn.PgError:
			for _, errorCode := range config.DB.SQLReconnectOnErrors {
				if dbErr.Code == errorCode {
					isDBConnRetryErr = true
				}
			}
			dbErrCode = "PgxError:" + dbErr.Code

		case *net.OpError:
			switch opErr := dbErr.Err.(type) {
			case *os.SyscallError:
				if opErr.Err == syscall.ECONNREFUSED {
					isDBConnRetryErr = true
				}
				dbErrCode = "SyscallError:" + opErr.Error()

			case *net.DNSError:
				isDBConnRetryErr = true
				dbErrCode = "DNSError:" + opErr.Error()

			default:
				dbErrCode = "OPError:" + opErr.Error()
			}
		default:
			dbErrCode = reflect.TypeOf(err).String() + ":" + err.Error()
		}

		//DB non-retriable connection errors
		if !isDBConnRetryErr {
			log.WithFields(logrus.Fields{"event": "DBConnectionError", "errorCode": dbErrCode}).
				Error("Failed to connect Database: ", err.Error())
			break
		}
		//DB retriable connection errors
		if (connRetryCount % connLimitLog) == 0 {
			log.WithFields(logrus.Fields{"event": "DBConnectionError", "errorCode": dbErrCode}).
				Error("Failed to connect Database: ", err.Error())
		}
		time.Sleep(config.DB.InitConnectSleep)
	}
	return err
}

// Connect Connect database
func (mdb *Database) Connect(config *Config) error {
	var err error

	mdb.Db, err = sql.Open("pgx", mdb.genConnectionString(config))
	if err != nil {
		return err
	}

	db, err := sqlx.Connect("postgres", mdb.genConnectionString(config))

	mdb.Dbx = db
	mdb.Db.SetMaxIdleConns(1)
	mdb.Db.SetConnMaxIdleTime(0)
	mdb.Db.SetMaxOpenConns(10)
	//Check if database connection is established
	err = mdb.Db.Ping()

	if err != nil {
		return err
	}
	
	return nil
}

// Close close db
func (mdb *Database) Close(config *Config) {
	log := mdb.logger

	log.Error("Database connection closed")
	err := mdb.Db.Close()

	if err != nil {
		log.Error(" Failed db.Close(): ", err.Error())
	}
}
