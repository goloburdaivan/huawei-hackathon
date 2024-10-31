package structs

type PortInfo struct {
	Index              int
	Name               string
	OID                string
	InOctets           uint
	OutOctets          uint
	InErrors           uint
	OutErrors          uint
	InUcastPkts        uint
	OutUcastPkts       uint
	InMulticastPkts    uint
	OutMulticastPkts   uint
	InBroadcastPkts    uint
	OutBroadcastPkts   uint
	AdminStatus        string
	OperStatus         string
	InOctetsPkts       uint
	OutOctetsPkts      uint
	InBandwidthUtil    float32
	OutBandwidthUtil   float32
	InBandwidthActual  float32
	OutBandwidthActual float32
}
