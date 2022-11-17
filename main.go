package main

import (
	"content_autogen/config"
	"content_autogen/drivers"
	"content_autogen/worker"
)

func main() {
	config.InitConfig()
	drivers.InitializeDrivers()

	worker.Consume()
}
