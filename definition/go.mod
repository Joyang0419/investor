module definition

go 1.22.2

require (
	go.mongodb.org/mongo-driver v1.14.0
	gorm.io/gorm v1.25.8
	tools v0.0.0
)

replace tools v0.0.0 => ../tools

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
