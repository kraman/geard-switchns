package main

import (
	"fmt"
	"github.com/kraman/geard-switchns/switchns"
	"github.com/kraman/geard-switchns/switchns/uid_map"
	"github.com/kraman/geard-switchns/switchns/uid_map/docker"
	// "github.com/kraman/geard-switchns/switchns/uid_map/nspawn"
	"os"
	"os/user"
	"strconv"
)

func usage() {
	usage := "Switch into container namespace and execute command.\n\nUsage:\n\t" + os.Args[0] + "<docker container name> <args>...\n\n"
	fmt.Printf(usage)
}

func main() {
	var container_name string
	var err error
	uid := os.Getuid()
	args := []string{"/bin/bash"}

	if uid == 0 {
		if len(os.Args) < 2 {
			usage()
			os.Exit(1)
		}
		container_name = os.Args[1]
		args = os.Args[2:]
	} else {
		user_obj, err := user.LookupId(strconv.Itoa(uid))
		if err != nil {
			fmt.Printf("Unable to find user %s", err)
			os.Exit(1)
		}
		container_name = user_obj.Username
	}

	//mappers := []uid_map.IDMapper{ new(docker.DockerNameToUUIDMap), new(nspawn.NSpawnNameToUIDMap) }
	mappers := []uid_map.IDMapper{new(docker.DockerNameToUIDMap)}

	for _, m := range mappers {
		container_name, err = m.MapContainerName(container_name)
		if err != nil {
			fmt.Printf("Unable to map container name %v", err)
			os.Exit(1)
		}
	}

	pid, err := strconv.Atoi(container_name)
	if err != nil {
		fmt.Printf("Unable to find container PID %v: %v", container_name, err)
		os.Exit(1)
	}
	switchns.JoinContainer("", pid, args, []string{})
}
