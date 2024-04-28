package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"tomasweigenast.com/forwarding_tables/packets"
)

func main() {
	// defer func() {
	// 	default_pubsub.close()
	// 	default_logger.close()
	// }()

	env := new_environment()
	data := load_yaml_data()

	load_environment(env, data)

	host1 := env.get_device("host1")
	//err := host1.send(net.ParseIP("10.0.0.2"), packets.NewICMP(8, 0))
	err := host1.send(net.ParseIP("10.128.0.2"), packets.NewICMP(8, 0))
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

func load_yaml_data() DataFile {
	data, err := os.ReadFile("input2.yaml")
	if err != nil {
		panic(err)
	}

	jsonData := DataFile{}
	err = yaml.Unmarshal(data, &jsonData)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func load_environment(env *environment, data DataFile) {
	for _, device := range data.Devices {
		if device.DeviceType == DeviceRouter {

			router := new_router(device.DeviceName)
			for _, ftable_entry := range device.Table {
				destination, next_hop_ip, next_hop_interface := parse_forwarding_table_entry(device, ftable_entry)
				router.ftable.add(destination, next_hop_ip, next_hop_interface)
			}

			for interf_name, interf_ip := range device.Interfaces {
				router.add_interface(interf_name, interf_ip)
			}

			env.devices = append(env.devices, router)

		} else if device.DeviceType == DeviceHost {
			host := new_host(device.DeviceName)
			for _, ftable_entry := range device.Table {
				destination, next_hop_ip, next_hop_interface := parse_forwarding_table_entry(device, ftable_entry)
				host.ftable.add(destination, next_hop_ip, next_hop_interface)
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

func parse_forwarding_table_entry(device DeviceDefinition, i string) (destination, next_hop_ip, next_hop_interface string) {
	parts := strings.Split(i, ">")
	if len(parts) != 2 {
		panic(fmt.Errorf("forwarding table entry has a wrong format: %s [device %q]", i, device.DeviceName))
	}

	destination = strings.TrimSpace(parts[0])
	next_hop := strings.Split(parts[1], ":")
	if len(next_hop) > 2 || len(next_hop) < 1 {
		panic(fmt.Errorf("forwarding table entry has a wrong format: %s [device %q]", i, device.DeviceName))
	}

	if len(next_hop) == 2 {
		next_hop_ip = next_hop[0]
		next_hop_interface = next_hop[1]
	} else {
		next_hop_ip = device.Interfaces[next_hop[0]]
		next_hop_interface = next_hop[0]
	}

	return destination, strings.TrimSpace(next_hop_ip), strings.TrimSpace(next_hop_interface)
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
