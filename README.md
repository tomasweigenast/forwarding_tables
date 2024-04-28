# Forwarding Tables

A simple Golang project that simulates (a part of) the OSI network layer by implementing packets and redirection with static forwarding tables. It also implements a basic implementation of the ICMP protocol.

### Read this

This project *does not* fully implement the OSI Network layer. It misses the ARP protocol, IPv6, and a full IPv4 datagram packet
format implementation. At this moment, you can test your static forwarding tables.

### Incoming
- [ ] Data Link layer
  - [ ] Network discovery
  - [ ] MAC support
  - [ ] The "Switch" device
- [ ] Full Network layer
  - [ ] Full IPv4 packet 
  - [ ] Full ICMP commands
  - [ ] ARP protocol
  - [ ] IPv6 packets
  - [ ] Dynamic forwarding tables
    - [ ] DHCP protocol
       
### Notes
This won't implement (at least at the moment of writing this) the physical layer. Now it tries to simulate "wires" between devices using the pub/sub (or event bus) pattern
allowing devices to subscribe and publish to topics (networks).
