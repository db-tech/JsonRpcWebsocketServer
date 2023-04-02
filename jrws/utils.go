package jrws

import (
	"encoding/json"
	"fmt"
	"log"
)

func CreateParamsObject(params interface{}, targetObject interface{}) error {
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(paramBytes, targetObject)
	if err != nil {
		return err
	}
	return nil
}

func WriteNotificationToUser(wsServer *WebsocketServer, username string, method string, data any) error {
	memberWs, ok := wsServer.GetConcurrentWebsocket(username)
	if ok {
		log.Printf(fmt.Sprintf("found member  %v, WriteNotification", username))
		return WriteNotification(memberWs, method, data)
	}
	return fmt.Errorf("member with name %v is not registered", username)
}
