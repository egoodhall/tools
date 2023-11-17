//go:build unix && !darwin

package daemon

import "github.com/kardianos/service"

func options() service.KeyValue {
	return service.KeyValue{}
}
