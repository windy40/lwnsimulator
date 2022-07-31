package device

import (
	"github.com/windy40/lwnsimulator/simulator/util"

	"github.com/brocaar/lorawan"
	act "github.com/windy40/lwnsimulator/simulator/components/device/activation"
	"github.com/windy40/lwnsimulator/simulator/components/device/classes"
	dl "github.com/windy40/lwnsimulator/simulator/components/device/frames/downlink"

	"github.com/windy40/lwnsimulator/socket"
)

func (d *Device) ProcessDownlink(phy lorawan.PHYPayload) (*dl.InformationDownlink, error) {

	var payload *dl.InformationDownlink
	var err error

	mtype := phy.MHDR.MType
	err = nil

	switch mtype {

	case lorawan.JoinAccept:
		Ja, err := act.DecryptJoinAccept(phy, d.Info.DevNonce, d.Info.JoinEUI, d.Info.AppKey)
		if err != nil {
			return nil, err
		}

		return d.ProcessJoinAccept(Ja)

	case lorawan.UnconfirmedDataDown:

		payload, err = dl.GetDownlink(phy, d.Info.Configuration.DisableFCntDown, d.Info.Status.FCntDown,
			d.Info.NwkSKey, d.Info.AppSKey)
		if err != nil {
			return nil, err
		}

		if d.Info.Status.LinkedDev {
			d.Info.Status.BufferDataDownlinks = append(d.Info.Status.BufferDataDownlinks, *payload)
			d.ReturnLoraEvent(socket.RX_PACKET_EVENT)
		}

	case lorawan.ConfirmedDataDown: //ack

		payload, err = dl.GetDownlink(phy, d.Info.Configuration.DisableFCntDown, d.Info.Status.FCntDown,
			d.Info.NwkSKey, d.Info.AppSKey)
		if err != nil {
			return nil, err
		}

		d.SendAck()

		if d.Info.Status.LinkedDev {
			d.Info.Status.BufferDataDownlinks = append(d.Info.Status.BufferDataDownlinks, *payload)
			d.ReturnLoraEvent(socket.RX_PACKET_EVENT)
		}

	}

	d.Info.Status.FCntDown = (d.Info.Status.FCntDown + 1) % util.MAXFCNTGAP

	switch d.Class.GetClass() {

	case classes.ClassA:
		d.Info.Status.DataUplink.AckMacCommand.CleanFOptsRXParamSetupAns()
		d.Info.Status.DataUplink.AckMacCommand.CleanFOptsRXTimingSetupAns()
		break

	case classes.ClassC:
		d.Info.Status.InfoClassC.SetACK(false) //Reset

	}

	msg := d.Info.Status.DataUplink.ADR.Reset()
	if msg != "" {
		d.Print(msg, nil, util.PrintBoth)
	}

	d.Info.Status.DataUplink.AckMacCommand.CleanFOptsDLChannelAns()

	return payload, err
}
