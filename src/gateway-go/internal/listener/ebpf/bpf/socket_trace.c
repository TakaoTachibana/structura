// +build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct socket_event {
	__u32 pid;
	char comm[16];
	__u64 timestamp_ns;
};

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} events SEC(".maps");

SEC("kprobe/tcp_v4_connect")
int kprobe__tcp_v4_connect(void *ctx) {
	struct socket_event *event;

	event = bpf_ringbuf_reserve(&events, sizeof(*event), 0);
	if (!event) {
		return 0;
	}

	event->pid = bpf_get_current_pid_tgid() >> 32;
	bpf_get_current_comm(&event->comm, sizeof(event->comm));
	event->timestamp_ns = bpf_ktime_get_ns();

	bpf_ringbuf_submit(event, 0);
	return 0;
}

char LICENSE[] SEC("license") = "GPL";

