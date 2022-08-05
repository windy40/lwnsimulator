package socket

const (
	DevEventLog   = "dev-log"
	DevEventError = "dev-error"

	DevEventLinkDev      = "link-dev"
	DevEventUnlinkDev    = "unlink-dev"
	DevEventJoinRequest  = "join-request"
	DevEventSendUplink   = "send-uplink"
	DevEventRecvDownlink = "recv-downlink"

	DevEventExecuteCmd  = "execute-cmd"
	DevEventResponseCmd = "response-cmd"

	DevEventAckCmd = "ack-cmd"
	DevEventLoRa   = "lora-event"
)

const (
	DevCmdLinkDev      = "link-dev"
	DevCmdUnlinkDev    = "unlink-dev"
	DevCmdJoinRequest  = "join-request"
	DevCmdSendUplink   = "send-uplink"
	DevCmdRecvDownlink = "recv-downlink"
)

const (
	RX_PACKET_EVENT   = 1
	TX_PACKET_EVENT   = 2
	TX_FAILED_EVENT   = 4
	JOIN_ACCEPT_EVENT = 16
)
