package sourcery

import (
	"fmt"
)

// MessageHandler is registered through [Receive] to handle
// messages for a specific address.
type MessageHandler = func(...any) (any, error)

// Receive registers a message handler to an address. Panic
// ensues if the address is already taken.
//
// Receive may be called during runtime.
func Receive(address string, handler MessageHandler) {
	if _, ok := app.msgHandlers[address]; ok {
		err := fmt.Sprintf(
			"Message handler already exists for address: %s",
			address,
		)
		panic(err)
	}

	app.msgHandlers[address] = handler
}

// Resign unregisters the messsage handler at the specified
// address freeing it to be used for another handler.
//
// Resign may be called during runtime. However, there are
// very few scenarios were calling this function is
// justified since the message system was designed with
// permanent addressing in mind.
func Resign(address string) {
	delete(app.msgHandlers, address)
}

// Query sends a synchronous request to the addressed
// message handler and returns the hanlder's response.
// Panic ensues if no handler is registered at the address.
func Query(address string, args ...any) (any, error) {
	// TODO
	return nil, nil
}
