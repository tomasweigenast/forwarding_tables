package main

import (
	"fmt"
	"net"
)

func main() {
	routerA := new_router("routerA")
	routerB := new_router("routerB")
	host1 := new_host("host1")
	host2 := new_host("host2")

	routerA.add_interface("int1", "10.192.0.1/10")
	routerA.add_interface("int2", "10.0.0.1/30")
	routerA.ftable.add("10.0.0.0/30", "10.0.0.1", "int2")

	host1.add_interface("host1_int", "10.192.0.5/10")
	host1.ftable.add("0.0.0.0/0", "10.192.0.5", "host1_int")

	routerB.add_interface("int1", "10.0.0.2/30")
	routerB.add_interface("int2", "10.0.0.13/30")
	host2.add_interface("host2_int", "10.0.0.14/30")

	err := host1.send(net.ParseIP("10.0.0.2"), []byte("hola routerB"))
	if err != nil {
		fmt.Printf("unable to send packet: %s", err)
	}

	for {
	}

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
