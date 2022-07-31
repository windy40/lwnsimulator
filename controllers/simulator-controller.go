package controllers

import (
	"github.com/brocaar/lorawan"
	socketio "github.com/googollee/go-socket.io"
	"github.com/windy40/lwnsimulator/models"
	repo "github.com/windy40/lwnsimulator/repositories"
	dev "github.com/windy40/lwnsimulator/simulator/components/device"
	gw "github.com/windy40/lwnsimulator/simulator/components/gateway"
	e "github.com/windy40/lwnsimulator/socket"
)

//SimulatorController interfaccia controller
type SimulatorController interface {
	Run() bool
	Stop() bool
	GetIstance()

	AddWebSocket(*socketio.Conn)
	// windy40 dev sockets
	DevExecuteLinkDev(*socketio.Conn, int)
	DeleteDevSocket(string)
	DevExecuteJoinRequest(int)
	DevExecuteSendUplink(int, e.DevExecuteSendUplink)
	DevExecuteRecvDownlink(int, e.DevExecuteRecvDownlink)
	// windy40
	SaveBridgeAddress(models.AddressIP) error
	GetBridgeAddress() models.AddressIP
	GetGateways() []gw.Gateway
	AddGateway(*gw.Gateway) (int, int, error)
	UpdateGateway(*gw.Gateway) (int, error)
	DeleteGateway(int) bool
	AddDevice(*dev.Device) (int, int, error)
	GetDevices() []dev.Device
	UpdateDevice(*dev.Device) (int, error)
	DeleteDevice(int) bool
	ToggleStateDevice(int)
	SendMACCommand(lorawan.CID, e.MacCommand)
	ChangePayload(e.NewPayload) (string, bool)
	SendUplink(e.NewPayload)
	ChangeLocation(e.NewLocation) bool
	ToggleStateGateway(int)
}

type simulatorController struct {
	repo repo.SimulatorRepository
}

//NewSimulatorController return il controller
func NewSimulatorController(repo repo.SimulatorRepository) SimulatorController {
	return &simulatorController{
		repo: repo,
	}
}

func (c *simulatorController) GetIstance() {
	c.repo.GetIstance()
}

func (c *simulatorController) AddWebSocket(socket *socketio.Conn) {
	c.repo.AddWebSocket(socket)
}

// windy40 dev sockets

func (c *simulatorController) DevExecuteLinkDev(socket *socketio.Conn, Id int) {
	c.repo.DevExecuteLinkDev(socket, Id)
}

func (c *simulatorController) DeleteDevSocket(SId string) {
	c.repo.DeleteDevSocket(SId)
}

func (c *simulatorController) DevExecuteJoinRequest(Id int) {
	c.repo.DevExecuteJoinRequest(Id)
}

func (c *simulatorController) DevExecuteSendUplink(Id int, data e.DevExecuteSendUplink) {
	c.repo.DevExecuteSendUplink(Id, data)
}
func (c *simulatorController) DevExecuteRecvDownlink(Id int, data e.DevExecuteRecvDownlink) {
	c.repo.DevExecuteRecvDownlink(Id, data)
}

// windy40

func (c *simulatorController) Run() bool {
	return c.repo.Run()
}

func (c *simulatorController) Stop() bool {
	return c.repo.Stop()
}

func (c *simulatorController) IsRunning() bool {
	return c.repo.Stop()
}

func (c *simulatorController) SaveBridgeAddress(addr models.AddressIP) error {
	return c.repo.SaveBridgeAddress(addr)
}

func (c *simulatorController) GetBridgeAddress() models.AddressIP {
	return c.repo.GetBridgeAddress()
}

func (c *simulatorController) GetGateways() []gw.Gateway {
	return c.repo.GetGateways()
}

func (c *simulatorController) AddGateway(gateway *gw.Gateway) (int, int, error) {
	return c.repo.AddGateway(gateway)
}

func (c *simulatorController) UpdateGateway(gateway *gw.Gateway) (int, error) {
	return c.repo.UpdateGateway(gateway)
}

func (c *simulatorController) DeleteGateway(Id int) bool {
	return c.repo.DeleteGateway(Id)
}

func (c *simulatorController) AddDevice(device *dev.Device) (int, int, error) {
	return c.repo.AddDevice(device)
}

func (c *simulatorController) GetDevices() []dev.Device {
	return c.repo.GetDevices()
}

func (c *simulatorController) UpdateDevice(device *dev.Device) (int, error) {
	return c.repo.UpdateDevice(device)
}

func (c *simulatorController) DeleteDevice(Id int) bool {
	return c.repo.DeleteDevice(Id)
}

func (c *simulatorController) ToggleStateDevice(Id int) {
	c.repo.ToggleStateDevice(Id)
}

func (c *simulatorController) SendMACCommand(cid lorawan.CID, data e.MacCommand) {
	c.repo.SendMACCommand(cid, data)
}

func (c *simulatorController) ChangePayload(pl e.NewPayload) (string, bool) {
	return c.repo.ChangePayload(pl)
}

func (c *simulatorController) SendUplink(pl e.NewPayload) {
	c.repo.SendUplink(pl)
}

func (c *simulatorController) ChangeLocation(loc e.NewLocation) bool {
	return c.repo.ChangeLocation(loc)
}

func (c *simulatorController) ToggleStateGateway(Id int) {
	c.repo.ToggleStateGateway(Id)
}
