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
)

var (
	DB gorm.DB
	Driver string
	Database string
	User  string
	Password  string
	Charset  string
)

func init() {
	var (
		sqlConn string
		err     error
	)
	sqlConn =  User + ":" + Password + "@/" + Database + "?charset=" + Charset+ "&parseTime=True"
	DB, err = gorm.Open("mysql", sqlConn)
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
)

var (
	DB gorm.DB
	Driver string
	Database string
	User  string
)

func init() {  
	var err error
    DB, err :=  gorm.Open(Driver, "user=" + User + " sslmode=disable"){
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
)

var (
	DB gorm.DB
	Driver string
	Database string
)

func init() {  
	var err error
    DB, err := gorm.Open(Driver, Database)  
    if err != nil{
    	panic(err.Error())
    }
    db.LogMode(true)
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
