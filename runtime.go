package pkg

import (
	"log"

	"go.uber.org/automaxprocs/maxprocs"
)

func ConfigMaxprocs() {
	maxprocs.Set(maxprocs.Logger(log.Printf))
}
