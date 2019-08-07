package main

import (
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"ble-mesh/utils"
	"reflect"
	"strconv"
	"strings"

	"ble-mesh/driver"
	"ble-mesh/mesh"

	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
)

type (
	Application struct {
		mode int // 0:Adv 1:Gatt
	}
)

const adapterID = "hci1"

func createCmd() {

}

var (
	app    Application
	drv    *driver.Driver
	logger = utils.CreateLogger("main")
)

func cleanup() {
	mesh.OnClose()
}

func startMesh() {
	var bear mesh.Bear
	bear = &mesh.AdvertisingBear{}
	drv.Handle(driver.AdvertisementReceived(func(d []byte) {
		bear.OnPduReceived(d)
	}))
	bear.SetWriteHandle(drv.Advertise)
	mesh.SetNetworkBear(bear)
	mesh.StartMeshNetwork()
}

func provision(node string) {
	var bear mesh.Bear
	bear = &mesh.GattProxyBear{}
	session := drv.OpenProvisionGatt(node)
	if session == nil {
		return
	}
	session.RegisterGattDataEventHandler(bear.OnPduReceived)
	bear.SetWriteHandle(session.Write)
	bear.SetMTU(69)
	mesh.SetProvisionBear(bear)
	time.Sleep(time.Second)
	mesh.StartMeshProvision(node, func() {
		logger.Debug("provision finished")
		session.OnProvisionFinished()
	})
}

func main() {
	var err error
	homeDir, _ := homedir.Dir()
	if homeDir == "/root" {
		homeDir = "/home/xxx"
	}
	mesh.Init(homeDir + "/.config/ble-mesh")
	drv, err = driver.StartDiscovery()
	if err != nil {
		logger.Fatalf("Failed to open device, err: %s\n", err)
	}

	drv.Handle(driver.DevicePoweredOn(func() {
		startMesh()
	}))
	drv.Handle(driver.DeviceUnavailable(func() {

	}))
	drv.Handle(driver.UnProvNodeDiscovered(func(n *driver.UnprovisionedNode) {
		provision(n.UUID)
	}))

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	go startRouter()

	// select {}
	// test := true
	// for {
	// 	time.Sleep(3 * time.Second)
	// 	if test {
	// 		mesh.GenericOnOffSet(0x1000, 1)
	// 		// mesh.ConfigNetKeyAdd(1, 0x1000)
	// 		// mesh.ConfigGetComposition(0x1000)
	// 		test = false
	// 	}
	// }
	p := prompt.New(
		executor,
		completer,
		prompt.OptionHistory([]string{
			"ConfigBeaconSet 1000 0",
			"GenericOnOffGet 1000",
			"GenericOnOffSet 1000 0",
			// "prov 90a9f980-2c69-b64e-9a19-531a8a0ddc35",
			"ConfigBeaconGet 1000",
		}),
	)
	p.Run()
}

func callApi(fname string, args []string) {
	logger.Debugf("call func:%s, params:%+#v", fname, args)
	if f, ok := mesh.Functions[fname]; ok {
		t := f.Type()
		argVals := []reflect.Value{}
		var value reflect.Value
		if len(args) != t.NumIn() {
			logger.Error("wrong arguments")
			return
		}
		for i := 0; i < t.NumIn(); i++ {
			in := t.In(i)
			w := args[i]
			switch in.Kind() {
			case reflect.Bool:
				res := w == "1"
				value = reflect.ValueOf(res)
			case reflect.Uint:
				n, _ := strconv.ParseInt(w, 16, 32)
				res := uint(n)
				value = reflect.ValueOf(res)
			case reflect.Uint8:
				n, _ := strconv.ParseInt(w, 16, 32)
				res := byte(n)
				value = reflect.ValueOf(res)
			case reflect.Float32:
				n, _ := strconv.ParseFloat(w, 32)
				res := float32(n)
				value = reflect.ValueOf(res)
			}
			argVals = append(argVals, value)
		}
		// logger.Debugf("%+#v", args)
		ret := f.Call(argVals)
		err := ret[0].Interface()
		if err != nil {
			logger.Error(err)
		}
	}
}

func executor(t string) {
	words := strings.Split(t, " ")
	if t == "exit" {
		cleanup()
		os.Exit(0)
	} else if t[:4] == "prov" {
		provision(words[1])
	}
	callApi(words[0], words[1:])
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	text := t.TextBeforeCursor()
	spaces := strings.Count(text, " ")
	suggests := []prompt.Suggest{}
	if spaces == 0 {
		for name := range mesh.Functions {
			suggests = append(suggests, prompt.Suggest{Text: name})
		}
		suggests = append(suggests, prompt.Suggest{Text: "exit"})
	} else if spaces == 1 {
		suggests = append(suggests, prompt.Suggest{Text: "1000"})
	}
	// logger.Debugf("%+#v", suggests)
	sep := `\w*`
	regex := sep
	for _, c := range t.GetWordBeforeCursor() {
		regex += string(c) + sep
	}
	res := funk.Filter(suggests, func(x interface{}) bool {
		m, _ := regexp.Match(regex, []byte(strings.ToLower(x.(prompt.Suggest).Text)))
		return m
	})
	return res.([]prompt.Suggest)
	// return prompt.FilterHasPrefix(suggests, t.GetWordBeforeCursor(), true)
}

func startRouter() {
	router := gin.Default()
	router.GET("/api", func(c *gin.Context) {
		params := c.Query("params")
		callApi(c.Query("f"), strings.Split(params, ","))
		c.Status(http.StatusOK)
	})
	router.GET("/log", func(c *gin.Context) {
		c.String(http.StatusOK, utils.ReadAllLogs())
	})
	router.GET("/getdb", func(c *gin.Context) {
		db := &mesh.Mesh{}
		copier.Copy(db, mesh.GetDb())
		for _, k := range db.NetKeys {
			k.OldKey = nil
			k.NewKey = nil
		}
		for _, k := range db.AppKeys {
			k.OldKey = nil
			k.NewKey = nil
		}
		for _, n := range db.Nodes {
			for _, e := range n.Elements {
				e.Node = nil
				for _, m := range e.Models {
					m.Element = nil
				}
			}
		}
		c.JSON(http.StatusOK, db)
	})
	router.GET("/getnode", func(c *gin.Context) {
		node := utils.HexStringToUint(c.Query("n"))
		if node != 0 {
			n := mesh.GetNode(node)
			for _, e := range n.Elements {
				e.Node = nil
				for _, m := range e.Models {
					m.Element = nil
				}
			}
			c.JSON(http.StatusOK, n)
		} else {
			c.String(http.StatusNotFound, "{}")
		}
	})
	router.POST("/setnode", func(c *gin.Context) {
		var node mesh.Node
		err := c.BindJSON(&node)
		if err != nil {
			c.String(http.StatusInternalServerError, "invalid settings")
			return
		}
		err = mesh.SetNode(&node)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			logger.Error(err)
			return
		}
		c.JSON(http.StatusOK, nil)
	})
	router.Run(":15031")
}
