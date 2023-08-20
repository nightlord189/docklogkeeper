package docker

import (
	"fmt"
	"os"
)

func (a *Adapter) GetContainers() []ContainerInfo {
	entries, err := os.ReadDir(a.Config.LogDirectory)
	if err != nil {
		fmt.Println("error read directory", err)
		return nil
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}
	return nil
}
