Introduction
------------

This project combines the Gin web framework, Gorm ORM, Go-Redis and Viper for building a web application. 
It provides a simple example of how to use martini to create a web application.

Installation
------------

### Install Martini

To install Gin, run the following command in your terminal:

```bash
go get -u github.com/healer1219/martini
```

Usage
-----

### Create Configuration File

Create a new file called `config.yaml` and add the following code:
The "app" and "log" are mandatory items in the configuration file.

```yaml

app:
  env: dev
  port: 8090
  name: martini

log:
  level: info # 日志等级
  root_dir: logs # 日志根目录
  filename: app.log # 日志文件名称
  format: json # 写入格式 可选json
  show_line: true # 是否显示调用行
  max_backups: 3 # 旧文件的最大个数
  max_size: 500 # 日志文件最大大小（MB）
  max_age: 28 # 旧文件的最大保留天数
  compress: true # 是否压缩

cloud:
  ip: 127.0.0.1
  port: 8500

database:
  database_type: mysql
  ip: 127.0.0.1
  port: 3306
  database_name: DB_NAME
  username: root
  password: root
  charset: utf8mb4
  max_idle_conns: 10 # 空闲连接池中连接的最大数量
  max_open_conns: 100 # 打开数据库连接的最大数量
  log_mode: info # 日志级别
  enable_file_log_writer: true # 是否启用日志文件
  log_filename: sql.log # 日志文件名称

dbs:
  default:
    database_type: mysql
    ip: 127.0.0.1
    port: 3306
    database_name: DB_NAME
    username: root
    password: root
    charset: utf8mb4
    max_idle_conns: 10 # 空闲连接池中连接的最大数量
    max_open_conns: 100 # 打开数据库连接的最大数量
    log_mode: info # 日志级别
    enable_file_log_writer: true # 是否启用日志文件
    log_filename: sql.log # 日志文件名称
  custom:
    database_type: mysql
    ip: 127.0.0.1
    port: 3306
    database_name: DB_NAME
    username: root
    password: root
    charset: utf8mb4
    max_idle_conns: 10 # 空闲连接池中连接的最大数量
    max_open_conns: 100 # 打开数据库连接的最大数量
    log_mode: info # 日志级别
    enable_file_log_writer: true # 是否启用日志文件
    log_filename: sql.log # 日志文件名称

custom:
  custom_key_1: oohhhh
  customMap:
    mapKey1: 1
    mapKey2: 2

```

### Create Router

Create a method that returns []bootstrap.RouteOption

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/healer1219/martini/bootstrap"
)

func Init() []bootstrap.RouteOption {
	options := []bootstrap.RouteOption{
		Routers,
	}
	return options
}

func Routers(engine *gin.Engine) {
	group := engine.Group("/awesome-app/v1")
	{
		group.GET("/user", queryUser)
		group.PUT("/user", updateUser)
		group.POST("/user", addUser)
	}
}
```
### Start your application

```go
package main

import (
	apis "awesome-app/api"
	"firewalld-manage/auth"
	"github.com/gin-gonic/gin"
	"github.com/healer1219/martini/bootstrap"
)

func main() {

	bootstrap.Default().
		Router(apis.Init()...).
		StartFunc(
			auth.InitAuth,  // start hook
		).
		ShutDownFunc(
			auth.ShutDown,  // shut down hook
		).
		Use(
			auth.JwtHandler,  // gin middleware
		).BootUp()
}

```

### use Gorm

make sure the config.yaml `database` or `dbs` is not empty

```go
package db

import (
	"github.com/healer1219/martini/global"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string
	Age       int
	Email     string `gorm:"uniqueIndex"`
}

func AutoMigrate() error {
	err = global.DB().AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

func SaveUser(user *User) error {
	if user == nil {
		user = User{Name: "John Doe", Age: 30, Email: "john@example.com"}
	}
	
	err = global.DB().Create(&user)
	if err != nil {
		return err
	}
	return nil
}

func FindUser(id uint) (*User, error) {
	var user = User{}
	err = global.DB().First(&user, user.ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	err = global.DB().Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(user *User) error {
	err = global.DB().Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

```

