package nspawn

import (
	"fmt"
	"github.com/guelfey/go.dbus"
	"strconv"
)

type NSpawnNameToUIDMap struct{}

func (n *NSpawnNameToUIDMap) MapContainerName(container_name string) (string, error) {
	var dconn *dbus.Conn
	var machine_path dbus.ObjectPath
	var err error
	var uid dbus.Variant

	dconn, err = dbus.SystemBusPrivate()
	if err != nil {
		return "", err
	}

	//SetUID program. So be explicit about user
	err = dconn.Auth([]dbus.Auth{dbus.AuthExternal("root"), dbus.AuthCookieSha1("root", "/root")})
	if err != nil {
		dconn.Close()
		return "", err
	}

	err = dconn.Hello()
	if err != nil {
		dconn.Close()
		return "", err
	}

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return "", err
	}
	manager := dconn.Object("org.freedesktop.machine1", "/org/freedesktop/machine1")
	manager.Call("org.freedesktop.machine1.Manager.GetMachine", 0, container_name).Store(&machine_path)
	if err != nil {
		fmt.Println("Unable to find machine:", err.Error())
		return "", err
	}

	machine := dconn.Object("org.freedesktop.machine1", machine_path)
	uid, err = machine.GetProperty("org.freedesktop.machine1.Machine.Leader")
	if err != nil {
		fmt.Println("Unable to find process in container", err.Error())
		return "", err
	}

	return strconv.Itoa(int(uid.Value().(uint32))), nil
}
