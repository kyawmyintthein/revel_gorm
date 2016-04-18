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
	Short:     "create new config for Revel application",
	Long: `
Create new database.conf file in app/conf of your application.

It puts all necessary files as import. 

--driver is required. The configuration will be changed based on --driver parameter.

For example:
    revel db:setup
`,
}

var mysqlDBTplStr = `package database
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB gorm.DB
)

func init() {
	var (
		sqlConn string
		err     error
	)
	sqlConn =  "{{username}}:{{password}}@/{{dbName}}?charset=utf8&parseTime=True"
	if DB, err = gorm.Open("mysql", sqlConn); err != nil{
		panic(err.Error())
	}
	DB.LogMode(true)
}
`
func init() {
	cmdDBSetup.Run = dbSetup
}

func dbSetup(cmd *Command, args []string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("$GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}

	// Determine the run mode.
	mode := "dev"
	if len(args) >= 2 {
		mode = args[1]
	}

	pwd, _ := os.Getwd()

	config, err := config.ReadDefault(path.Join(pwd, "conf", "database.conf"))
	if err != nil || config == nil {
		log.Fatalln("Failed to load database.conf:", err)
	}

	databaseFolder := path.Join(pwd, "app", "models", "database")
	databaseFile := path.Join(databaseFolder, "database.go")
	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.MkdirAll(databaseFolder, 0777)
	}

	if _, err := os.Stat(databaseFile); !os.IsNotExist(err) {
		if err = os.Remove(databaseFile); err != nil{
			log.Fatalln("Failed to remove existing file database.go:", err)
		}
	}
	if file, err := os.OpenFile(databaseFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer file.Close()
			user,_ :=  config.String(mode, "user");
			password, _ :=  config.String(mode, "password");
			dbName, _ :=  config.String(mode, "name");
			driverType, _ :=  config.String(mode, "driver");

			switch driverType {
		    case "mysql":
		      	content := strings.Replace(mysqlDBTplStr, "{{username}}",user, -1)
				content = strings.Replace(content, "{{password}}",password, -1)
				content = strings.Replace(content, "{{dbName}}", dbName, -1)
				file.WriteString(content)
			default: 
				errorf("datbase driver not found.\nRun 'revel help db' for usage.\n")
		    }
	} else {
		log.Println(err)
		errorf("Missing database.go.\nRun 'revel help db' for usage.\n")
	}	
    
}
