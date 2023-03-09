package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/breda/logly/internal/grpc"
	logglyHttp "github.com/breda/logly/internal/http"
	"github.com/breda/logly/internal/logly"
)

var (
	memoryStore bool
	fileStore   bool

	memoryIndex  bool
	binTreeIndex bool
)

func init() {
	flag.BoolVar(&memoryStore, "memory", false, "Choose in-memory storage system")
	flag.BoolVar(&fileStore, "file", false, "Choose file storage system (data.db file)")

	flag.BoolVar(&memoryIndex, "idx-memory", false, "Use an in-memory (hash) index")
	flag.BoolVar(&binTreeIndex, "idx-bintree", false, "Use a binary-tree based index")

	flag.Parse()
}

func main() {
	// Init loggly
	loggly := getLoggly()

	grpcServer := grpc.NewGrpcServer(loggly)
	grpcLn, err := net.Listen("tcp", ":3332")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started gRPC server on port 3332")
	go grpcServer.Serve(grpcLn)

	httpServer := logglyHttp.New(loggly)
	http.HandleFunc("/append", httpServer.HandleAppend)
	http.HandleFunc("/fetch", httpServer.HandleFetch)

	log.Println("started HTTP server on port 3333")
	go http.ListenAndServe(":3333", nil)

	// Block forever
	<-make(chan struct{})
}

func getLoggly() *logly.Logly {
	if memoryStore && fileStore {
		log.Fatal("Cannot choose two storage systems at the same time")
	}

	if memoryStore {
		log.Println("using in-memory storage")
		return logly.InMemory()
	}

	if fileStore {
		log.Println("using file storage (data.db)")

		if memoryIndex && binTreeIndex {
			log.Fatal("Cannot choose two index systems at the same time")
		}

		return logly.File()
	}

	return nil
}
