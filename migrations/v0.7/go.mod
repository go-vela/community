module github.com/go-vela/community/migrations/v0.7

go 1.15

replace github.com/go-vela/server => ../../../server

require (
	github.com/go-vela/server v0.7.3
	github.com/go-vela/types v0.7.4-0.20210225205732-6bf075d597f6
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.3.0
	github.com/sirupsen/logrus v1.8.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)
