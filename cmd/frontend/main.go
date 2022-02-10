// SPDX-License-Identifier: MIT
//
// Copyright 2022 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"time"

	"bursavich.dev/dynamictls"
	"bursavich.dev/dynamictls/grpctls"
	"bursavich.dev/dynamictls/tlsprom"
	"bursavich.dev/graceful"
	"bursavich.dev/grpcprom"
	"bursavich.dev/httpprom"
	"bursavich.dev/zapr"
	"bursavich.dev/zapr/zaprprom"

	"bursavich.dev/demo/pkg/frontend"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	bepb "bursavich.dev/grpcprom/testdata/backend"
	fepb "bursavich.dev/grpcprom/testdata/frontend"
)

func main() { os.Exit(realMain()) }

func realMain() (code int) {
	flags := flag.NewFlagSet(path.Base(os.Args[0]), flag.ExitOnError)

	debugAddr := flags.String("debug-addr", ":9091", "Listen address for HTTP debug server.")
	grpcAddr := flags.String("grpc-addr", ":8443", "Listen address for gRPC server.")
	backendAddr := flags.String("backend-addr", "backend.example.svc.cluster.local", "Dial address of backend server.")

	grpcCertFile := flags.String("grpc-cert", "/etc/example/grpc/tls.crt", "TLS certificate file for gRPC clients and server.")
	grpcKeyFile := flags.String("grpc-key", "/etc/example/grpc/tls.key", "TLS key file for gRPC clients and server.")
	grpcCAFile := flags.String("grpc-ca", "/etc/example/grpc/ca.crt", "TLS certificate authority file for gRPC clients and server.")

	shutdownDelay := flags.Duration("shutdown-delay", 10*time.Second,
		"Delay before starting graceful shutdown after receiving a shutdown signal. "+
			"This allows loadbalancer updates before the server stops accepting new requests.",
	)
	shutdownGrace := flags.Duration("shutdown-grace", 30*time.Second,
		"Maximum time allowed for pending requests to complete before shutdown.",
	)

	zaprObserver := zaprprom.NewObserver()
	zaprOptions := zapr.AllOptions(zapr.WithObserver(zaprObserver))
	zapr.RegisterFlags(flags, zaprOptions...)

	flags.Parse(os.Args[1:])

	// Create logger with metrics.
	log, sink := zapr.NewLogger(zaprOptions...)
	errLog := zapr.NewStdErrorLogger(sink)
	defer sink.Flush()

	// Create Prometheus registry and register collectors.
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		zaprObserver,
	)

	// Create dynamic TLS gRPC credentials with metrics.
	grpcTLSObserver, err := tlsprom.NewObserver(
		tlsprom.WithGRPC(),
		tlsprom.WithLogger(log),
	)
	if err != nil {
		log.Error(err, "Failed to create TLS Prometheus metrics")
		return 1
	}
	registry.MustRegister(grpcTLSObserver)
	grpcTLSConfig, err := dynamictls.NewConfig(
		dynamictls.WithObserver(grpcTLSObserver),
		dynamictls.WithBase(&tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS13,
		}),
		dynamictls.WithCertificate(*grpcCertFile, *grpcKeyFile),
		dynamictls.WithRootCAs(*grpcCAFile),
		dynamictls.WithClientCAs(*grpcCAFile),
		dynamictls.WithHTTP2(),
		dynamictls.WithLogger(log),
	)
	if err != nil {
		log.Error(err, "Failed to create dynamic TLS config for gRPC credentials")
		return 1
	}
	defer grpcTLSConfig.Close()
	grpcCreds, err := grpctls.NewCredentials(grpcTLSConfig)
	if err != nil {
		log.Error(err, "Failed to create gRPC credentials")
		return 1
	}

	// Create gRPC client connections with credentials and metrics.
	grpcClientMetrics := grpcprom.NewClientMetrics()
	registry.MustRegister(grpcClientMetrics)
	backendConn, err := grpc.Dial(
		*backendAddr,
		grpc.WithTransportCredentials(grpcCreds),
		grpc.WithStatsHandler(grpcClientMetrics.StatsHandler()),
		grpc.WithStreamInterceptor(grpcClientMetrics.StreamInterceptor()),
		grpc.WithUnaryInterceptor(grpcClientMetrics.UnaryInterceptor()),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		log.Error(err, "Failed to create gRPC backend connection")
		return 1
	}

	// Create gRPC server with credentials and metrics.
	grpcServerMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcServerMetrics)
	grpcSrv := grpc.NewServer(
		grpc.Creds(grpcCreds),
		grpc.StatsHandler(grpcServerMetrics.StatsHandler()),
		grpc.StreamInterceptor(grpcServerMetrics.StreamInterceptor()),
		grpc.UnaryInterceptor(grpcServerMetrics.UnaryInterceptor()),
	)
	fepb.RegisterFrontendServer(grpcSrv, &frontend.Server{
		BackendClient: bepb.NewBackendClient(backendConn),
	})
	reflection.Register(grpcSrv)
	grpcServerMetrics.Init(grpcSrv, codes.OK)

	// Create HTTP server with metrics.
	httpMux := httpprom.NewServeMux(httpprom.WithCode())
	registry.MustRegister(httpMux.Collector())
	httpSrv := &http.Server{
		Handler:  httpMux,
		ErrorLog: errLog,
	}

	// Create graceful server.
	srv := graceful.DualServerConfig{
		InternalServer: graceful.FromHTTP(httpSrv),
		ExternalServer: graceful.FromGRPC(grpcSrv),
		ShutdownDelay:  *shutdownDelay,
		ShutdownGrace:  *shutdownGrace,
		Logger:         log,
	}

	// Register debug handlers and metrics.
	httpMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry: registry,
		ErrorLog: errLog,
	}))
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	httpMux.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpMux.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-srv.ShuttingDown():
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})

	// Listen and serve with graceful shutdown.
	if err := srv.ListenAndServe(context.Background(), *debugAddr, *grpcAddr); err != nil {
		log.Error(err, "Failed to gracefully listen and serve")
		return 1
	}
	return 0
}
