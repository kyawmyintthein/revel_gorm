package main

import(
	"os"
	"path"
	"strings"
	"log"
)

var cmdDBConfig = &Command{
	UsageLine: "db:config --driver=[driver_name]",
	Short:     "create new config for Revel application",
	Long: `
Create new database.conf file in app/conf of your application.

It puts all necessary files as import. 

--driver is required. The configuration will be changed based on --driver parameter.

For example:

    revel db:config --driver=[mysql]
`,
}

var driver flagValue

var dbConfigTplStr =  `[dev]
		driver = {{databaseDriver}}
		name = {{dbName}}
		user = root
		password = 
		maxConn = 30
		maxIdle = 30
		host = 127.0.0.1
		port = 3306`

func init() {
	cmdDBConfig.Run = dbConfig
	cmdDBConfig.Flag.Var(&driver, "driver", "database driver: mysql, mongodb, etc.")
}

func dbConfig(cmd *Command, args []string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("$GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}

	if len(args) == 0 {
		errorf("No driver given.\nRun 'revel help db' for usage.\n")
	}else{
		cmd.Flag.Parse(args[0:])
	}

	// Determine the run mode.
	mode := "dev"
	if len(args) >= 2 {
		mode = args[1]
	}

	pwd, _ := os.Getwd()
    destFile := path.Join(pwd, "conf", "database.conf")
    if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		if err = os.Remove(destFile); err != nil{
			log.Fatalln("Failed to remove existing file database.go:", err)
		}
	}
    if file, err := os.OpenFile(destFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer file.Close()
		content := strings.Replace(dbConfigTplStr, "[env]", mode, -1)
		content = strings.Replace(content, "{{databaseDriver}}", driver.String(), -1)
		content = strings.Replace(content, "{{dbName}}", "test", -1)
		file.WriteString(content)
	} else {
		errorf("No driver given.\nRun 'revel help db' for usage.\n")
	}
}
