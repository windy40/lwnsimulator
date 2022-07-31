package socket

const (
	DevEventLog   = "dev-log"
	DevEventError = "dev-error"
	DevEventTest  = "test"

	DevEventLinkDev      = "link-dev"
	DevEventJoinRequest  = "join-request"
	DevEventSendUplink   = "send-uplink"
	DevEventRecvDownlink = "recv-downlink"
	/*
		DevEventExecuteCmd  = "execute-cmd"
	*/
	DevEventResponseCmd = "response-cmd"

	DevEventAckCmd = "ack-cmd"
	DevEventLoRa   = "lora-event"
)

const (
	DevCmdLinkDev      = "link-dev"
	DevCmdJoinRequest  = "join-request"
	DevCmdSendUplink   = "send-uplink"
	DevCmdRecvDownlink = "recv-downlink"
)

const (
	RX_PACKET_EVENT = iota
	TX_PACKET_EVENT
	TX_FAILED_EVENT
	JOIN_ACCEPT_EVENT
)
