# Message transfer studies

Final intent is to understand how go + zmq can quickly move messages between containers in a k8s or container environment.

Develop an "adapter" in go that can handle metadata modification + data transfer to another adapter over a configured interface.
Eventually add in:

- Improved profiling capabilities
  - Memory
  - CPU
  - Data transfer rates
  - Latencies
- Performance maximization + visualization
- Metadata modification of data packets
- Data signing with SSL certs
- Containerization study: img size, vulnerabilities, performance differences (if any)
  - Alpine
  - Ubuntu/Debian
  - RHEL/RockyLinux
- Version management via templates (ideally YAML)
  - Investigate runtime vs compile time performance of version definitions
- Refactor into library which can be injected into code for handling transport

## Demo applications

Provide some demo apps which can be run to produce & consume message packets on either side of the adapters

Currently this is baked into the code. To run as a sidecar will require a refactor
