package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf socketTrace bpf/socket_trace.c
