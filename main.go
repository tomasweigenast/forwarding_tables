package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"tomasweigenast.com/forwarding_tables/packets"
)

func main() {
	defer func() {
		default_pubsub.close()
		default_logger.close()
	}()

	env := new_environment()
	data := load_json_data()

	load_environment(env, data)

	host1 := env.get_device("host1")
	err := host1.send(net.ParseIP("10.0.0.2"), packets.NewICMP(8, 0))
	if err != nil {
		panic(err)
	}

	for {
	}

	/*routerA := new_router("routerA")
	routerB := new_router("routerB")
	host1 := new_host("host1")
	host2 := new_host("host2")

	routerA.add_interface("ra_int1", "10.192.0.1/10")
	routerA.add_interface("ra_int2", "10.0.0.1/30")
	routerA.ftable.add("10.0.0.0/30", "10.0.0.1", "ra_int2")

	host1.add_interface("host1_int", "10.192.0.5/10")
	host1.ftable.add("0.0.0.0/0", "10.192.0.5", "host1_int")

	routerB.add_interface("rb_int1", "10.0.0.2/30")
	routerB.add_interface("rb_int2", "10.0.0.13/30")
	host2.add_interface("host2_int", "10.0.0.14/30")

	err := host1.send(net.ParseIP("10.0.0.2"), []byte("hola routerB"))
	if err != nil {
		default_logger.infof("unable to send packet: %s", err)
	}

	for {
	}*/

	// rl, err := readline.New("> ")
	// if err != nil {
	// 	panic(err)
	// }

	// defer rl.Close()

	// for {
	// 	line, err := rl.Readline()
	// 	if err != nil {
	// 		break
	// 	}

	// 	lineParts := strings.Split(line, " ")
	// 	senderName := lineParts[0]
	// 	destName := lineParts[1]
	// 	data := lineParts[2:]

	// 	switch senderName {

	// 	}
	// }
}

func load_json_data() JsonFile {
	data, err := os.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	jsonData := make(JsonFile, 0)
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func load_environment(env *environment, data JsonFile) {
	for _, device := range data {
		if device.DeviceType == DeviceRouter {

			router := new_router(device.DeviceName)
			for _, ftable_entry := range device.Table {
				router.ftable.add(ftable_entry.Destination, ftable_entry.NextHop, ftable_entry.Interface)
			}

			for interf_name, interf_ip := range device.Interfaces {
				router.add_interface(interf_name, interf_ip)
			}

			env.devices = append(env.devices, router)

		} else if device.DeviceType == DeviceHost {
			host := new_host(device.DeviceName)
			for _, ftable_entry := range device.Table {
				host.ftable.add(ftable_entry.Destination, ftable_entry.NextHop, ftable_entry.Interface)
			}

			for interf_name, interf_ip := range device.Interfaces {
				host.add_interface(interf_name, interf_ip)
			}

			env.devices = append(env.devices, host)
		} else {
			panic(fmt.Errorf("invalid device type: %s", device.DeviceType))
		}
	}

	default_logger.infof("loaded %d devices", len(env.devices))
	for _, device := range env.devices {
		default_logger.infof("	device %s: %s", device.name(), device.id())
	}
	default_logger.infof("-------------")
}

func (env *environment) get_device(name string) Device {
	for _, device := range env.devices {
		if device.name() == name {
			return device
		}
	}

	return nil
}

type environment struct {
	devices []Device
}

func new_environment() *environment {
	return &environment{devices: make([]Device, 0)}
}
