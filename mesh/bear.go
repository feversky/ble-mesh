package mesh

type (
	Bear interface {
		Start()
		Stop()
		OnPduReceived(pdu []byte)
		SetWriteHandle(handle func([]byte) error)
		SetMTU(mtu uint)
		SendNetPdu(pdu []byte)
		SendProvPdu(pdu []byte)
	}
)
