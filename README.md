# Demo

This provides a demo of various `bursavich.dev` packages.

See the [frontend](cmd/frontend/main.go) command for a unified usage example, from which the
following sample metrics are pulled.

---

## [bursavich.dev/graceful](https://bursavich.dev/graceful)

Package graceful provides graceful shutdown for servers.

It provides simple integrations with HTTP and gRPC servers.

It supports a single server and coordinated split internal/external servers (e.g. health checks,
debugging, and metrics on the internal server and primary services on the external server).

When a shutdown signal is received it delays some time before starting the shutdown to give load
balancers a chance to remove the instance before shutting down. After this delay, it starts the
graceful shutdown where clients are told the server is going away, pending requests are allowed
to complete, and new requests are rejects. After a grace period, any incomplete requests are
forcibly cancelled and the server exits.

Repeated shutdown signals short-circuit the delay and grace periods.

---

## [bursavich.dev/httpprom](https://bursavich.dev/httpprom)

Package httpprom provides Prometheus instrumentation for HTTP servers.

### Sample HTTP Metrics

Optionally, labels may be added for HTTP code and method. In this case only code is added.

By default the path of the handler is used for its name, but this can be overridden.

```bash
# HELP http_server_requests_pending Number of HTTP server requests currently pending.
# TYPE http_server_requests_pending gauge
http_server_requests_pending{handler="/debug/pprof/"} 0
http_server_requests_pending{handler="/debug/pprof/cmdline"} 0
http_server_requests_pending{handler="/health/liveness"} 0
http_server_requests_pending{handler="/health/readiness"} 0
http_server_requests_pending{handler="/metrics"} 1
# HELP http_server_requests_total Total number of HTTP server requests completed.
# TYPE http_server_requests_total gauge
http_server_requests_total{code="200",handler="/debug/pprof/"} 2
http_server_requests_total{code="200",handler="/debug/pprof/cmdline"} 1
http_server_requests_total{code="200",handler="/health/liveness"} 23
http_server_requests_total{code="200",handler="/health/readiness"} 23
http_server_requests_total{code="200",handler="/metrics"} 2
```

---

## [bursavich.dev/dynamictls](https://bursavich.dev/dynamictls)

Package dynamictls watches the filesystem and updates TLS configuration when certificate changes
occur.

It provides simple integrations with HTTP/1.1, HTTP/2, gRPC, and Prometheus.

### Sample DynamicTLS Metrics

Optionally, metrics names may be prefixed with "grpc", "http", or a custom namespace.
In this example the gRPC option is used.

```bash
# HELP grpc_tls_config_certificate_verify_error Indicates if there was an error verifying the TLS configuration's certificates and expirations.
# TYPE grpc_tls_config_certificate_verify_error gauge
grpc_tls_config_certificate_verify_error 0
# HELP grpc_tls_config_earliest_certificate_expiration_time_seconds Earliest expiration time of the TLS configuration's certificates in seconds since the Unix epoch.
# TYPE grpc_tls_config_earliest_certificate_expiration_time_seconds gauge
grpc_tls_config_earliest_certificate_expiration_time_seconds 1.70750508e+09
# HELP grpc_tls_config_update_error Indicates if there was an error updating the TLS configuration.
# TYPE grpc_tls_config_update_error gauge
grpc_tls_config_update_error 0
```

---

## [bursavich.dev/grpcprom](https://bursavich.dev/grpcprom)

Package grpcprom provides Prometheus instrumentation for gRPC clients and servers.

### Sample gRPC Client Metrics

Optionally, each metric can be disabled and each histogram's buckets can be overridden. In this
example the defaults are used.

