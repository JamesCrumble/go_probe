prog-test:
	go test -v

pprof-heap:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/heap.prof

pprof-allocs:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/allocs.prof

pprof-goroutine:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/goroutine.prof

pprof-block:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/block.prof

pprof-threadcreate:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/threadcreate.prof

# use top in interactive mode 
pprof-inuse-space:
	go tool pprof -nodefraction=0 -inuse_space ./.prof/heap.prof

pprof-inuse-objects:
	go tool pprof -nodefraction=0 -inuse_objects ./.prof/heap.prof

run:
	go run .

build:
	go build main.go -o build