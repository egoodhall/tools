package daemon

import "github.com/kardianos/service"

func options() service.KeyValue {
	return service.KeyValue{
		"RunAtLoad": "true",
	}
}