```bash
# HELP grpc_client_connections_open Number of gRPC client connections open.
# TYPE grpc_client_connections_open gauge
grpc_client_connections_open 1
# HELP grpc_client_connections_total Total number of gRPC client connections opened.
# TYPE grpc_client_connections_total counter
grpc_client_connections_total 1
# HELP grpc_client_latency_seconds Latency of gRPC client requests.
# TYPE grpc_client_latency_seconds histogram
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.001"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.0025"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.005"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.01"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.025"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.05"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.1"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.25"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="0.5"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="1"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="2.5"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="5"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="10"} 7
grpc_client_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary",le="+Inf"} 7
grpc_client_latency_seconds_sum{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 0.0023483840000000002
grpc_client_latency_seconds_count{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
# HELP grpc_client_recv_bytes_count Bytes received in gRPC client responses count.
# TYPE grpc_client_recv_bytes_count counter
grpc_client_recv_bytes_count{grpc_frame="Header",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
grpc_client_recv_bytes_count{grpc_frame="Payload",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
grpc_client_recv_bytes_count{grpc_frame="Trailer",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
# HELP grpc_client_recv_bytes_sum Bytes received in gRPC client responses sum.
# TYPE grpc_client_recv_bytes_sum counter
grpc_client_recv_bytes_sum{grpc_frame="Header",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 26
grpc_client_recv_bytes_sum{grpc_frame="Payload",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 35
grpc_client_recv_bytes_sum{grpc_frame="Trailer",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 36
# HELP grpc_client_requests_pending Number of gRPC client requests pending.
# TYPE grpc_client_requests_pending gauge
grpc_client_requests_pending{grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 0
# HELP grpc_client_requests_total Total number of gRPC client requests completed.
# TYPE grpc_client_requests_total counter
grpc_client_requests_total{grpc_code="OK",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
# HELP grpc_client_sent_bytes_count Bytes sent in gRPC client requests count.
# TYPE grpc_client_sent_bytes_count counter
grpc_client_sent_bytes_count{grpc_frame="Header",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
grpc_client_sent_bytes_count{grpc_frame="Payload",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 7
# HELP grpc_client_sent_bytes_sum Bytes sent in gRPC client requests sum.
# TYPE grpc_client_sent_bytes_sum counter
grpc_client_sent_bytes_sum{grpc_frame="Header",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 0
grpc_client_sent_bytes_sum{grpc_frame="Payload",grpc_method="Query",grpc_service="backend.Backend",grpc_type="Unary"} 35
```

### Sample gRPC Server Metrics

Optionally, each metric can be disabled and each histogram's buckets can be overridden. In this
example the defaults are used.

```bash
# HELP grpc_server_connections_open Number of gRPC server connections open.
# TYPE grpc_server_connections_open gauge
grpc_server_connections_open 1
# HELP grpc_server_connections_total Total number of gRPC server connections opened.
# TYPE grpc_server_connections_total counter
grpc_server_connections_total 2
# HELP grpc_server_latency_seconds Latency of gRPC server requests.
# TYPE grpc_server_latency_seconds histogram
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.001"} 6
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.0025"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.005"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.01"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.025"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.05"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.1"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.25"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="0.5"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="1"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="2.5"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="5"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="10"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary",le="+Inf"} 7
grpc_server_latency_seconds_sum{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0.004028765
grpc_server_latency_seconds_count{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.001"} 0
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.0025"} 0
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.005"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.01"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.025"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.05"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.1"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.25"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="0.5"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="1"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="2.5"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="5"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="10"} 1
grpc_server_latency_seconds_bucket{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream",le="+Inf"} 1
grpc_server_latency_seconds_sum{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0.003289488
grpc_server_latency_seconds_count{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 1
# HELP grpc_server_recv_bytes_count Bytes received in gRPC server requests count.
# TYPE grpc_server_recv_bytes_count counter
grpc_server_recv_bytes_count{grpc_frame="Header",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_recv_bytes_count{grpc_frame="Header",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 1
grpc_server_recv_bytes_count{grpc_frame="Payload",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_recv_bytes_count{grpc_frame="Payload",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 4
grpc_server_recv_bytes_count{grpc_frame="Trailer",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0
grpc_server_recv_bytes_count{grpc_frame="Trailer",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0
# HELP grpc_server_recv_bytes_sum Bytes received in gRPC server requests sum.
# TYPE grpc_server_recv_bytes_sum counter
grpc_server_recv_bytes_sum{grpc_frame="Header",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 68
grpc_server_recv_bytes_sum{grpc_frame="Header",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 97
grpc_server_recv_bytes_sum{grpc_frame="Payload",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 35
grpc_server_recv_bytes_sum{grpc_frame="Payload",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 87
grpc_server_recv_bytes_sum{grpc_frame="Trailer",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0
grpc_server_recv_bytes_sum{grpc_frame="Trailer",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0
# HELP grpc_server_requests_pending Number of gRPC server requests pending.
# TYPE grpc_server_requests_pending gauge
grpc_server_requests_pending{grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0
grpc_server_requests_pending{grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0
# HELP grpc_server_requests_total Total number of gRPC server requests completed.
# TYPE grpc_server_requests_total counter
grpc_server_requests_total{grpc_code="OK",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_requests_total{grpc_code="OK",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 1
# HELP grpc_server_sent_bytes_count Bytes sent in gRPC server responses count.
# TYPE grpc_server_sent_bytes_count counter
grpc_server_sent_bytes_count{grpc_frame="Header",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_sent_bytes_count{grpc_frame="Header",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 1
grpc_server_sent_bytes_count{grpc_frame="Payload",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_sent_bytes_count{grpc_frame="Payload",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 4
grpc_server_sent_bytes_count{grpc_frame="Trailer",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 7
grpc_server_sent_bytes_count{grpc_frame="Trailer",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 1
# HELP grpc_server_sent_bytes_sum Bytes sent in gRPC server responses sum.
# TYPE grpc_server_sent_bytes_sum counter
grpc_server_sent_bytes_sum{grpc_frame="Header",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0
grpc_server_sent_bytes_sum{grpc_frame="Header",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0
grpc_server_sent_bytes_sum{grpc_frame="Payload",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 35
grpc_server_sent_bytes_sum{grpc_frame="Payload",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 2189
grpc_server_sent_bytes_sum{grpc_frame="Trailer",grpc_method="Query",grpc_service="frontend.Frontend",grpc_type="Unary"} 0
grpc_server_sent_bytes_sum{grpc_frame="Trailer",grpc_method="ServerReflectionInfo",grpc_service="grpc.reflection.v1alpha.ServerReflection",grpc_type="BidiStream"} 0
```

