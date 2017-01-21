package main

import(
	"os"
	"path"
	"strings"
	"log"
	"github.com/robfig/config"
)

var cmdDBSetup = &Command{
	UsageLine: "db:setup",
	Short:     "setup database based on database.conf for Revel application",
	Long: `
Based on database.conf config. It puts all necessary files as import. 

--driver is required. The configuration will be changed based on --driver parameter.

For example:
    revel db:setup
`,
}

var mysqlDatabaseTpl = `package {{packageName}}
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
    "github.com/robfig/config"
    "path"
    "errors"
)

var (
	DB gorm.DB
)

func InitDB() {
	var (
        err error
        driver string
        database string
        port string
        host string
        username string
        password string
    )

    //load config
    configPath := path.Join(revel.BasePath, "conf", "database.conf")
    config, err := config.ReadDefault(configPath)
    if err != nil || config == nil {
        panic(err)
    }

  
    switch revel.RunMode{
    case "dev":
        driver, _ =  config.String("dev", "driver");
        database, _ =  config.String("dev", "database");
        port, _ =  config.String("dev", "port");
        username, _ =  config.String("dev", "username");
        password, _ =  config.String("dev", "password");
        host, _ =  config.String("dev", "host");
    case "prod":
        driver, _ =  config.String("prod", "driver");
        database, _ =  config.String("prod", "database");
        port, _ =  config.String("prod", "port");
        username, _ =  config.String("prod", "username");
        password, _ =  config.String("prod", "password");
        host, _ =  config.String("prod", "host");
    default:
        panic(errors.New("Invalid RunMode"))
    }

    sqlConn := ""
    if port != "" && host != ""{
        sqlConn := username + ":" + password + "@tcp(" + host + ":" + port + ")/"+ database + "?charset=utf8"+ "&parseTime=True"
    }else{
        sqlConn := username + ":" + password + "@tcp/"+ database + "?charset=utf8"+ "&parseTime=True"
    }
    DB, err = gorm.Open(driver, sqlConn)
	if err != nil{
		panic(err.Error())
	}
	DB.LogMode(true)
}
`

var postgresqlDatabseTpl = `package {{packageName}}

import (  
 	"github.com/jinzhu/gorm"
 	_ "github.com/lib/pq"
 	"github.com/revel/revel"
    "github.com/robfig/config"
    "path"
    "errors"
)

var (
	DB gorm.DB
)

func InitDB() {  
	var (
        err error
        driver string
        database string
        port string
        host string
        username string
        password string
    )

    //load config
    configPath := path.Join(revel.BasePath, "conf", "database.conf")
    config, err := config.ReadDefault(configPath)
    if err != nil || config == nil {
        panic(err)
    }

  
    switch revel.RunMode{
    case "dev":
        driver, _ =  config.String("dev", "driver");
        database, _ =  config.String("dev", "database");
        port, _ =  config.String("dev", "port");
        username, _ =  config.String("dev", "username");
        password, _ =  config.String("dev", "password");
        host, _ =  config.String("dev", "host");
    case "prod":
        driver, _ =  config.String("prod", "driver");
        database, _ =  config.String("prod", "database");
        port, _ =  config.String("prod", "port");
        username, _ =  config.String("prod", "username");
        password, _ =  config.String("prod", "password");
        host, _ =  config.String("prod", "host");
    default:
        panic(errors.New("Invalid RunMode"))
    }

    DB, err :=  gorm.Open(driver, "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + database)
    if err != nil{
    	panic(err.Error())
    }
    db.LogMode(true)
}
`

var sqliteDatabseTpl = `package {{packageName}}

import (  
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/revel/revel"
    "github.com/robfig/config"
    "path"
    "errors"
)

var (
	DB gorm.DB
)

func InitDB(){  
    var (
        err error
        driver string
        dbname string
    )

    //load config
    configPath := path.Join(revel.BasePath, "conf", "database.conf")
    config, err := config.ReadDefault(configPath)
    if err != nil || config == nil {
        panic(err)
    }

    
    switch revel.RunMode{
    case "dev":
        driver, _ =  config.String("dev", "driver");
        dbname, _ =  config.String("dev", "database");
    case "prod":
        driver, _ =  config.String("prod", "driver");
        dbname, _ =  config.String("prod", "database");
    default:
        panic(errors.New("Invalid RunMode"))
    }

    DB, err := gorm.Open(driver, dbname)  
    if err != nil{
    	panic(err.Error())
    }
    DB.LogMode(true)
}`


func init() {
	cmdDBSetup.Run = dbSetup
}

func dbSetup(cmd *Command, args []string) {

	// get $GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ColorLog("[ERRO] $GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}

	// get curret folder
	pwd, _ := os.Getwd()

	// check config
	configPath := path.Join(pwd, "conf", "database.conf")
	config, err := config.ReadDefault(configPath)
	if err != nil || config == nil {
		ColorLog("[ERRO] database.conf not found in conf '%s'.\n", err)
		os.Exit(2)
	}
	ColorLog("[SUCC] '%s' is using as database config file.\n",configPath)

	databaseFolder := path.Join(pwd, "app", "models", "database")
	databaseFile := path.Join(databaseFolder, "database.go")

	//create database folder under models path
	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		os.MkdirAll(databaseFolder, 0777)
	}

	// generate files
	if _, err := os.Stat(databaseFile); !os.IsNotExist(err) {
		if err = os.Remove(databaseFile); err != nil{
			ColorLog("[ERRO] database.go is already exist. '%s'\n", err)
			os.Exit(2)
		}
	}
	if file, err := os.OpenFile(databaseFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer file.Close()
			driverType, _ :=  config.String("dev", "driver");
			ColorLog("[SUCC] '%s' is using as database driver.\n",driverType)
			switch driverType {
		    case "mysql":
		      	content := strings.Replace(mysqlDatabaseTpl, "{{packageName}}","database", -1)
				file.WriteString(content)
			case "sqlite3":
		      	content := strings.Replace(sqliteDatabseTpl, "{{packageName}}","database", -1)
				file.WriteString(content)
			case "postgresql":
		      	content := strings.Replace(postgresqlDatabseTpl, "{{packageName}}","database", -1)
				file.WriteString(content)
			default: 
				ColorLog("[ERRO] datbase driver not found.\nRun 'revel_gorm help' for usage.\n")
				os.Exit(2)
		    }
	} else {
		log.Println(err)
		ColorLog("[ERRO] Missing database.go.\nRun 'revel_gorm help' for usage.\n")
		os.Exit(2)
	}	

	ColorLog("[SUCC] database file is generated as '%s'.\n",databaseFile)
    
}
