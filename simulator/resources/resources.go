package resources

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type Resources struct {
	ExitGroup       sync.WaitGroup        `json:"-"`
	WebSocket       socketio.Conn         `json:"-"`
	LinkedDevSocket map[int]socketio.Conn `json:"-"`
	ConnDevSocket   map[string]int        `json:"-"`
	//	LinkedDevEUItoId map[string]int        `json:"-"`
}

func (r *Resources) AddWebSocket(WebSocket *socketio.Conn) {
	r.WebSocket = *WebSocket
}

func (r *Resources) AddDevSocket(devSocket *socketio.Conn, Id int) {
	r.LinkedDevSocket[Id] = *devSocket
	SId := (*devSocket).ID()
	r.ConnDevSocket[SId] = Id
}

func (r *Resources) DeleteDevSocket(SId string) {
	if _, ok := r.ConnDevSocket[SId]; ok {
		delete(r.LinkedDevSocket, r.ConnDevSocket[SId])
		delete(r.ConnDevSocket, SId)
	}

}
