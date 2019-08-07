package driver

import (
	"ble-mesh/mesh"
	"ble-mesh/utils"
	"encoding/binary"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/linux/cmd"
	"github.com/google/uuid"
)

type (
	OnDataReceived func(data []byte)
	Option         func(*Driver) error
	Handler        func(*Driver)

	UnprovisionedNode struct {
		Name string
		Mac  string
		UUID string
		OOB  uint16
		Hash uint32
		Adv  bool
		Gatt bool
		p    gatt.Peripheral
	}

	Session struct {
		node    *UnprovisionedNode
		driver  *Driver
		p       gatt.Peripheral
		uuidSvc string
		ch      chan bool

		onDataRecvd func(data []byte)

		charaOut *gatt.Characteristic
		charaIn  *gatt.Characteristic
	}

	Driver struct {
		dev           gatt.Device
		connecting    bool
		sessions      map[string]*Session
		activeSession *Session

		unprovNodes map[string]*UnprovisionedNode

		devicePoweredOn      func()
		deviceUnavailable    func()
		advertismentReceived func(data []byte)
		unProvNodeDiscovered func(*UnprovisionedNode)
	}
)

const (
	UUID_MESH_PROXY        = "1828"
	UUID_MESH_PROVISIONING = "1827"

	UUID_MESH_PROV_DATA_IN   = "2adb"
	UUID_MESH_PROV_DATA_OUT  = "2adc"
	UUID_MESH_PROXY_DATA_IN  = "2add"
	UUID_MESH_PROXY_DATA_OUT = "2ade"
)

var defaultClientOptions = []gatt.Option{
	gatt.LnxMaxConnections(1),
	gatt.LnxDeviceID(-1, true),
}
var logger *logrus.Entry

func (d *Driver) onStateChanged(dev gatt.Device, s gatt.State) {
	logger.Debug("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		logger.Debug("scanning...")
		dev.Scan([]gatt.UUID{}, false)
		if d.devicePoweredOn != nil {
			d.devicePoweredOn()
		}
		return
	default:
		dev.StopScanning()
		dev.StopAdvertising()
		if d.deviceUnavailable != nil {
			d.deviceUnavailable()
		}
	}
}

func (d *Driver) saveUnprovNode(uuid, name, mac string, oob uint16, adv, gatt bool, p gatt.Peripheral) {
	node, ok := d.unprovNodes[uuid]
	if !ok {
		node = &UnprovisionedNode{}
	}

	node.Name = name
	node.Mac = mac
	node.UUID = uuid
	node.OOB = oob
	node.p = p
	node.Adv = adv
	node.Gatt = gatt
	if !ok {
		d.unprovNodes[uuid] = node
		logger.Debugf("unprovisoned node: %+#v", node)
		if d.unProvNodeDiscovered != nil {
			d.unProvNodeDiscovered(node)
		}
	}

}

func (d *Driver) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if len(a.Raw) > 2 && a.Raw[1] == mesh.BT_LE_ADV_BEACON && a.Raw[2] == 0 {
		// todo: hash
		uuid, _ := uuid.FromBytes(a.Raw[3:19])
		if n, ok := d.unprovNodes[uuid.String()]; ok && n.Adv {
			return
		}
		oob := binary.BigEndian.Uint16(a.Raw[19:21])
		d.saveUnprovNode(uuid.String(), p.Name(), p.ID(), oob, true, false, p)
	}
	for _, s := range a.ServiceData {
		if s.UUID.String() == UUID_MESH_PROVISIONING {
			uuid, _ := uuid.FromBytes(s.Data[:16])
			if n, ok := d.unprovNodes[uuid.String()]; ok && n.Gatt {
				return
			}
			oob := binary.BigEndian.Uint16(s.Data[16:])

			d.saveUnprovNode(uuid.String(), p.Name(), p.ID(), oob, false, true, p)
		}
	}
	if d.advertismentReceived != nil {
		d.advertismentReceived(a.Raw)
	}
}

func (d *Driver) onPeriphConnected(p gatt.Peripheral, err error) {
	logger.Debug("Connected")
	defer func() { d.connecting = false }()
	// defer p.Device().CancelConnection(p)
	session := d.activeSession
	if err := p.SetMTU(69); err != nil {
		logger.Errorf("Failed to set MTU, err: %s\n", err)
		session.ch <- false
		return
	}

	// Discovery services
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		logger.Errorf("Failed to discover services, err: %s\n", err)
		session.ch <- false
		return
	}

	for _, s := range ss {
		if s.UUID().String() != session.uuidSvc {
			continue
		}
		var uuidIn, uuidOut string
		if s.UUID().String() == UUID_MESH_PROVISIONING {
			uuidIn = UUID_MESH_PROV_DATA_IN
			uuidOut = UUID_MESH_PROV_DATA_OUT
		} else {
			uuidIn = UUID_MESH_PROXY_DATA_IN
			uuidOut = UUID_MESH_PROXY_DATA_OUT
		}

		// Discovery characteristics
		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			logger.Errorf("Failed to discover characteristics, err: %s\n", err)
			session.ch <- false
			return
		}

		for _, c := range cs {
			if c.UUID().String() == uuidIn {
				session.charaIn = c
				if c.Properties()&(gatt.CharWrite|gatt.CharWriteNR) == 0 {
					logger.Errorf("characteristics dose not support write")
					session.ch <- false
					return
				}
			} else if c.UUID().String() == uuidOut {
				// Discovery descriptors
				_, err := p.DiscoverDescriptors(nil, c)
				if err != nil {
					logger.Errorf("Failed to discover descriptors, err: %s\n", err)
					session.ch <- false
					return
				}
				session.charaOut = c
				if (c.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
					f := func(c *gatt.Characteristic, b []byte, err error) {
						if session.onDataRecvd != nil {
							session.onDataRecvd(b)
						}
					}
					if err := p.SetNotifyValue(c, f); err != nil {
						logger.Debugf("Failed to subscribe characteristic, err: %s\n", err)
						session.ch <- false
						return
					}
				}
				logger.Debug("subscribe successful")
			}
		}
		session.p = p
		session.ch <- true
	}

}

