package main

import(
	"os"
	"path"
	"strings"
)

var cmdDBConfig = &Command{
	UsageLine: "db:config -driver=[driver_name] -dbname=[dbname]",
	Short:     "generate gorm database config for Revel application",
	Long: `
Create new database.conf file in app/conf of your application.
It puts all necessary files as import. 
--driver is required. The configuration will be changed based on --driver parameter.
For example:
    revel db:config -driver=[mysql] -dbName=[dbname]
`,
}

var driver flagValue
var dbName flagValue


func init() {
	cmdDBConfig.Run = generateDBConfig
	cmdDBConfig.Flag.Var(&driver, "driver", "database driver: mysql, sqlite3, postgresql, etc.")
	cmdDBConfig.Flag.Var(&dbName, "dbname", "database name: test, app_name, etc.")
}

func generateDBConfig(cmd *Command, args []string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ColorLog("[ERRO] $GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}

	// check driver params
	if len(args) < 0 {
		ColorLog("[ERRO] No driver given.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}else{
		cmd.Flag.Parse(args[0:])
	}

	// check driver params
	if len(args) == 1 {
		cmd.Flag.Parse(args[1:])
	}

	if dbName.String() == ""{
		dbName = "test"
	}

	ColorLog("[SUCC] database driver is using as '%s'.\n",driver.String())
	ColorLog("[SUCC] database name is using as '%s'.\n",dbName.String())

	// get current path
	pwd, _ := os.Getwd()
    destFile := path.Join(pwd, "conf", "database.conf")
    if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		if err = os.Remove(destFile); err != nil{
			ColorLog("[WARN] Failed to remove existing file database.go:", err)
			os.Exit(2)
		}
	}

    if file, err := os.OpenFile(destFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer file.Close()

		switch driver{
			case "mysql":
				content := strings.Replace(mysqlDatabaseConfigTpl, "{{dbName}}", dbName.String(), -1)
				file.WriteString(content)
			case "sqlite3":
				content := strings.Replace(sqliteDatabaseConfigTpl, "{{dbName}}", dbName.String(), -1)
				file.WriteString(content)
			case "postgresql":
				content := strings.Replace(postgresDatabaseConfigTpl, "{{dbName}}", dbName.String(), -1)
				file.WriteString(content)
			default: 
			  ColorLog("[ERRO] Invalid database driver.\n");
			  os.Exit(2) 
		}
		
	} else {
		ColorLog("[ERRO] No driver given.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}
	ColorLog("[SUCC] grom database config file is generated as '%s'.\n",destFile)

}


// templates
var mysqlDatabaseConfigTpl =  
`[dev]
driver: mysql
encoding: utf8
database: {{dbName}}_dev
maxConn = 30
maxIdle = 30
username: root
password:
host = 127.0.0.1
port = 3306


[prod]
driver: mysql
encoding: utf8
database: {{dbName}}
maxConn = 30
maxIdle = 30
username: root
password:
host = 127.0.0.1
port = 3306


[test]
driver: mysql
encoding: utf8
database: {{dbName}}_test	
maxConn = 30
maxIdle = 30
username: root
password:
host = 127.0.0.1
port = 3306`
		
var sqliteDatabaseConfigTpl =  
`[dev]
driver: sqlite3
database: db/{{dbName}}_dev.sqlite3
maxConn = 30
maxIdle = 30
timeout: 5000


[prod]
driver: sqlite3
database: db/{{dbName}}.sqlite3
maxConn = 30
maxIdle = 30
timeout: 5000


[test]
driver: sqlite3
database: db/{{dbName}}_test.sqlite3
maxConn = 30
maxIdle = 30
timeout: 5000`
		
var postgresDatabaseConfigTpl =  
`[dev]
driver: postgres
encoding: unicode
database: {{dbName}}_dev
maxConn = 30
maxIdle = 30
username: root
password:


[prod]
driver: postgres
encoding: unicode
database: {{dbName}}
maxConn = 30
maxIdle = 30
username: root
password:


[test]
driver: postgres
encoding: unicode
database: {{dbName}}_test
maxConn = 30
maxIdle = 30
username: root
password:`

