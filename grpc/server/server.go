package main

import (
	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/pkg/api/proto/pvz/v1/pvz/v1"

	"HOMEWORK-1/internal/middleware"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"
	gw "HOMEWORK-1/pkg/api/proto/pvz/v1/pvz/v1" 

	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = 50051
	httpPort = ":63342"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("ошибка файла конфигурации: %v", err)
	}
}

func connectDB() *pgxpool.Pool {
	dbPassword := viper.GetString("database.password")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbName := viper.GetString("database.dbname")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("ошибка url: %v", err)
	}
	config.MaxConns = 30

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("нет соединения: %v", err)
	}
	return pool
}

func main() {
	initConfig()
	pool := connectDB()
	defer pool.Close()
	tm := &transactor.TransactionManager{Pool: pool}
	repo := repository.NewRepository(tm)
	PVZ := module.NewModule(module.Deps{
		Repository: repo,
		Transactor: tm,
	})

	commands := cli.NewCLI(cli.Deps{Module: PVZ})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	log.Println("Start")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer( // grpc сервер (aka http.Serever)
		grpc.ChainUnaryInterceptor(middleware.Logging),
	)
	pvz.RegisterPVZServer(grpcServer, commands)

	//
	//reflection.Register(grpcServer) // Рефлексия! (Повзоляет получать описание rpc функционала нашего сервиса. Полезно для Postman)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	errGw := gw.RegisterPVZHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if errGw != nil {
		log.Fatalf("failed to RegisterTelephoneHandlerFromEndpoint: %v", errGw)
		return
	}

	go func() {
		gwServer := &http.Server{ // Создаем HTTP gateway сервер
			Addr:    httpPort,
			Handler: middleware.WithHTTPLoggingMiddleware(mux), // middleware
		}

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		errHttp := gwServer.ListenAndServe()
		if errHttp != nil {
			log.Println(errHttp)
			return
		}
	}()

	if err = grpcServer.Serve(lis); err != nil { // запускаем grpc сервер
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("Done")
}
