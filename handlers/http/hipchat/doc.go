/*
Package hipchat allows for log messages to be sent to a hipchat room.

Example

NOTE: "/notification" is added to the host url automatically.

	package main

	import (
		"github.com/go-playground/log"
		"github.com/go-playground/log/handlers/http/hipchat"
	)

	func main() {

		// NOTE: ROOM TOKEN must have view permissions for room
		hc, err := hipchat.New(hipchat.APIv2, "https://api.hipchat.com/v2/room/{ROOM ID or NAME}", "application/json", "{ROOM TOKEN}")
		if err != nil {
			log.Error(err)
		}
		hc.SetFilenameDisplay(log.Llongfile)

		log.RegisterHandler(hc, log.WarnLevel, log.AlertLevel, log.PanicLevel)

		log.WithFields(log.F("error", "something went wrong")).StackTrace().Alert("ALERT!")
	}
*/
package hipchat
