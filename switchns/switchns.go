package switchns

import (
	"fmt"
	"github.com/kraman/libcontainer"
	"github.com/kraman/libcontainer/namespaces"
	"github.com/kraman/libcontainer/utils"
	"os"
	"path"
	"strconv"
)

func CreateContainer(containerName string, nsPid int, args []string, env []string) (*libcontainer.Container, error) {
	container := new(libcontainer.Container)
	container.ID = containerName
	container.NsPid = nsPid
	container.Command = &libcontainer.Command{args, env}
	container.Namespaces = []libcontainer.Namespace{
		libcontainer.CLONE_NEWNS,
		libcontainer.CLONE_NEWUTS,
		libcontainer.CLONE_NEWIPC,
		libcontainer.CLONE_NEWPID,
		libcontainer.CLONE_NEWNET,
	}
	container.Capabilities = []libcontainer.Capability{
		libcontainer.CAP_SETPCAP,
		libcontainer.CAP_SYS_MODULE,
		libcontainer.CAP_SYS_RAWIO,
		libcontainer.CAP_SYS_PACCT,
		libcontainer.CAP_SYS_ADMIN,
		libcontainer.CAP_SYS_NICE,
		libcontainer.CAP_SYS_RESOURCE,
		libcontainer.CAP_SYS_TIME,
		libcontainer.CAP_SYS_TTY_CONFIG,
		libcontainer.CAP_MKNOD,
		libcontainer.CAP_AUDIT_WRITE,
		libcontainer.CAP_AUDIT_CONTROL,
		libcontainer.CAP_MAC_OVERRIDE,
		libcontainer.CAP_MAC_ADMIN,
	}
	netns_path := path.Join("/proc", strconv.Itoa(nsPid), "ns", "net")
	f, err := os.Open(netns_path)
	if err != nil {
		return nil, err
	}
	container.NetNsFd = f.Fd()

	return container, nil
}

func JoinContainer(containerName string, nsPid int, args []string, env []string) error {
	container, err := CreateContainer(containerName, nsPid, args, env)
	if err != nil {
		return fmt.Errorf("error creating container %s", err)
	}

	pid, err := namespaces.ExecIn(container, container.Command)
	if err != nil {
		return fmt.Errorf("error exexin container %s", err)
	}
	exitcode, err := utils.WaitOnPid(pid)
	if err != nil {
		return fmt.Errorf("error waiting on child %s", err)
	}
	os.Exit(exitcode)
	return nil
}
