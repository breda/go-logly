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
	addr string
	port string

	certFile string
	keyFile  string

	memoryStore bool
	fileStore   bool

	memoryIndex  bool
	binTreeIndex bool
)

func init() {
	// connection settings
	flag.StringVar(&addr, "addr", "localhost", "Network address to bind to")
	flag.StringVar(&port, "port", "3333", "Network port to bind to")

	// TLS
	flag.StringVar(&certFile, "cert", "./config/localhost.pem", "Certificate file used for TLS connections")
	flag.StringVar(&keyFile, "key", "./config/localhost-key.pem", "Certificate ket file used for TLS connections")

	// Flags
	flag.BoolVar(&memoryStore, "memory", false, "Choose in-memory storage system")
	flag.BoolVar(&fileStore, "file", false, "Choose file storage system (data.db file)")

	flag.BoolVar(&memoryIndex, "idx-memory", false, "Use an in-memory (hash) index")
	flag.BoolVar(&binTreeIndex, "idx-bintree", false, "Use a binary-tree based index")

	flag.Parse()
}

func main() {
	// Init Logly
	logger := logger.New("main")
	logly := getLogly(logger)

	// Start servers
	go startGRPCServer(logly)
	go startHttpServer(logly)

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

	logly.Logger.Info().Msg(fmt.Sprintf("started HTTP secure server on %s:%s", addr, "3332"))
	http.ListenAndServeTLS(net.JoinHostPort(addr, "3332"), certFile, keyFile, nil)
}

func getLogly(logger *zerolog.Logger) *logly.Logly {
	if memoryStore && fileStore {
		logger.Fatal().Str("error", "cannot choose two store systems at the same time").Send()
	}

	if memoryStore {
		logger.Info().Msg("using in-memory storage")
		return logly.InMemory()
	}

	if fileStore {
		logger.Info().Msg("using file storage (data.db)")

		if memoryIndex && binTreeIndex {
			logger.Fatal().Str("error", "cannot choose two indexing systems at the same time").Send()
		}

		return logly.File()
	}

	return nil
}
