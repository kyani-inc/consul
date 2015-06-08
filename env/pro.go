// +build production

package env

import (
	"log"
)

func debug(msg string) {
	log.Println(msg)
}
