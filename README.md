# termtrix-netscout

A small ARP network scanner written in Go. It sweeps a local subnet by sending
raw ARP requests and prints the IP → MAC address of every host that replies.

No external packet libraries — Ethernet and ARP frames are built by hand and
sent over a raw `AF_PACKET` socket.

## How it works

1. Reads the IP and MAC address of the local network interface.
2. Builds a list of candidate hosts from the target CIDR range.
3. Spawns 10 worker goroutines that broadcast an ARP request per host.
4. A receiver goroutine listens for ARP replies (opcode 2) and prints each
   sender's IP and MAC.

## Requirements

- **Linux** — uses `AF_PACKET` raw sockets, which are not available on macOS or Windows.
- **Go 1.25+**
- **root privileges** — raw sockets require `CAP_NET_RAW`.

## Usage

```bash
sudo make server
```

Or directly:

```bash
cd cmd
go build -o netscout .
sudo ./netscout
```

Example output:

```
192.168.225.1 --> 00 1a 2b 3c 4d 5e
192.168.225.18 --> a4 5e 60 c1 88 09
192.168.225.30 --> 3c 22 fb 7d 11 04
```

## Configuration

Two values are currently hardcoded and need to be edited to match your machine:

| Value | Location |
| --- | --- |
| Interface name (`ens33`) | [internals/interface.go](internals/interface.go), [cmd/main.go](cmd/main.go) |
| Target subnet (`192.168.225.25/24`) | [cmd/main.go](cmd/main.go) |

Find your interface name with `ip addr`.

## Project layout

```
cmd/main.go                   entry point: socket setup, workers, reply listener
internals/interface.go        local interface IP/MAC lookup
internals/network.go          Ethernet + ARP struct definitions
internals/ethernet_frame.go   serializes structs into a 42-byte ARP frame
```

## Note

Only scan networks you own or have permission to test.
