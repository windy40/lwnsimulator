package codes

const (
	DevCodeNOK = iota
	DevCodeOK
)

const (
	DevCodeLinkedDevOK = iota
	DevCodeDevJoined
)

const (
	DevCmdOK = iota
	DevCmdTimeout
	DevErrorNoDeviceWithDevEUI
	DevErrorNIY
	DevErrorDeviceNotLinked
	DevErrorDeviceTurnedOFF
	DevErrorDeviceNotJoined
	DevErrorDeviceAlreadyJoined
	DevErrorRecvBufferEmpty
)
