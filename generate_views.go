package main



import(
	"os"
	"path"
	"strings"
	 "github.com/yosssi/gohtml"
)

// generate controller
func generateViews(mname, fields, crupath string) {
	// get controller name and package 
	_, f := path.Split(mname)

	// set controller name to uppercase
	modelName := strings.Title(f)

	header, err := GetTableHeaders(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate views : %s\n", err)
		os.Exit(2)
	}

	body, err := GetIndexTableBody(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate views: %s\n", err)
		os.Exit(2)
	}


	showTable, err := GetShowTableBody(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate views: %s\n", err)
		os.Exit(2)
	}

	form, err := GetFormAttributes(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate views: %s\n", err)
		os.Exit(2)
	}

	updateForm, err := GetUpdateFormAttributes(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate views: %s\n", err)
		os.Exit(2)
	}


	ColorLog("[INFO] Using '%s' is generated in views path\n", modelName)

	// create controller folders
	filePath := path.Join(crupath ,"app", "views", modelName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(filePath, 0777); err != nil {
			ColorLog("[ERRO] Could not create views directory: %s\n", err)
			os.Exit(2)
		}
	}

	viewFiles := []string{"Index", "Show", "New", "Edit"}
	for _, filename := range viewFiles{
		currentFile := path.Join(crupath, "app", "views" , modelName, filename + ".html")
		if _, err := os.Stat(currentFile); os.IsNotExist(err) {
			if cf, err := os.OpenFile(currentFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
				defer cf.Close()
				switch filename{
				case "Index":
					content := strings.Replace(indexTpl, "{{pageTitle}}", modelName, -1)		
					content = strings.Replace(content, "{{pageHeader}}", modelName + "s", -1)
					content = strings.Replace(content, "{{tableHeaders}}", header, -1)	
					content = strings.Replace(content, "{{body}}", body, -1)	
					content = strings.Replace(content, "{{newMethod}}", "New", -1)	
					content = strings.Replace(content, "{{modelName}}", modelName, -1)		
					gohtml.FormatWithLineNo(content)			
					cf.WriteString(content)
				case "Show":
					content := strings.Replace(showTpl, "{{pageTitle}}", modelName, -1)		
					content = strings.Replace(content, "{{pageHeader}}", modelName, -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)
					content = strings.Replace(content, "{{body}}", showTable, -1)	
					content = strings.Replace(content, "{{editMethod}}", "Edit", -1)	
					content = strings.Replace(content, "{{indexMethod}}", "Index", -1)		
					gohtml.FormatWithLineNo(content)				
					cf.WriteString(content)
				case "New":
					content := strings.Replace(newTpl, "{{pageTitle}}", modelName, -1)		
					content = strings.Replace(content, "{{pageHeader}}", "New " + modelName, -1)	
					content = strings.Replace(content, "{{method}}", "POST", -1)	
					content = strings.Replace(content, "{{action}}", strings.ToLower(modelName) +"s", -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)	
					content = strings.Replace(content, "{{buttonName}}", "Create " + strings.ToLower(modelName), -1)
					content = strings.Replace(content, "{{formAttributes}}", form, -1)	
					content = strings.Replace(content, "{{Action}}", "Create" , -1)			
					gohtml.FormatWithLineNo(content)				
					cf.WriteString(content)
				case "Edit":
					content := strings.Replace(editTpl, "{{pageTitle}}", modelName, -1)		
					content = strings.Replace(content, "{{pageHeader}}", "Edit " + modelName, -1)
					content = strings.Replace(content, "{{method}}", "POST", -1)	
					content = strings.Replace(content, "{{action}}", strings.ToLower(modelName) + "s/update", -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)	
					content = strings.Replace(content, "{{buttonName}}", "Save " + strings.ToLower(modelName), -1)	
					content = strings.Replace(content, "{{formAttributes}}", updateForm, -1)	
					content = strings.Replace(content, "{{Action}}", "Update" , -1)		
					gohtml.FormatWithLineNo(content)				
					cf.WriteString(content)
				}	
				// gofmt generated source code
				// FormatSourceCode(currentFile)

				ColorLog("[INFO] '%s' view file are generated as: %s\n", filename,currentFile)
			} else {
				// error creating file
				ColorLog("[ERRO] Could not create views for '%s': %s\n",modelName,err)
				os.Exit(2)
			}
		}

	} 
}


// remove existing vuews file
func deleteViews(mname, crupath string) {
	_, f := path.Split(mname)
	modelName := strings.Title(f)
	filePath := path.Join(crupath, "app", "views", modelName)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err = os.Remove(filePath)
		if err != nil{
			ColorLog("[ERRO] Could not delete views: %s\n", err)
			os.Exit(2)	
		}
		ColorLog("[INFO] views files are deleted: %s\n", filePath)	
	}

}


var indexTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
        <div class="panel-header">
	  		<a href="{{url "{{modelName}}.{{newMethod}}"}}" class="btn btn-default">New</a>
	    </div>
    	<div class="panel panel-default">
		  <div class="panel-body">
			<table class="table table-striped">
				<thead>
					<tr>
						{{tableHeaders}}
					</tr>
				</thead>
				<tbody>
					{{body}}
				</tbody>
			</table>
          </div>
		</div>
    </div>
  </div>
</div>
{{template "footer.html" .}}
`

var showTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
    	<div class="panel panel-default">
		  <div class="panel-body">
		    <table class="table table-striped">
				<tbody>
					{{body}}
				</tbody>
			</table>
		  </div>
		  <div class="panel-footer">
		  		<a href="{{url "{{modelName}}.{{editMethod}}"}}" class="btn btn-info">Edit</a>
		  		<a href="{{url "{{modelName}}.{{indexMethod}}"}}" class="btn btn-default">{{modelName}}s</a>
    	  </div>
		</div>
    </div>
  </div>
</div>
{{template "footer.html" .}}
`


var editTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
    	<form method="POST" action="{{ url "{{modelName}}.{{Action}}"}}">
			{{formAttributes}}
    		<input type="submit" class="btn btn-success" value="{{buttonName}}" />
		</form>
    </div>
  </div>
</div>
{{template "footer.html" .}}`

var newTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
    	<form method="POST" action="{{ url "{{modelName}}.{{Action}}"}}">
			{{formAttributes}}
    		<input type="submit" class="btn btn-success" value="{{buttonName}}" />
		</form>
    </div>
  </div>
</div>
{{template "footer.html" .}}`
