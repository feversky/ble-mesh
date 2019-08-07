module ble-mesh

go 1.12

replace github.com/bettercap/gatt => ./driver/gatt

require (
	github.com/ChengjinWu/aescrypto v0.0.0-20181106090456-e3bda2891c3d
	github.com/aead/cmac v0.0.0-20160719120800-7af84192f0b1
	github.com/aead/ecdh v0.0.0-20190219101236-85c03e91d99a
	github.com/bettercap/gatt v0.0.0-20190418085356-fac16c0ad797
	github.com/c-bata/go-prompt v0.2.3
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/stretchr/testify v1.2.2
	github.com/thoas/go-funk v0.4.0
	github.com/ugorji/go v1.1.4 // indirect
	golang.org/x/sys v0.0.0-20190429094411-2cc0cad0ac78 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
