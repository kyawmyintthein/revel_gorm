package main



import(
	"os"
	"path"
	"strings"
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

	body, err := GetTableBody(modelName, fields)
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

	viewFiles := []string{"Index", "Show", "New", "Edit", "form"}
	for _, filename := range viewFiles{
		currentFile := path.Join(crupath, "app", "views" , modelName, filename + ".html")
		if _, err := os.Stat(currentFile); os.IsNotExist(err) {
			if cf, err := os.OpenFile(currentFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
				defer cf.Close()
				switch filename{
				case "Index":
					content := strings.Replace(indexTpl, "{{pageTitle}}", modelName, -1)		
					content = strings.Replace(content, "{{pageHeader}}", modelName, -1)
					content = strings.Replace(content, "{{tableHeaders}}", header, -1)	
					content = strings.Replace(content, "{{body}}", body, -1)				
					cf.WriteString(content)
				case "Show":
					content := strings.Replace(indexTpl, "{{title}}", modelName, -1)		
					// content = strings.Replace(indexTpl, "{{header}}", modelName, -1)
					// content = strings.Replace(indexTpl, "{{tds}}", header, -1)	
					// content = strings.Replace(indexTpl, "{{body}}", body, -1)				
					cf.WriteString(content)
				case "New":
					content := strings.Replace(indexTpl, "{{title}}", modelName, -1)		
					// content = strings.Replace(indexTpl, "{{header}}", modelName, -1)
					// content = strings.Replace(indexTpl, "{{tds}}", header, -1)	
					// content = strings.Replace(indexTpl, "{{body}}", body, -1)				
					cf.WriteString(content)
				case "Edit":
					content := strings.Replace(indexTpl, "{{title}}", modelName, -1)		
					// content = strings.Replace(indexTpl, "{{header}}", modelName, -1)
					// content = strings.Replace(indexTpl, "{{tds}}", header, -1)	
					// content = strings.Replace(indexTpl, "{{body}}", body, -1)				
					cf.WriteString(content)
				case "form":
					content := strings.Replace(indexTpl, "{{title}}", modelName, -1)		
					// content = strings.Replace(indexTpl, "{{header}}", modelName, -1)
					// content = strings.Replace(indexTpl, "{{tds}}", header, -1)	
					// content = strings.Replace(indexTpl, "{{body}}", body, -1)				
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
{{template "footer.html" .}}
`

var showTpl = `
{{set . "title" "Home"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{header}}</h1>
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
		<table class="table table-striped">
			<thead>
				<th>
					{{tds}}
				</th>
			</thead>
			<tbody>
				{{trs}}
			</tbody>
		</table>
    </div>
  </div>
</div>
{{template "footer.html" .}}
`


var htmlModelTpl = `
{{set . "title" "Home"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{header}}</h1>
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
		<form method="{{method}}" action="{{action}}">
		</form>
    </div>
  </div>
</div>
{{template "footer.html" .}}
`
var formTpl = `
<form method="{{method}}" action="{{action}}">
		</form>
`
