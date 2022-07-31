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
	DevErrorNoDeviceWithDevEUI
	DevErrorNIY
	DevErrorDeviceNotLinked
	DevErrorDeviceTurnedOFF
	DevErrorDeviceNotJoined
)