func (d *Driver) onPeriphDisconnected(p gatt.Peripheral, err error) {
}

func (d *Driver) GetUnprovNodes() []*UnprovisionedNode {
	nodes := []*UnprovisionedNode{}
	for _, n := range d.unprovNodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (d *Driver) OpenProvisionGatt(uuid string) *Session {
	// d.dev.StopScanning()
	// d.dev.StopAdvertising()
	defer func() { d.activeSession = nil }()
	if d.connecting {
		logger.Error("connecting to other device is not finished, please retry later")
		return nil
	}
	if n, ok := d.unprovNodes[uuid]; ok {
		d.connecting = true
		session := &Session{uuidSvc: UUID_MESH_PROVISIONING, ch: make(chan bool)}
		d.activeSession = session
		d.dev.Connect(n.p)

		var success bool
		to := time.NewTimer(time.Second * 5)
		select {
		case success = <-d.activeSession.ch:
			if !success {
				return nil
			}
		case <-to.C:
			logger.Error("connect timeout")
			return nil
		}

		if d.sessions == nil {
			d.sessions = map[string]*Session{}
		}
		session.node = n
		session.driver = d
		d.sessions[uuid] = session
		return session
	}
	return nil
}

func (s *Session) OnProvisionFinished() {
	if s.p != nil {
		s.p.Device().CancelConnection(s.p)
	}
	delete(s.driver.unprovNodes, s.node.UUID)
	delete(s.driver.sessions, s.node.UUID)
}

func (s *Session) GetMTU() uint {
	return 2
}

func (d *Driver) Advertise(data []byte) error {
	packet := &gatt.AdvPacket{}
	packet.AppendField(data[1], data[2:])
	err := d.dev.Advertise(packet)
	return err
}

func (d *Driver) Stop() {
	d.dev.Stop()
}

func (d *Driver) Handle(hh ...Handler) {
	for _, h := range hh {
		h(d)
	}
}

func DevicePoweredOn(f func()) Handler {
	return func(d *Driver) { d.devicePoweredOn = f }
}

func DeviceUnavailable(f func()) Handler {
	return func(d *Driver) { d.deviceUnavailable = f }
}

func AdvertisementReceived(f func([]byte)) Handler {
	return func(d *Driver) { d.advertismentReceived = f }
}

func UnProvNodeDiscovered(f func(*UnprovisionedNode)) Handler {
	return func(d *Driver) { d.unProvNodeDiscovered = f }
}

func (s *Session) RegisterGattDataEventHandler(h OnDataReceived) {
	s.onDataRecvd = h
}

func (s *Session) Write(data []byte) error {
	err := s.p.WriteCharacteristic(s.charaIn, data, true)
	return err
}

func StartDiscovery() (*Driver, error) {
	logger = utils.CreateLogger("driver")
	driver := &Driver{}
	driver.unprovNodes = map[string]*UnprovisionedNode{}
	d, err := gatt.NewDevice(defaultClientOptions...)
	if err != nil {
		return nil, err
	}

	d.Handle(
		gatt.PeripheralDiscovered(driver.onPeriphDiscovered),
		gatt.PeripheralConnected(driver.onPeriphConnected),
		gatt.PeripheralDisconnected(driver.onPeriphDisconnected),
	)

	d.Init(driver.onStateChanged)

	o := gatt.LnxSetAdvertisingParameters(&cmd.LESetAdvertisingParameters{
		AdvertisingIntervalMin:  0x100,     // [0x0800]: 0.625 ms * 0x0800 = 1280.0 ms
		AdvertisingIntervalMax:  0x300,     // [0x0800]: 0.625 ms * 0x0800 = 1280.0 ms
		AdvertisingType:         0x03,      // [0x00]: ADV_IND, 0x01: DIRECT(HIGH), 0x02: SCAN, 0x03: NONCONN, 0x04: DIRECT(LOW)
		OwnAddressType:          0x00,      // [0x00]: public, 0x01: random
		DirectAddressType:       0x00,      // [0x00]: public, 0x01: random
		DirectAddress:           [6]byte{}, // Public or Random Address of the device to be connected
		AdvertisingChannelMap:   0x7,       // [0x07] 0x01: ch37, 0x02: ch38, 0x04: ch39
		AdvertisingFilterPolicy: 0x00,
	})
	d.Option(o)

	driver.dev = d

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	driver.Provision("a8017135-0200-0089-ebc0-07da78000000")
	// }()

	return driver, nil
}
