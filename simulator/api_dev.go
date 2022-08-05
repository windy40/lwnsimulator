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

func LinkedDevName(d *dev.Device) string {
	name := d.Info.Name
	if d.Info.Status.LinkedDev {
		name = fmt.Sprintf("*%s", d.Info.Name)
	}
	return name
}
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
	d := s.Devices[Id]
	if d.Info.Status.LinkedDev {
		// only need to update dev socket
		log.Println(fmt.Sprintf("DEV[%s] dev socket updated %s", LinkedDevName(d), (*DevSocket).ID()))
		s.Resources.AddDevSocket(DevSocket, Id)

		if s.State == util.Running && !d.IsOn() {
			s.turnONDevice(Id)
		}
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdLinkDev, Error: codes.DevCmdOK})

		return
	}
	// device not yet linked to external dev
	if d.IsOn() {
		s.turnOFFDevice(Id)
	}

	d.Info.Status.LinkedDev = true
	s.Resources.AddDevSocket(DevSocket, Id)

	log.Println(fmt.Sprintf("DEV[%s] linked to external dev with socket %s", LinkedDevName(d), (*DevSocket).ID()))

	if s.State == util.Running {
		s.turnONDevice(Id)
	}

	s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdLinkDev, Error: codes.DevCmdOK})

	return
}

func (s *Simulator) DevExecuteUnlinkDev(DevSocket *socketio.Conn, Id int) {
	d := s.Devices[Id]
	if !d.Info.Status.LinkedDev {

		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdUnlinkDev, Error: codes.DevErrorDeviceNotLinked})

		return
	}
	// if device is on, first unjoin then turnOFF
	if d.IsOn() {
		if d.Info.Status.Joined {
			d.UnJoined()
			log.Println(fmt.Sprintf("DEV[%s] unjoined", LinkedDevName(d)))

		}
		s.turnOFFDevice(Id)
	}

	d.Info.Status.LinkedDev = true // keep it as Linkable device and turnedOFF

	log.Println(fmt.Sprintf("DEV[%s] unlinked from external dev on socket %s", LinkedDevName(d), (*DevSocket).ID()))

	s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdUnlinkDev, Error: codes.DevCmdOK})

	return
}

func (s *Simulator) DevExecuteJoinRequest(Id int) {
	d := s.Devices[Id]
	if !d.IsOn() {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", LinkedDevName(d), socket.DevCmdJoinRequest))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdJoinRequest, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	if d.Info.Status.Joined {
		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device already joined", LinkedDevName(d), socket.DevCmdJoinRequest))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdJoinRequest, Error: codes.DevErrorDeviceAlreadyJoined})
		return
	}

	go d.DevJoinAndProcessUplink()

	return
}

func (s *Simulator) DevExecuteSendUplink(Id int, data socket.DevExecuteSendUplink) {
	d := s.Devices[Id]
	mtype := data.MType
	payload := data.Payload
	if !d.IsOn() {

		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", LinkedDevName(d), data.Cmd))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: data.Cmd, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	if !d.Info.Status.Joined {

		log.Println(fmt.Sprintf("DEV[%s][CMD %s]Error dev not joined", LinkedDevName(d), socket.DevCmdSendUplink))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdSendUplink, Error: codes.DevErrorDeviceNotJoined})
		return
	}

	MType := lorawan.UnconfirmedDataUp
	if mtype == "ConfirmedDataUp" {
		MType = lorawan.ConfirmedDataUp
	}

	d.NewUplink(MType, payload)
	d.UplinkWaiting <- struct{}{}

	//d.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdSendUplink, Error: codes.DevCmdOK})
	//	s.Resources.WebSocket.Emit(socket.EventResponseCommand, "Uplink queued")

	return
}

func (s *Simulator) DevExecuteRecvDownlink(Id int, data socket.DevExecuteRecvDownlink) {

	d := s.Devices[Id]
	BufferSize := data.BufferSize

	if !d.IsOn() {

		log.Println(fmt.Sprintf("DEV[%s][CMD %s] device turned off", LinkedDevName(d), socket.DevCmdRecvDownlink))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdRecvDownlink, Error: codes.DevErrorDeviceTurnedOFF})
		return
	}

	if !d.Info.Status.Joined {

		log.Println(fmt.Sprintf("DEV[%s][CMD %s]Error dev not joined", LinkedDevName(d), socket.DevCmdRecvDownlink))
		s.Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseCmd{Cmd: socket.DevCmdRecvDownlink, Error: codes.DevErrorDeviceNotJoined})
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

	} else {

		s.Devices[d.Id].Resources.LinkedDevSocket[d.Id].Emit(socket.DevEventResponseCmd, socket.DevResponseRecvDownlink{Cmd: socket.DevEventRecvDownlink, Error: codes.DevErrorRecvBufferEmpty, MType: "", Payload: ""})

	}
	return
}

func (s *Simulator) DeleteDevSocket(SId string) {
	if Id, ok := s.Resources.ConnDevSocket[SId]; ok {
		d := s.Devices[Id]
		if d.Info.Status.Joined {
			d.UnJoined()
		}
		s.Resources.DeleteDevSocket(SId)
	}
}
