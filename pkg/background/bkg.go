package background

import (
	"log"
)

// Background
func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		fn()
	}()
}
