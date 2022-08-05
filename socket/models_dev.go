package socket

type DevExecuteCmdInter interface {
	GetCmd() string
	GetAck() bool
	GetDevEUI() string
}

type DevExecuteCmd struct {
	Cmd    string `json:"cmd"`
	Ack    bool   `json:"ack"`
	DevEUI string `json:"devEUI"`
}

func (c DevExecuteCmd) GetCmd() string {
	return c.Cmd
}
func (c DevExecuteCmd) GetAck() bool {
	return c.Ack
}
func (c DevExecuteCmd) GetDevEUI() string {
	return c.DevEUI
}

type DevExecuteRecvDownlink struct {
	Cmd        string `json:"cmd"`
	Ack        bool   `json:"ack"`
	DevEUI     string `json:"devEUI"`
	BufferSize int    `json:"buffersize"`
}

func (c DevExecuteRecvDownlink) GetCmd() string {
	return c.Cmd
}
func (c DevExecuteRecvDownlink) GetAck() bool {
	return c.Ack
}
func (c DevExecuteRecvDownlink) GetDevEUI() string {
	return c.DevEUI
}

type DevExecuteSendUplink struct {
	Cmd     string `json:"cmd"`
	DevEUI  string `json:"devEUI"`
	Ack     bool   `json:"ack"`
	MType   string `json:"mtype"`
	Payload string `json:"payload"`
}

func (c DevExecuteSendUplink) GetCmd() string {
	return c.Cmd
}
func (c DevExecuteSendUplink) GetAck() bool {
	return c.Ack
}
func (c DevExecuteSendUplink) GetDevEUI() string {
	return c.DevEUI
}

type DevAckCmd struct {
	Cmd  string `json:"cmd"`
	Args string `json:"args"`
}

type DevResponseCmd struct {
	Cmd   string `json:"cmd"`
	Error int    `json:"error"`
}

type DevResponseRecvDownlink struct {
	Cmd     string `json:"cmd"`
	Error   int    `json:"error"`
	MType   string `json:"mtype"`
	Payload string `json:"payload"`
}

type DevLoRaEvent struct {
	Event int `json:"event"`
}
