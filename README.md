# Websocket JsonRPC Server

This project is a simple WebSocket server that supports JSON-RPC requests. It allows users to connect to the server without authentication and interact with it using JSON-RPC methods.

## Disclaimer

This project was created for my own hobby/test project and is not intended to be used in production environments. The code may not be well-documented, thoroughly tested, or optimized. Use at your own risk.

## Features

- WebsocketServer struct that holds WebSocket connection handlers and user information
- AddHandler method to register JSON-RPC request handlers
- GetConcurrentWebsocket method to get a user's WebSocket connection
- WriteNotificationToAllMembers and WriterNotificationToMember methods to send notifications to all or a specific user
- AddUser and RemoveUser methods to add or remove users from the server
- StartListening method to start the WebSocket server and listen for requests

## Usage

Just checkout the main.go file in the example/server directory !

Once the server is running, users can connect to it using a WebSocket client and send JSON-RPC requests. The only predefined request handler is "login", which requires a "username" parameter in the request object. This handler must be called by clients before they can receive notifications.

To add custom request handlers, use the `AddHandler` method to register a function that takes a `models.Request` object and a `*ConcurrentWebsocket` object as arguments, and returns an interface{} and an error.

To send notifications to users, use the `WriteNotificationToAllMembers` or `WriterNotificationToMember` methods.

To add or remove users, use the `AddUser` and `RemoveUser` methods.



## Contributing

If you have any suggestions or bug reports, feel free to open an issue or pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
