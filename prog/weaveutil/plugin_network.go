package main

import (
	"fmt"

	"github.com/weaveworks/weave/common/docker"
)

// API 1.21 is the first version that supports docker network commandsk
const DOCKER_API_VERSION = "1.21"

func removeNetwork(args []string) error {
	if len(args) != 1 {
		cmdUsage("remove-network", "<network-name>")
	}
	networkName := args[0]

	d, err := docker.NewVersionedClientFromEnv()
	if err != nil {
		return err
	}

	err = d.RemoveNetwork(networkName)
	if _, ok := err.(*docker.NoSuchNetwork); !ok && err != nil {
		if info, err2 := d.NetworkInfo(networkName); err2 == nil {
			if len(info.Containers) > 0 {
				containers := ""
				for container := range info.Containers {
					containers += fmt.Sprintf("  %.12s ", container)
				}
				return fmt.Errorf(`WARNING: the following containers are still attached to network %q:
%s
Docker operations involving those containers may pause or fail
while Weave is not running`, networkName, containers)
			}
		}
		return fmt.Errorf("unable to remove network: %s", err)
	}
	return nil
}
