package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/Android12349/food_recomendation/profile_managment_service/config"
	server "github.com/Android12349/food_recomendation/profile_managment_service/internal/api/profile_management_api"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/profile_management_api"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.ProfileManagementAPI, cfg *config.Config) {
	go func() {
		if err := runGRPCServer(api, cfg); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(cfg); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.ProfileManagementAPI, cfg *config.Config) error {
	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	profile_management_api.RegisterProfileManagementServiceServer(s, &api)

	slog.Info("gRPC-server server listening on " + grpcAddr)
	return s.Serve(lis)
}

func runGatewayServer(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := profile_management_api.RegisterProfileManagementServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)

	httpAddr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	slog.Info("gRPC-Gateway server listening on " + httpAddr)
	return http.ListenAndServe(httpAddr, r)
}
