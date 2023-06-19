package internal

import "os"

func HostName() string {
	name, err := os.Hostname()
	if err != nil {

	}
	return name
}