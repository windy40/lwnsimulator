package gateway

import (
	"fmt"
	"log"
	"time"

	f "github.com/windy40/lwnsimulator/simulator/components/forwarder"
	"github.com/windy40/lwnsimulator/simulator/components/gateway/models"
	res "github.com/windy40/lwnsimulator/simulator/resources"
	"github.com/windy40/lwnsimulator/simulator/resources/communication/buffer"

	"github.com/windy40/lwnsimulator/simulator/util"
	"github.com/windy40/lwnsimulator/socket"
)

type Gateway struct {
	Id   int                `json:"id"`
	Info models.InfoGateway `json:"info"`

	State int `json:"-"`

	Resources *res.Resources `json:"-"` //is a pointer
	Forwarder *f.Forwarder   `json:"-"` //is a pointer

	Stat models.Stat `json:"-"`

	BufferUplink buffer.BufferUplink `json:"-"`
}

func (g *Gateway) CanExecute() bool {

	if g.State == util.Stopped {
		return false
	}

	return true

}

func (g *Gateway) Print(content string, err error, printType int) {

	now := time.Now()
	message := ""
	messageLog := ""
	event := socket.EventGw

	if err == nil {
		message = fmt.Sprintf("[ %s ] GW[%s]: %s", now.Format(time.Stamp), g.Info.Name, content)
		messageLog = fmt.Sprintf("GW[%s]: %s", g.Info.Name, content)
	} else {
		message = fmt.Sprintf("[ %s ] GW[%s] [ERROR]: %s", now.Format(time.Stamp), g.Info.Name, err)
		messageLog = fmt.Sprintf("GW[%s] [ERROR]: %s", g.Info.Name, err)
		event = socket.EventError
	}

	data := socket.ConsoleLog{
		Name: g.Info.Name,
		Msg:  message,
	}

	switch printType {

	case util.PrintBoth:
		g.Resources.WebSocket.Emit(event, data)
		log.Println(messageLog)

	case util.PrintOnlySocket:
		g.Resources.WebSocket.Emit(event, data)

	case util.PrintOnlyConsole:
		log.Println(messageLog)

	}

}
