package actions

import (
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// DispatchMessage transforms a message into an action
func DispatchMessage(msg *message.Message, registry *runtime.Registry) Action {
	switch msg.Action {
	case "touch":
		log.Printf("Running Touch action\n")

		return NewTouchAction(registry)
	case "rebuild":
		log.Printf("Running Rebuild action\n")
		return NewRebuildAction(registry)
	}
	log.Panicf("Unknown action '%s'. Allowed actions: touch, rebuild.", msg.Action)
	return nil
}
