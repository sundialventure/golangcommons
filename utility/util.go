package utility

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	//"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/lib/pq"
	"github.com/sundialventure/golangcommons/cfg"
	"github.com/sundialventure/golangcommons/sys"
	"github.com/xo/dburl"
	"golang.org/x/crypto/bcrypt"
	//log15 "gopkg.in/inconshreveable/log15.v2"
)

//DB ...
var DB gorm.DB
var DBSql *sql.DB

var conn *gorm.DB
var oracleConn *sql.DB
var dberr error

//DataStoreObject ...
var DataStoreObject *DataStore

//Keebler ...
var Keebler *securecookie.SecureCookie

var newLogger logger.Interface


//SetupDBNoEntities ... Setup DB without any entity
func SetupDBNoEntities(config map[string]interface{}, dbtype string) {
	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logger.Silent, // Log level
		  IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		  Colorful:                  false,          // Disable color
		},
	  )
	
	ctx := context.Background()
	

	if dbtype == "mssql" {
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", cfg.DB_SERVER, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_PORT, cfg.DB_NAME)
		//conn, dberr = gorm.Open("mssql", connString) //sqlx.Open("mssql", connString)
		conn, dberr = gorm.Open(sqlserver.New(sqlserver.Config{
			DriverName: "my_sqlserver_driver",
			DSN:        connString, //"gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		}),&gorm.Config{
			Logger: newLogger,
		  })
	}
	if dbtype == "postgres" {
		connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME) //host=%s userid=%s port=%d password=%s dbname=%s
		
		conn, dberr = gorm.Open(postgres.New(postgres.Config{
			DriverName: "my_postgres_driver",
			DSN:        connString, //"gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		}), &gorm.Config{
			Logger: newLogger,
		  })
	}

	if dbtype == "oracle" {
		connString := fmt.Sprintf("%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=tcp)(HOST=%s)(PORT=%d)))(CONNECT_DATA=(SERVICE_NAME=%s)))", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME)
		oracleConn, dberr = sql.Open("goracle", connString)
		oracleConn.SetMaxOpenConns(20)
		oracleConn.SetMaxIdleConns(10)
	}

	DataStoreObject = &DataStore{}
	DataStoreObject.StoreType = "rdbms"
	if dbtype != "oracle" {
		
		DB = *conn

		if dberr != nil {
			newLogger.Error(ctx, "Error on opening database connection: %s", dberr.Error())
			panic(dberr)
		} else {
			newLogger.Info(ctx, "DB has been pinged")
		}

	} else {
		DBSql = oracleConn
		dberr := DBSql.Ping()
		if dberr != nil {
			newLogger.Error(ctx, "Error on opening database connection: %s", dberr.Error())
			panic(dberr)
		} else {
			newLogger.Info(ctx, "DB has been pinged")
		}

	}

	

	DataStoreObject.RDBMS = RDBMSImpl{conn, ctx}

	//PanicIf(dberr)

}

//SetupDB ...
func SetupDB(ctx context.Context, config map[string]interface{}, entites []interface{}, dbtype string) {

	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logger.Silent, // Log level
		  IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		  Colorful:                  false,          // Disable color
		},
	  )

	if dbtype == "mssql" {
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", cfg.DB_SERVER, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_PORT, cfg.DB_NAME)
		//conn, dberr = gorm.Open("mssql", connString) //sqlx.Open("mssql", connString)
		conn, dberr = gorm.Open(sqlserver.New(sqlserver.Config{
			DriverName: "my_sqlserver_driver",
			DSN:        connString, //"gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		}),&gorm.Config{
			Logger: newLogger,
		  })
	}
	if dbtype == "postgres" {
		connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME) //host=%s userid=%s port=%d password=%s dbname=%s
		
		conn, dberr = gorm.Open(postgres.New(postgres.Config{
			DriverName: "my_postgres_driver",
			DSN:        connString, //"gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		}), &gorm.Config{
			Logger: newLogger,
		  })
	}

	if cfg.DB_TYPE == "oracle" {
		connString := fmt.Sprintf("%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=tcp)(HOST=%s)(PORT=%d)))(CONNECT_DATA=(SERVICE_NAME=%s)))", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME)
		oracleConn, dberr = sql.Open("goracle", connString)
	}

	DataStoreObject = &DataStore{}
	DataStoreObject.StoreType = "rdbms"
	//conn.LogMode(true)
	//DB.LogMode(true)
	DB = *conn
	DataStoreObject.RDBMS = RDBMSImpl{conn, ctx}

	DataStoreObject.InitSchema(entites)
	//PanicIf(dberr)

	if dberr != nil {
		log.Fatalf("Error on opening database connection: %s", dberr.Error())
	}
	log.Printf("DB has been pinged")
}

