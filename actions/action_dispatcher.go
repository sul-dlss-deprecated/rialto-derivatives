package actions

import (
	"log"

	"github.com/sul-dlss/rialto-derivatives/message"
	"github.com/sul-dlss/rialto-derivatives/runtime"
)

// DispatchMessage transforms a message into an action
func DispatchMessage(msg *message.Message, registry *runtime.Registry) Action {
	switch msg.Action {
	case "touch":
		return NewTouchAction(registry)
	}
	log.Panicf("Unknown action '%s'. Allowed actions: touch.", msg.Action)
	return nil
}
