package main

import (
	"flag"
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
	memoryStore bool
	fileStore   bool

	memoryIndex  bool
	binTreeIndex bool

	certFile string
	keyFile  string
)

func init() {
	flag.BoolVar(&memoryStore, "memory", false, "Choose in-memory storage system")
	flag.BoolVar(&fileStore, "file", false, "Choose file storage system (data.db file)")

	flag.BoolVar(&memoryIndex, "idx-memory", false, "Use an in-memory (hash) index")
	flag.BoolVar(&binTreeIndex, "idx-bintree", false, "Use a binary-tree based index")

	// TLS
	flag.StringVar(&certFile, "cert", "./config/localhost.pem", "Certificate file used for TLS connections")
	flag.StringVar(&keyFile, "key", "./config/localhost-key.pem", "Certificate ket file used for TLS connections")

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

	grpcLn, err := net.Listen("tcp", ":3332")
	if err != nil {
		logly.Logger.Fatal().Err(err)
	}

	logly.Logger.Info().Msg("started gRPC secure server on port 3332")
	grpcServer.Serve(grpcLn)
}

func startHttpServer(logly *logly.Logly) {
	httpServer := logglyHttp.New(logly)

	http.HandleFunc("/append", httpServer.HandleAppend)
	http.HandleFunc("/fetch", httpServer.HandleFetch)
	http.HandleFunc("/", httpServer.HandleIndex)

	logly.Logger.Info().Msg("started HTTP secure server on port 3333")
	http.ListenAndServeTLS(":3333", certFile, keyFile, nil)
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