//SetupDBReturn ...
func SetupDBReturn(config map[string]interface{}, entites []interface{}, dbtype string) (db *gorm.DB, er error) {
	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	ctx := context.Background()
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logger.Silent, // Log level
		  IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		  Colorful:                  false,          // Disable color
		},
	  )
	/*settings := "user=" + cfg.DB_USER + " password=" + cfg.DB_PASSWORD + " dbname=" + cfg.DB_NAME + " sslmode=disable"
	db, err := sql.Open("postgres", settings)

	PanicIf(err)
	DB = db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")*/
	connString := fmt.Sprintf("host=%s user=%s port=%d password=%s dbname=%s sslmode=disable", cfg.DB_SERVER, cfg.DB_USER, cfg.DB_PORT, cfg.DB_PASSWORD, cfg.DB_NAME)
	fmt.Println(connString)
	conn, dberr = gorm.Open(postgres.New(postgres.Config{
		DriverName: "my_postgres_driver",
		DSN:        connString, //"gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
	}), &gorm.Config{
		Logger: newLogger,
	  })
	DataStoreObject = &DataStore{}
	DataStoreObject.StoreType = dbtype
	
	DB = *db
	DataStoreObject.RDBMS = RDBMSImpl{db, ctx}

	DataStoreObject.InitSchema(entites)
	PanicIf(dberr)

	if dberr != nil {
		log.Fatalf("Error on opening database connection: %s", dberr.Error())
		return nil, dberr
	}
	log.Printf("DB has been pinged")
	return db, nil
}

//SetupDB ...
func SetupDB2(config map[string]interface{}) {
	cfg.DB_SERVER = config["address"].(string)
	cfg.DB_USER = config["userid"].(string)
	cfg.DB_TYPE = config["dbtype"].(string)
	cfg.DB_NAME = config["dbname"].(string)
	cfg.DB_PASSWORD = config["password"].(string)
	cfg.DB_PORT = int64(config["port"].(float64))
	/*settings := "user=" + cfg.DB_USER + " password=" + cfg.DB_PASSWORD + " dbname=" + cfg.DB_NAME + " sslmode=disable"
	db, err := sql.Open("postgres", settings)

	PanicIf(err)
	DB = db
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	log.Printf("DB has been pinged")*/
	connString := fmt.Sprintf("pgsql://u%s:%s@%s:%d/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_SERVER, cfg.DB_PORT, cfg.DB_NAME)
	db, err := dburl.Open(connString)
	PanicIf(err)
	DBSql = db
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

	log.Printf("DB has been pinged")
	//DBSql.LogMode(true)
}

//SetupKeebler ...
func SetupKeebler() {
	Keebler = securecookie.New([]byte(cfg.HASHKEY), []byte(cfg.BLOCKKEY))
	log.Printf("Keebler has been created.... ready to make cookies...")
}

//GenHash ...
func GenHash(c int) string {
	//c := 8
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

//CookieHasAccess ...
func CookieHasAccess(c *gin.Context) (bool, string, string) {
	if cookie, err := c.Request.Cookie(cfg.COOKIE_NAME); err == nil {
		value := make(map[string]string)
		if err = Keebler.Decode(cfg.COOKIE_NAME, cookie.Value, &value); err == nil {
			log.Printf("The value of access is " + value[sys.COOKIE_APP_ACCESS])
			log.Printf("The value of email is " + value[sys.COOKIE_EMAIL])
			log.Printf("Access value of [COOKIE_APP_ACCESS]: %v", value[sys.COOKIE_APP_ACCESS])
			if value[sys.COOKIE_APP_ACCESS] == sys.ACCESS_OK {
				return true, value[sys.COOKIE_EMAIL], value[sys.COOKIE_USERID]
			}
		}
	}
	log.Printf("Where is the cookie? %v", cfg.COOKIE_NAME)
	return false, "", ""
}

//PanicIf ...
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// Schedule is a generic timer func that will run a func after a delay
func Schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
	// stop is the channel that will stop if you do: stop <- true
}

//HashPassword ...
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash, err
}

//Round ...
func Round(val float64, prec int) float64 {
	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}

// Btc_Sat converts BTC float to Satochi Int; 100 million sat = 1BTC
// int64 value +/- 9223372036854775807
//Btc_Sat ...
func Btc_Sat(btc float64) (sat int64, err error) {
	num := btc * 100000000
	if (num > 9223372036854775807) || (num < -9223372036854775807) {
		num = -1
		err = errors.New("Out of scale for int64")
	}
	return int64(num), err
}

// TODO write check for float64 size
func Sat_Btc(sat int64) (btc float64, err error) {
	num := float64(sat)
	if (num > 9223372036854775807) || (num < -9223372036854775807) {
		num = -1
		err = errors.New("Out of scale for float64")
	}
	btc = num / 100000000
	return btc, err
}

//RandStr ...
func RandStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	if randType == "all" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	if randType == "number" {
		dictionary = "0123456789"
	}
	if randType == "all" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

//GenerateTransRef ... Generate transaction reference for each country
func GenerateTransRef(countryid int) string {
	time_now := time.Now()
	transRef := fmt.Sprintf(time_now.Format("20060102150405"))
	if countryid == 1 {
		transRef = fmt.Sprintf(time_now.Format("20060102150405"))
	}
	return transRef
}
