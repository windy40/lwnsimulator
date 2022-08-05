package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/brocaar/lorawan"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rakyll/statik/fs"
	cnt "github.com/windy40/lwnsimulator/controllers"
	"github.com/windy40/lwnsimulator/models"
	dev "github.com/windy40/lwnsimulator/simulator/components/device"
	rp "github.com/windy40/lwnsimulator/simulator/components/device/regional_parameters"
	mrp "github.com/windy40/lwnsimulator/simulator/components/device/regional_parameters/models_rp"
	gw "github.com/windy40/lwnsimulator/simulator/components/gateway"
	"github.com/windy40/lwnsimulator/socket"
	_ "github.com/windy40/lwnsimulator/webserver/statik"
)

//WebServer type
type WebServer struct {
	Address      string
	Port         int
	Router       *gin.Engine
	ServerSocket *socketio.Server
}

var (
	simulatorController cnt.SimulatorController
	configuration       *models.ServerConfig
)

func NewWebServer(config *models.ServerConfig, controller cnt.SimulatorController) *WebServer {

	serverSocket := newServerSocket()

	configuration = config
	simulatorController = controller

	go func() {

		err := serverSocket.Serve()

		if err != nil {
			log.Fatal(err)
		}

	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = []string{"Origin", "Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers", "Content-type"}
	configCors.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))

	router.Use(gin.Recovery())

	ws := WebServer{
		Address:      configuration.Address,
		Port:         configuration.Port,
		Router:       router,
		ServerSocket: serverSocket,
	}

	staticFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticGroup := router.Group("/dashboard")
	staticGroup.StaticFS("/", staticFS)
	//router.Use(static.Serve("/", staticFS))

	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/start", startSimulator)
		apiRoutes.GET("/stop", stopSimulator)
		apiRoutes.GET("/bridge", getRemoteAddress)
		apiRoutes.GET("/gateways", getGateways)
		apiRoutes.GET("/devices", getDevices)
		apiRoutes.POST("/add-device", addDevice)
		apiRoutes.POST("/up-device", updateDevice)
		apiRoutes.POST("/del-device", deleteDevice)
		apiRoutes.POST("/del-gateway", deleteGateway)
		apiRoutes.POST("/add-gateway", addGateway)
		apiRoutes.POST("/up-gateway", updateGateway)
		apiRoutes.POST("/bridge/save", saveInfoBridge)

	}

	router.GET("/socket.io/*any", gin.WrapH(serverSocket))
	router.POST("/socket.io/*any", gin.WrapH(serverSocket))

	router.GET("/", func(context *gin.Context) { context.Redirect(http.StatusMovedPermanently, "/dashboard") })

	return &ws
}

func startSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.Run())
}

func stopSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.Stop())
}

func saveInfoBridge(c *gin.Context) {

	var ns models.AddressIP
	c.BindJSON(&ns)

	c.JSON(http.StatusOK, gin.H{"status": simulatorController.SaveBridgeAddress(ns)})
}

func getRemoteAddress(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.GetBridgeAddress())
}

func getGateways(c *gin.Context) {

	gws := simulatorController.GetGateways()
	c.JSON(http.StatusOK, gws)
}

func addGateway(c *gin.Context) {

	var g gw.Gateway
	c.BindJSON(&g)

	code, id, err := simulatorController.AddGateway(&g)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code, "id": id})

}

func updateGateway(c *gin.Context) {

	var g gw.Gateway
	c.BindJSON(&g)

	code, err := simulatorController.UpdateGateway(&g)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code})

}

func deleteGateway(c *gin.Context) {

	Identifier := struct {
		Id int `json:"id"`
	}{}

	c.BindJSON(&Identifier)

	c.JSON(http.StatusOK, gin.H{"status": simulatorController.DeleteGateway(Identifier.Id)})

}

func getDevices(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.GetDevices())
}

func addDevice(c *gin.Context) {

	var device dev.Device
	c.BindJSON(&device)

	code, id, err := simulatorController.AddDevice(&device)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code, "id": id})

}

func updateDevice(c *gin.Context) {

	var device dev.Device
	c.BindJSON(&device)

	code, err := simulatorController.UpdateDevice(&device)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code})

}

func deleteDevice(c *gin.Context) {

	Identifier := struct {
		Id int `json:"id"`
	}{}

	c.BindJSON(&Identifier)

	c.JSON(http.StatusOK, gin.H{"status": simulatorController.DeleteDevice(Identifier.Id)})
}

func newServerSocket() *socketio.Server {

	serverSocket := socketio.NewServer(nil)

	serverSocket.OnConnect("/", func(s socketio.Conn) error {
		// windy40 dev socket : distinguish connection from web interface and from device
		remote_hdr := s.RemoteHeader()
		v, ok := remote_hdr["User-Agent"]

		if ok && strings.Contains(v[0], "Mozilla") {

			log.Println("[WS]: Socket connected\n")

			s.SetContext("")
			simulatorController.AddWebSocket(&s)
		}

		return nil

	})

	serverSocket.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.Close()
	})

	serverSocket.OnEvent("/", socket.EventToggleStateDevice, func(s socketio.Conn, Id int) {
		simulatorController.ToggleStateDevice(Id)
	})

	serverSocket.OnEvent("/", socket.EventToggleStateGateway, func(s socketio.Conn, Id int) {
		simulatorController.ToggleStateGateway(Id)
	})

	serverSocket.OnEvent("/", socket.EventMacCommand, func(s socketio.Conn, data socket.MacCommand) {

		switch data.CID {
		case "DeviceTimeReq":
			simulatorController.SendMACCommand(lorawan.DeviceTimeReq, data)
		case "LinkCheckReq":
			simulatorController.SendMACCommand(lorawan.LinkCheckReq, data)
		case "PingSlotInfoReq":
			simulatorController.SendMACCommand(lorawan.PingSlotInfoReq, data)

		}

	})

	serverSocket.OnEvent("/", socket.EventChangePayload, func(s socketio.Conn, data socket.NewPayload) (string, bool) {
		return simulatorController.ChangePayload(data)
	})

	serverSocket.OnEvent("/", socket.EventSendUplink, func(s socketio.Conn, data socket.NewPayload) {
		simulatorController.SendUplink(data)
	})

	serverSocket.OnEvent("/", socket.EventGetParameters, func(s socketio.Conn, code int) mrp.Informations {
		return rp.GetInfo(code)
	})

	serverSocket.OnEvent("/", socket.EventChangeLocation, func(s socketio.Conn, info socket.NewLocation) bool {
		return simulatorController.ChangeLocation(info)
	})

	setupDevEventHandler(serverSocket)

	return serverSocket
}

func (ws *WebServer) Run() {

	log.Println("[WS]: Listen [", ws.Address+":"+strconv.Itoa(ws.Port), "]")

	err := ws.Router.Run(ws.Address + ":" + strconv.Itoa(ws.Port))
	if err != nil {
		log.Println("[WS] [ERROR]:", err.Error())
	}

}
