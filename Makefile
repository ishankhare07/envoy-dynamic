run-control-plane:
	go run main.go
run-envoy:
	envoy -c envoy-dynamic-control-plane.yaml -l debug
start-grpcui:
	grpcui --plaintext localhost:18000
