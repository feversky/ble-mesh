package db

import (
	"encoding/json"
	"io/ioutil"
)

type Mesh struct {
	MeshName       string      `json:"meshName"`
	NetKeys        []NetKey    `json:"netKeys"`
	AppKeys        []AppKey    `json:"appKeys"`
	Groups         []Group     `json:"groups"`
	Provisioner    Provisioner `json:"provisioner"`
	IVindex        uint        `json:"IVindex"`
	IVupdate       uint        `json:"IVupdate"`
	SequenceNumber uint        `json:"sequenceNumber"`
}
type NetKey struct {
	Index           uint   `json:"index"`
	KeyRefreshPhase uint   `json:"keyRefreshPhase"`
	Key             string `json:"key"`
	OldKey          string `json:"oldkey"`
}
type AppKey struct {
	Index           uint   `json:"index"`
	KeyRefreshPhase uint   `json:"keyRefreshPhase"`
	Key             string `json:"key"`
	OldKey          string `json:"oldkey"`
}
type DevKey struct {
	AppKey
}
type Provisioner struct {
	ProvisionerName string `json:"provisionerName"`
	UnicastAddress  string `json:"unicastAddress"`
	LowAddress      string `json:"lowAddress"`
	HighAddress     string `json:"highAddress"`
}

type Node struct {
	DeviceKey                    string         `json:"deviceKey"`
	BindedNetKeys                []BindedNetKey `json:"bindedNetKeys"`
	UnicastAddress               string         `json:"unicastAddress"`
	Elements                     []Element      `json:"elements"`
	Cid                          int            `json:"cid"`
	Pid                          int            `json:"pid"`
	Vid                          int            `json:"vid"`
	Crpl                         int            `json:"crpl"`
	Features                     Features       `json:"features"`
	SequenceNumber               uint           `json:"sequenceNumber"`
	Mac                          string         `json:"mac"`
	UUID                         string         `json:"uuid"`
	LPN                          bool           `json:"lpn"`
	Friend                       string         `json:"friend"`
	TTL                          uint           `json:"ttl"`
	Relay                        uint           `json:"relay"`
	RelayRetransmitCount         uint           `json:"relayRetransmitCount"`
	RelayRetransmitIntervalSteps uint           `json:"relayRetransmitIntervalSteps"`
	AttentionTimer               uint           `json:"attentionTimer"`
	SecureNetworkBeacon          bool           `json:"secureNetworkBeacon"`
	GATTProxyState               uint           `json:"gattProxyState"`
	FriendState                  uint           `json:"friendState"`
	KeyRefreshPhaseState         uint           `json:"keyRefreshPhaseState"`
	NetworkTransmitCount         uint           `json:"networkTransmitCount"`
	NetworkTransmitIntervalSteps uint           `json:"networkTransmitIntervalSteps"`
	CurrentFault                 uint           `json:"currentFault"`
}

type Group struct {
	GroupAddress string `json:"groupAddress"`
	Name         string `json:"name"`
}
type BindedNetKey struct {
	NetKeyIndex       uint   `json:"netKeyIndex"`
	BindedAppKeys     []uint `json:"bindedAppKeys"`
	NodeIdentityState uint   `json:"NodeIdentityState"`
}
type Model struct {
	ModelID       string     `json:"modelId"`
	BindedAppKeys []uint     `json:"bindedAppKeys"`
	PubSetting    PubSetting `json:"publish"`
	SubAddresses  []string   `json:"subAddresses"`
	State         string     `json:"state"`
}
type Element struct {
	ElementIndex   int     `json:"elementIndex"`
	Location       int     `json:"Location"`
	UnicastAddress string  `json:"unicastAddress"`
	Models         []Model `json:"models"`
}
type Features struct {
	Relay  bool `json:"relay"`
	Proxy  bool `json:"proxy"`
	Friend bool `json:"friend"`
	Lpn    bool `json:"lpn"`
}
type PubSetting struct {
	PublishAddress                 uint `json:"publishAddress"`
	AppKeyIndex                    uint `json:"appKeyIndex"`
	CredentialFlag                 uint `json:"credentialFlag"`
	PublishTTL                     uint `json:"publishTTL"`
	PublishNumberOfSteps           uint `json:"publishNumberOfSteps"`
	PublishStepResolution          uint `json:"publishStepResolution"`
	PublishRetransmitCount         uint `json:"publishRetransmitCount"`
	PublishRetransmitIntervalSteps uint `json:"publishRetransmitIntervalSteps"`
}

func ReadFromDb(path string, obj interface{}) error {
	jsonData, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(jsonData, obj)
	return err
}

func WriteToDb(path string, obj interface{}) error {
	jsonData, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonData, 0744)
	return err
}
