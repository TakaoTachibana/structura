# Structura

> A high-throughput, polyglot telemetry processing pipeline designed with a **Zero-Allocation Memory Strategy** across C# (.NET) and Go.

---

##  System Architecture

```text
┌───────────────────────────┐         ┌───────────────────────────┐         ┌───────────────────────────┐
│  simulation-py (Python)   │   UDP   │     gateway-go (Go)       │   IPC   │   backend-core (C# .NET)  │
│  - asyncio Multi-Agent    ├────────►│  - sync.Pool Receiver     ├────────►│  - ArrayPool<byte>        │
│  - Dynamic Telemetry      │  :8080  │  - Zero-Alloc Buffer Pool │ (Socket)│  - ref struct Evaluator   │
└───────────────────────────┘         └───────────────────────────┘         └───────────────────────────┘
```

---

##  Key Highlights

- **Zero-Allocation Core (C# .NET)**: Built with `ArrayPool<byte>` and stack-allocated `ref struct` (`DeconstructiveEvaluator`) to eliminate heap allocations and GC pauses.
- **High-Throughput Gateway (Go)**: Uses `sync.Pool` for UDP buffer recycling, handling high-concurrency packet streaming without memory churn.
- **Asynchronous Chaos Generator (Python)**: Multi-agent telemetry simulator running on `asyncio` to stream dynamic, non-linear system state payloads.

---

##  Performance & Verification

Verified using native Go benchmark tooling (`go test -bench`) and .NET GC allocation APIs.

| Component | Metric | Result | Goal |
|---|---|---|---|
| **`gateway-go`** | Memory Allocation | **`0 B/op`** (0 allocs/op) | Zero GC pressure on packet ingest |
| **`gateway-go`** | Latency | **`11.20 ns/op`** | Ultra-low latency throughput |
| **`backend-core`** | GC Heap Allocation | **`0 bytes`** (100k eval loops) | Deterministic microsecond execution |

---

##  Monorepo Layout

```text
src/
├── backend-core/     # C# (.NET 8) Zero-allocation evaluation engine
├── gateway-go/       # Go 1.22 High-concurrency UDP gateway
└── simulation-py/   # Python 3.12 Multi-agent telemetry simulator
```

---

##  Quick Start

### 1. Start the Go Gateway
```bash
cd src/gateway-go
go run main.go
```

### 2. Launch the Python Simulator
```bash
cd src/simulation-py
source venv/bin/activate
python main.py
```
