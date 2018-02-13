all:
		go build cluster-top-node.go config.go data.go
		go build cluster-top-router.go config.go data.go
		go build cluster-top.go config.go data.go
		go build cluster-top-local.go config.go data.go
