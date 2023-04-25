package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/breda/logly/internal/grpc"
	logglyHttp "github.com/breda/logly/internal/http"
	"github.com/breda/logly/internal/logger"
	"github.com/breda/logly/internal/logly"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/credentials"
)

var (
	addr    string
	port    string
	name    string
	cluster string

	certFile string
	keyFile  string

	memoryStore bool
	fileStore   bool

	memoryIndex  bool
	binTreeIndex bool
)

func init() {
	// connection settings
	flag.StringVar(&addr, "addr", "127.0.0.1", "Network address to bind to")
	flag.StringVar(&port, "port", "3333", "Network port to bind to")
	flag.StringVar(&name, "name", "default", "Name of the node")
	flag.StringVar(&cluster, "cluster", "", "Cluster node address")

	// TLS
	flag.StringVar(&certFile, "cert", "./config/localhost.pem", "Certificate file used for TLS connections")
	flag.StringVar(&keyFile, "key", "./config/localhost-key.pem", "Certificate ket file used for TLS connections")

	// Flags
	flag.BoolVar(&memoryStore, "memory", false, "Choose in-memory storage system")
	flag.BoolVar(&fileStore, "file", true, "Choose file storage system (data.db file)")

	flag.BoolVar(&memoryIndex, "idx-memory", false, "Use an in-memory (hash) index")
	flag.BoolVar(&binTreeIndex, "idx-bintree", true, "Use a binary-tree based index")

	flag.Parse()
}

func main() {
	// Init Logly
	logger := logger.New("main")
	logly := getLogly(logger)

	// Start servers
	go startGRPCServer(logly)
	//go startHttpServer(logly)

	// And block forver
	<-make(chan struct{})
}

func startGRPCServer(logly *logly.Logly) {
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		logly.Logger.Fatal().Err(err)
	}

	grpcServer := grpc.NewGrpcServer(logly, creds)

	grpcLn, err := net.Listen("tcp4", net.JoinHostPort(addr, port))
	if err != nil {
		logly.Logger.Fatal().Err(err)
	}

	logly.Logger.Info().Msg(fmt.Sprintf("started gRPC secure server on %s", net.JoinHostPort(addr, port)))
	grpcServer.Serve(grpcLn)
}

func startHttpServer(logly *logly.Logly) {
	httpServer := logglyHttp.New(logly)

	http.HandleFunc("/append", httpServer.HandleAppend)
	http.HandleFunc("/fetch", httpServer.HandleFetch)
	http.HandleFunc("/", httpServer.HandleIndex)

	logly.Logger.Info().Msg(fmt.Sprintf("started HTTP secure server on %s:%s", addr, port))
	http.ListenAndServeTLS(net.JoinHostPort(addr, port), certFile, keyFile, nil)
}

func getLogly(logger *zerolog.Logger) *logly.Logly {
	if memoryStore {
		return logly.InMemory(name, addr, port, cluster)
	}

	if fileStore {
		return logly.File(name, addr, port, cluster)
	}

	return nil
}
