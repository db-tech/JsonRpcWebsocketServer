package jrws

import (
	"fmt"
	"github.com/db-tech/JsonRpcWebsocketServer/models"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	SERVICE_METHOD_LOGIN     = "login"
	SERVICE_METHOD_GET_USERS = "getUsers"
)

func SetGlobalLogger(logWriter io.Writer) {
	log.SetOutput(logWriter)
}

func init() {
	log.SetOutput(io.Discard)
}

type RequestQueueElement struct {
	Request             models.Request
	Username            string
	ConcurrentWebsocket *ConcurrentWebsocket
}

type WebsocketServer struct {
	WsHandlers                 map[string]func(request models.Request, ws *ConcurrentWebsocket) (interface{}, error)
	userToWebsocketMutex       sync.RWMutex
	UserToWebsocket            map[string]*ConcurrentWebsocket
	Path                       string
	Port                       int
	RequestElementQueueChannel chan RequestQueueElement
}

func NewWebsocketServer(path string, port int) *WebsocketServer {

	websocketServer := &WebsocketServer{
		WsHandlers:                 make(map[string]func(request models.Request, ws *ConcurrentWebsocket) (interface{}, error)),
		UserToWebsocket:            make(map[string]*ConcurrentWebsocket),
		RequestElementQueueChannel: make(chan RequestQueueElement, 1000),
		Path:                       path,
		Port:                       port,
	}
	websocketServer.AddHandler(SERVICE_METHOD_LOGIN, websocketServer.ServiceHandlerLogin)
	return websocketServer
}

func (wsServer *WebsocketServer) AddHandler(method string, handler func(request models.Request, ws *ConcurrentWebsocket) (interface{}, error)) {
	wsServer.WsHandlers[method] = handler
}

func (wsServer *WebsocketServer) StartListening() {
	go wsServer.RunRequestQueue()
	http.HandleFunc(wsServer.Path, wsServer.serveWs)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", wsServer.Port), nil))
}

func WriteJson(ws *ConcurrentWebsocket, data interface{}) error {
	err := ws.WriteJSON(data)
	if err != nil {
		log.Printf("write json error: %v\n", err)
		return err
	}
	return nil
}

func WriteNotification(ws *ConcurrentWebsocket, method string, data any) error {
	notification := models.NewJsonRpcNotification(method)
	notification.Params = data
	return WriteJson(ws, notification)
}

func (wsServer *WebsocketServer) GetConcurrentWebsocket(username string) (ws *ConcurrentWebsocket, exists bool) {
	wsServer.userToWebsocketMutex.RLock()
	defer wsServer.userToWebsocketMutex.RUnlock()
	ws, exists = wsServer.UserToWebsocket[username]
	return
}

func (wsServer *WebsocketServer) IterateUsers(f func(username string, ws *ConcurrentWebsocket) error) error {
	wsServer.userToWebsocketMutex.RLock()
	defer wsServer.userToWebsocketMutex.RUnlock()
	for user, ws := range wsServer.UserToWebsocket {
		err := f(user, ws)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wsServer *WebsocketServer) WriteNotificationToAllMembers(method string, data any) error {
	return wsServer.IterateUsers(func(userId string, client *ConcurrentWebsocket) error {
		err := WriteNotification(client, method, data)
		if err != nil {
			return fmt.Errorf("write notification error: %v", err)
		}
		return nil
	})
}

func (wsServer *WebsocketServer) WriterNotificationToMember(member string, method string, data any) error {
	return WriteNotificationToUser(wsServer, member, method, data)
}

func (wsServer *WebsocketServer) AddUser(username string, ws *ConcurrentWebsocket) {
	wsServer.userToWebsocketMutex.Lock()
	defer wsServer.userToWebsocketMutex.Unlock()
	wsServer.UserToWebsocket[username] = ws
	for _, client := range wsServer.UserToWebsocket {
		users := make([]string, 0, len(wsServer.UserToWebsocket))
		for user := range wsServer.UserToWebsocket {
			users = append(users, user)
		}
		err := WriteNotification(client, SERVICE_METHOD_GET_USERS, users)
		if err != nil {
			log.Printf("write notification error: %v\n", err)
			return
		}
	}
}

func (wsServer *WebsocketServer) RemoveUser(username string) {
	wsServer.userToWebsocketMutex.Lock()
	defer wsServer.userToWebsocketMutex.Unlock()
	delete(wsServer.UserToWebsocket, username)
	for _, client := range wsServer.UserToWebsocket {
		users := make([]string, 0, len(wsServer.UserToWebsocket))
		for user := range wsServer.UserToWebsocket {
			users = append(users, user)
		}
		err := WriteNotification(client, SERVICE_METHOD_GET_USERS, users)
		if err != nil {
			log.Printf("write notification error: %v\n", err)
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wsServer *WebsocketServer) serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	cws := &ConcurrentWebsocket{
		Ws: ws,
	}
	defer cws.Ws.Close()

	request := models.Request{}
	currentUser := ""

	for {
		err = cws.ReadJSON(&request)
		if err != nil {
			log.Println("ws closed", err)
			if currentUser != "" {
				wsServer.RemoveUser(currentUser)
				//delete(self.UserToWebsocket, currentUser)
				log.Println("deleted user from map", currentUser)
				return
			}
			break
		}
		log.Printf("ws Request Method: %v\n", request.Method)
		if len(request.Method) == 0 {
			log.Fatal("empty method string")
		}
		if request.Method == SERVICE_METHOD_LOGIN {
			loginParams := models.LoginParams{}
			err := CreateParamsObject(request.Params, &loginParams)
			if err != nil {
				log.Printf("create params object error: %v\n", err)
				return
			}
			currentUser = loginParams.Username
		}
		requestElement := RequestQueueElement{
			Request:             request,
			Username:            currentUser,
			ConcurrentWebsocket: cws,
		}
		wsServer.RequestElementQueueChannel <- requestElement
	}
}

func (wsServer *WebsocketServer) RunRequestQueue() {
	for {
		select {
		case requestElement := <-wsServer.RequestElementQueueChannel:
			request := requestElement.Request
			if len(request.Method) == 0 {
				log.Fatal("empty method string")
			}
			if handler, ok := wsServer.WsHandlers[request.Method]; ok {
				log.Printf("WS Start: %v", request.Method)
				response, err := handler(request, requestElement.ConcurrentWebsocket)
				if err != nil {
					err := requestElement.ConcurrentWebsocket.WriteJSON(models.NewJsonRpcResponseError(request.Id, -1, err.Error()))
					if err != nil {
						log.Printf("write json error: %v\n", err)
						return
					}
				} else {
					if response != nil {
						jsonResponse := models.NewJsonRpcResponseOk(request.Id)
						jsonResponse.Result = response
						err := requestElement.ConcurrentWebsocket.WriteJSON(jsonResponse)
						if err != nil {
							log.Printf("write json error: %v\n", err)
							return
						}
					} else {
						err := requestElement.ConcurrentWebsocket.WriteJSON(models.NewJsonRpcResponseOk(request.Id))
						if err != nil {
							log.Printf("write json error: %v\n", err)
							return
						}
					}
				}
				log.Printf("WS End: %v", request.Method)
			} else {
				err := requestElement.ConcurrentWebsocket.WriteJSON(models.NewJsonRpcResponseError(request.Id, -1, fmt.Sprintf("method %v not found", request.Method)))
				if err != nil {
					log.Printf("write json error: %v\n", err)
					return
				}
			}
		}
	}
}
