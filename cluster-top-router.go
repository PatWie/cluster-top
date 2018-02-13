package main

import (
	"github.com/patwie/cluster-top/messaging"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"sync"
)

// nice cluster struct
var clus Cluster

var allNodes map[string]Node

func main() {

	// load ports and ip-address
	cfg := LoadConfig()
	cfg.Print()

	allNodes = make(map[string]Node)
	var mutex = &sync.Mutex{}

	// receiving messages in extra thread
	go func() {
		SocketAddr := "tcp://*:" + cfg.Ports.Nodes
		log.Println("Now listening on", SocketAddr)
		node_socket, err := zmq4.NewSocket(zmq4.PULL)

		if err != nil {
			panic(err)
		}
		defer node_socket.Close()
		node_socket.Bind(SocketAddr)

		for {
			// read node information
			s, err := node_socket.RecvBytes(0)
			if err != nil {
				log.Println(err)
				continue
			}

			var node Node
			err = msgpack.Unmarshal(s, &node)
			if err != nil {
				log.Println(err)
				continue
			}

			mutex.Lock()
			if _, ok := allNodes[node.Name]; !ok {
				log.Printf("A new node \"%v\" connected\n", node.Name)
			}
			allNodes[node.Name] = node
			mutex.Unlock()
		}
	}()

	// outgoing messages (REQ-ROUTER)
	SocketAddr := "tcp://*:" + cfg.Ports.Clients
	log.Println("Router binds to", SocketAddr)
	router_socket, err := zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		panic(err)
	}
	defer router_socket.Close()
	router_socket.Bind(SocketAddr)

	for {

		// read request of client
		msg, err := messaging.ReceiveMultipartMessage(router_socket)
		if err != nil {
			panic(err)
		}

		mutex.Lock()
		// rebuild cluster struct from map
		clus := Cluster{}
		for _, n := range allNodes {
			clus.Nodes = append(clus.Nodes, n)
		}
		mutex.Unlock()

		// send cluster information to client
		body, err := msgpack.Marshal(&clus)
		msg.Body = body
		messaging.SendMultipartMessage(router_socket, &msg)

	}
}
