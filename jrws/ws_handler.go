package jrws

import (
	"fmt"
	"github.com/db-tech/JsonRpcWebsocketServer/models"
	"log"
)

func (wsServer *WebsocketServer) ServiceHandlerLogin(request models.Request, ws *ConcurrentWebsocket) (interface{}, error) {
	loginParams := models.LoginParams{}
	err := CreateParamsObject(request.Params, &loginParams)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("login user: %v", loginParams.Username))
	wsServer.AddUser(loginParams.Username, ws)
	return nil, nil
}
