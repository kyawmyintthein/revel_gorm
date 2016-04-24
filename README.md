# **revel_gorm**

## About
  **revel_gorm** is a code generator for **Revel web framework** https://revel.github.io/ with **gorm** https://github.com/jinzhu/gorm ORM library.
  It can also generate **RESTFul API**. 
  revel_gorm can setup database using mysql, postgresql, sqlite3 in your revel project easily and it is using "gorm" for ORM.
  It can also generate models, controllers and views for your revel project and also for restful controller.

## Installation
    go get github.com/kyawmyintthein/revel_gorm  

## Usage
#### Setup database config
    revel_gorm db:config -driver=[mysql,postgres,sqlite]

#### Setup database 
    revel_gorm db:setup 

#### Setup database 
    revel_gorm db:setup 

#### Update the database.conf for your application.
  *Add following code in your database.conf file and change your database configuration.*
   
      [dev]
	driver: mysql
	encoding: utf8
	database: example_dev
	maxConn = 30
	maxIdle = 30
	username: root
	password:
	host = 127.0.0.1
	port = 3306
	
	
	[prod]
	driver: mysql
	encoding: utf8
	database: example_dev
	maxConn = 30
	maxIdle = 30
	username: root
	password:
	host = 127.0.0.1
	port = 3306
	
	
	[test]
	driver: mysql
	encoding: utf8
	database: example_dev	
	maxConn = 30
	maxIdle = 30
	username: root
	password:
	host = 127.0.0.1
	port = 3306
	  

*Add this code under your init function of init.go file*

    import "project/app/models/database"
    revel.OnAppStart(database.InitDB)


#### Generate model 
    revel_gorm generate model ModelName -fields="fieldname:string,fieldname:int,fieldname:bool,,datetime"

#### Generate controller 
    revel_gorm generate controller ModelName


#### Generate rest-controller  for json response
    revel_gorm generate rest-controller ModelName


#### Generate views 
    revel_gorm generate views ModelName -fields="fieldname:string,fieldname:int,fieldname:bool,,datetime"

#### Scaffold model + views + controller 
    revel_gorm generate scaffold ModelName -fields="fieldname:string,fieldname:int,fieldname:bool,,datetime"
    
#### Scaffold model + views + controller  for API
    revel_gorm generate res-scaffold ModelName -fields="fieldname:string,fieldname:int,fieldname:bool,,datetime"
