PACKAGER="github.com/patwie/cluster-top/compiletimeconst"

# add compile-time constants to avoid path issues to configuration files
include cluster-top.env
LDFLAGS="-X ${PACKAGER}.ServerIp=${cluster_top_server_ip} -X ${PACKAGER}.PortGather=${cluster_top_server_port_gather} -X ${PACKAGER}.PortDistribute=${cluster_top_server_port_distribute} -X ${PACKAGER}.Tick=${cluster_top_tick_ms}"

all:
		go build -ldflags ${LDFLAGS} cluster-top.go config.go data.go
		go build -ldflags ${LDFLAGS} cluster-top-server.go config.go data.go
		go build -ldflags ${LDFLAGS} cluster-top-node.go config.go data.go

# PKG_CONFIG_PATH=/graphics/opt/opt_Ubuntu16.04/libzmq/dist/lib/pkgconfig \
# go build -v --ldflags '-extldflags "-static"' -a cluster-top-node.go