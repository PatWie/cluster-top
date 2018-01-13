package main

import (
	"github.com/patwie/cluster-top/compiletimeconst"
	"strconv"
	"time"
)

type Config struct {
	ServerIp             string        // ip of cluster-top-server
	ServerPortGather     string        // port of cluster-top-server, which nodes send to
	ServerPortDistribute string        // port of cluster-top-server, where clients subscribe to
	Tick                 time.Duration // tick between receiving data
	MaxDisplay           int           // top n processes per node (sorted by usage)
}

func CreateConfig() Config {

	c := Config{}
	c.ServerIp = compiletimeconst.ServerIp
	c.ServerPortGather = compiletimeconst.PortGather
	c.ServerPortDistribute = compiletimeconst.PortDistribute

	ms, _ := strconv.Atoi(compiletimeconst.Tick)
	c.Tick = time.Duration(ms) * time.Millisecond
	md, _ := strconv.Atoi(compiletimeconst.MaxDisplay)
	c.MaxDisplay = md
	return c
}
