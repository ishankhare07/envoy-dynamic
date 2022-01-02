run-control-plane:
	go run main.go
run-envoy:
	envoy -c github.com/ishankhare07/envoy-dynamic-control-plane.yaml -l debug
