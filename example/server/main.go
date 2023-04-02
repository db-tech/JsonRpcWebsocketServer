package main

import (
	"github.com/db-tech/JsonRpcWebsocketServer/jrws"
	"github.com/db-tech/JsonRpcWebsocketServer/models"
	"time"
)

type ListItemParams struct {
	ListItem string `json:"listItem"`
}

func main() {
	List := []string{"test", "test2", "test3"}

	server := jrws.NewWebsocketServer("/list", 8080)
	server.AddHandler("addListItem", func(request models.Request, ws *jrws.ConcurrentWebsocket) (interface{}, error) {
		listItemParams := ListItemParams{}
		err := jrws.CreateParamsObject(request.Params, &listItemParams)
		if err != nil {
			return nil, err
		}
		List = append(List, listItemParams.ListItem)
		return List, nil
	})

	server.AddHandler("getList", func(request models.Request, ws *jrws.ConcurrentWebsocket) (interface{}, error) {
		return List, nil
	})

	server.AddHandler("status", func(request models.Request, ws *jrws.ConcurrentWebsocket) (interface{}, error) {
		return "OK", nil
	})

	//Send notification every 4 seconds
	go func() {
		for {
			err := server.WriteNotificationToAllMembers("status", "OK")
			if err != nil {
				panic(err)
			}
			time.Sleep(4 * time.Second)
		}
	}()

	server.StartListening()

}