---

## [bursavich.dev/zapr](https://bursavich.dev/zapr)

Package zapr provides a [logr.LogSink](https://pkg.go.dev/github.com/go-logr/logr#LogSink)
implementation using [zap](https://pkg.go.dev/go.uber.org/zap).

It includes optional flag registration, Prometheus instrumentation, and a standard library
[*log.Logger](https://pkg.go.dev/log#Logger) adapter.

### Sample Zapr Flags

These are all of the available flags, but they may be enabled individually.

```txt
-log-caller
    Log caller file and line. (default true)
-log-caller-format value
    Log caller format (e.g. "full" or "short"). (default short)
-log-caller-key string
    Log caller key. (default "caller")
-log-development
    Log with development-friendly defaults.
-log-duration-format value
    Log duration format (e.g. "millis", "nanos", "secs", or "string"). (default secs)
-log-error-key string
    Log error key. (default "error")
-log-format value
    Log format (e.g. "console" or "json"). (default json)
-log-function-key string
    Log function key.
-log-level int
    Log verbosity level.
-log-level-format value
    Log level format (e.g. "color", "lower", or "upper"). (default upper)
-log-level-key string
    Log level key. (default "level")
-log-line-ending string
    Log line ending. (default "\n")
-log-message-key string
    Log message key. (default "message")
-log-name string
    Log name.
-log-name-key string
    Log name key. (default "logger")
-log-sampler-first int
    Log every call up to this count per tick. (default 100)
-log-sampler-thereafter int
    Log only one of this many calls after reaching the first sample per tick. (default 100)
-log-sampler-tick duration
    Sample logs over this duration. (default 1s)
-log-stacktrace
    Log stacktrace on error.
-log-stacktrace-key string
    Log stacktrace key. (default "stacktrace")
-log-time-format value
    Log time format (e.g. "iso8601", "millis", "nanos", "rfc3339", or "secs"). (default iso8601)
-log-time-key string
    Log time key. (default "time")
```

### Sample Zapr Metrics

Named loggers are tracked individually.

``` bash
# HELP log_bytes_total Total bytes of encoded log lines.
# TYPE log_bytes_total counter
log_bytes_total{level="error",name=""} 0
log_bytes_total{level="info",name=""} 312
# HELP log_encoder_errors_total Total number of log entry encoding failures.
# TYPE log_encoder_errors_total counter
log_encoder_errors_total{name=""} 0
# HELP log_lines_total Total number of log lines.
# TYPE log_lines_total counter
log_lines_total{level="error",name=""} 0
log_lines_total{level="info",name=""} 2
```
