package device

import (
	"github.com/windy40/lwnsimulator/simulator/components/device/classes"
	"github.com/windy40/lwnsimulator/simulator/util"
	"github.com/windy40/lwnsimulator/socket"
)

func (d *Device) DevJoinAndProcessUplink() {

	defer d.Resources.ExitGroup.Done()

	d.Print("trying to join ...", nil, util.PrintBoth)
	d.OtaaActivation()

	if d.Info.Status.Joined {
		if d.Info.Status.LinkedDev {
			d.ReturnLoraEvent(socket.JOIN_ACCEPT_EVENT)

		} else {
			d.Print("Could not join", nil, util.PrintBoth)
			return

		}

	}
	for {

		select {

		case <-d.UplinkWaiting:
			break

		case <-d.Exit:
			d.Print("Turn OFF", nil, util.PrintBoth)
			return
		}

		if d.CanExecute() {

			if d.Info.Status.Joined {

				if d.Info.Configuration.SupportedClassC {
					d.SwitchClass(classes.ClassC)
				} else if d.Info.Configuration.SupportedClassB {
					d.SwitchClass(classes.ClassB)
				}

				d.Execute()

			} else {
				d.OtaaActivation()
			}

		}
	}
}
