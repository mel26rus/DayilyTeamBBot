package platform

import "github.com/kardianos/service"

func IsService() (bool, error) {
	return !service.Interactive(), nil
}
