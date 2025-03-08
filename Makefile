prog-test:
	python test.py

pprof: prog-test
	python update_profiling.py

pprof-heap:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/heap.out

pprof-allocs:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/allocs.out

pprof-goroutine:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/goroutine.out

pprof-block:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/block.out

pprof-threadcreate:
	go tool pprof -http :9000 -edgefraction 0 -nodefraction 0 -nodecount 100000 ./.prof/threadcreate.out

# use top in interactive mode 
pprof-inuse-space:
	go tool pprof -nodefraction=0 -inuse_space ./.prof/heap.out

pprof-inuse-objects:
	go tool pprof -nodefraction=0 -inuse_objects ./.prof/heap.out

run:
	go run .

build:
	go build main.go -o build