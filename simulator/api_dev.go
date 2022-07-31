package simulator

import (
	"errors"
	"fmt"
	"log"

	"github.com/brocaar/lorawan"
	socketio "github.com/googollee/go-socket.io"
	"github.com/windy40/lwnsimulator/codes"
	dev "github.com/windy40/lwnsimulator/simulator/components/device"

	"github.com/windy40/lwnsimulator/simulator/util"
	"github.com/windy40/lwnsimulator/socket"
)

func (s *Simulator) getDeviceWithDevEUI(devEUIstr string) (dev *dev.Device, err error) {
	devices := s.Devices

	var devEUI lorawan.EUI64
	devEUI.UnmarshalText([]byte(devEUIstr))

	// get device Id for devEUI device if it exists
	for _, d := range devices {

		if d.Info.DevEUI == devEUI {
			return d, nil
		}
	}
	err = errors.New("No device with given devEUI")
	return nil, err
}

func (s *Simulator) DevExecuteLinkDev(DevSocket *socketio.Conn, Id int) {

	if s.Devices[Id].Info.Status.LinkedDev {
		// only need to update dev socket
		log.Println(fmt.Sprintf("DEV[*%s] dev socket updated %s", s.Devices[Id].Info.Name, (*DevSocket).ID()))
		s.Resources.AddDevSocket(DevSocket, Id)
		return
	}
	// device not yet linked to external dev
	if s.Devices[Id].IsOn() {
		s.turnOFFDevice(Id)
	}

	s.Devices[Id].Info.Status.LinkedDev = true
	s.Resources.AddDevSocket(DevSocket, Id)

	log.Println(fmt.Sprintf("DEV[%s] linked to external dev with socket %s", s.Devices[Id].Info.Name, (*DevSocket).ID()))

	if s.State == util.Running {
		s.turnONDevice(Id)
	}

	s.Resources.LinkedDevSocket[Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdLinkDev, Error: codes.DevCmdOK})

}

func (s *Simulator) DeleteDevSocket(SId string) {
	s.Resources.DeleteDevSocket(SId)
}

func (s *Simulator) DevExecuteJoinRequest(Id int) {
	d := s.Devices[Id]
	if !d.IsOn() {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", d.Info.Name, socket.DevCmdJoinRequest))
		d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdJoinRequest, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	go d.DevJoinAndProcessUplink()
}

func (s *Simulator) DevExecuteSendUplink(Id int, data socket.DevExecuteSendUplink) {
	d := s.Devices[Id]
	mtype := data.MType
	payload := data.Payload
	if !d.IsOn() {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", d.Info.Name, data.Cmd))
		d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: data.Cmd, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	if !d.Info.Status.Joined {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s]Error dev not joined", d.Info.Name, socket.DevCmdSendUplink))
		d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdSendUplink, Error: codes.DevErrorDeviceNotJoined})
		return
	}

	MType := lorawan.UnconfirmedDataUp
	if mtype == "ConfirmedDataUp" {
		MType = lorawan.ConfirmedDataUp
	}

	d.NewUplink(MType, payload)
	d.UplinkWaiting <- struct{}{}

	d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdSendUplink, Error: codes.DevCmdOK})

	//	s.Resources.WebSocket.Emit(socket.EventResponseCommand, "Uplink queued")

}

func (s *Simulator) DevExecuteRecvDownlink(Id int, data socket.DevExecuteRecvDownlink) {

	d := s.Devices[Id]
	BufferSize := data.BufferSize

	if !d.IsOn() {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", d.Info.Name, socket.DevCmdRecvDownlink))
		d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdRecvDownlink, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	if !d.Info.Status.Joined {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s]Error dev not joined", d.Info.Name, socket.DevCmdRecvDownlink))
		d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdRecvDownlink, Error: codes.DevErrorDeviceNotJoined})
		return
	}

	if len(d.Info.Status.BufferDataDownlinks) > 0 {

		payload := d.Info.Status.BufferDataDownlinks[0].DataPayload
		size := len(payload)
		if size > BufferSize {
			size = BufferSize
		}
		payload = payload[:size]

		mtype := "UnconfirmedDataDown"
		if d.Info.Status.BufferDataDownlinks[0].MType == lorawan.ConfirmedDataDown {
			mtype = "ConfirmedDataDown"
		}

		switch len(d.Info.Status.BufferDataDownlinks) {
		case 1:
			d.Info.Status.BufferDataDownlinks = d.Info.Status.BufferDataDownlinks[:0]

		default:
			d.Info.Status.BufferDataDownlinks = d.Info.Status.BufferDataDownlinks[1:]

		}

		s.Devices[d.Id].Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseRecvDownlink{Cmd: socket.DevEventRecvDownlink, Error: codes.DevCmdOK, MType: mtype, Payload: string(payload)})

	}

}
