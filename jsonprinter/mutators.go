package jsonprinter

import "github.com/silverstagtech/loggos/jsonmessage"

// Mutator is used to change a message before it gets sent into the printers buffer.
// The returned bool must indicate to the printer if it should carry on proccessing the
// log message. If your mutation doesn't work consider creating a new json message
// and replacing the data in the pointer and sending that instead.
type Mutator interface {
	Mutate(*jsonmessage.JSONMessage) bool
}
