package main

import(
	"os"
)

var cmdGenerate = &Command{
	UsageLine: "generate [Command]",
	Short:     "source code generator",
	Long: `
revel_grom generate scaffold [scaffoldname] [-fields=""]
    example: revel_grom generate scaffold post -fields="title:string,body:text"

revel_grom generate model [modelname] [-fields=""]
    example: revel_grom generate model post -fields="title:string,body:text"

revel_grom generate controller [controllerfile]
    generate RESTFul controllers 
    example: revel_grom generate controller post
`,
}

var fields flagValue

func init() {
	cmdGenerate.Run = generateCode
	cmdGenerate.Flag.Var(&fields, "fields", "specify the fields want to generate.")
}

func generateCode(cmd *Command, args []string){
	// get current path
	curpath, _ := os.Getwd()

	// check args
	if len(args) < 1 {
		ColorLog("[ERRO] command is missing\n")
		os.Exit(2)
	}

	// check $GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ColorLog("[ERRO] $GOPATH not found\n")
		ColorLog("[HINT] Set $GOPATH in your environment vairables\n")
		os.Exit(2)
	}

	gcmd := args[0]
	switch gcmd {
	case "api-scaffold":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate scaffold [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate scaffold [modelname] [-fields=\"title:string,body:text\"]\n")
			os.Exit(2)
		}
		sname := args[1]
		ColorLog("[INFO] Using '%s' as controller name\n", sname)
		ColorLog("[INFO] Using '%s' as controller name\n", sname + "Controller")

		//generate model and controller
		generateModel(sname, fields.String(), curpath)
		generateRestController(sname, curpath)
	case "api-controller":
		if len(args) == 2 {
			cname := args[1]
			generateRestController(cname, curpath)
		} else {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate controller [controllername]\n")
			os.Exit(2)
		}
	case "scaffold":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate scaffold [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: bee generate scaffold [modelname] [-fields=\"title:string,body:text\"]\n")
			os.Exit(2)
		}
		sname := args[1]
		ColorLog("[INFO] Using '%s' as controller name\n", sname)
		ColorLog("[INFO] Using '%s' as controller name\n", sname + "Controller")

		//generate model and controller
		generateModel(sname, fields.String(), curpath)
		generateController(sname, curpath)
	case "controller":
		if len(args) == 2 {
			cname := args[1]
			generateController(cname, curpath)
		} else {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate controller [controllername]\n")
			os.Exit(2)
		}
	case "model":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate model [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate model [modelname] [-fields=title:string,body:string]\n")
			os.Exit(2)
		}
		sname := args[1]
		generateModel(sname, fields.String(), curpath)
	case "views":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate model [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo generate model [modelname] [-fields=title:string,body:string]\n")
			os.Exit(2)
		}
		sname := args[1]
		generateViews(sname, fields.String(), curpath)
	default:
		ColorLog("[ERRO] command is missing\n")
		os.Exit(2)
	}
	ColorLog("[SUCC] generate successfully created!\n")
	
}
