# clingy
Self-hosted, end-to-end encrypted (soon) messaging from your terminal

<!-- preview gif here -->

---

## Setup

### Prerequisites

- Go 1.21+
- Node.js 20+

### Build

```sh
# Build TUI client binary - 
./clingy-client/build/build.sh

# Build remote server
./clingy-server/build/build.sh

# Generate TLS certs (required before running the server)
cd clingy-server && openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -nodes -subj "/"
```

### Run the relay server
From a network reachable by clients
```sh
./clingy-server/clingy-server
```

The server currently runs on HTTP/3 (QUIC) and holds an SSE stream open per registered user for inbound message delivery.

### Start the TUI client
```sh
./clingy
```

---

## Use

1. Register
On launch, the **Intro** page reads your saved server config. If one exists, it auto-registers with the relay and drops you into chat. If not, you will need to register through the Config modal.

1. Add a user
Open the Contacts modal and add a contact

1. Select a conversation
Select a user in the contacts modal to open the messages

1. Send a message

---

## Architecture

<!-- diagram here -->

There are three distinct parts to the clingy architecture

**Relay server (`clingy-server`)** — Go HTTP server. Holds one long-lived SSE connection per registered user and routes messages between them. Stateless aside from the in-memory connection map; no message persistence. Currently the server operates off of HTTP/3 but there is work in progress for setting HTTP/2 with a cli flag.

**Client API daemon (`clingy-client/api`)** — Go process running locally alongside the TUI. Owns the connection to the relay, manages local config and contacts, and exposes a small REST API on localhost. Separated the API at this layer with the intent of setting up other connections to send messages through this API, e.g. MCP or new UI.

**TUI (`clingy-client/ui`)** — Built on [OpenTUI](https://github.com/sst/opentui). Talks only to the local API daemon over HTTP/SSE — never to the relay directly. 

### Message flow

1. Sender's TUI → Local API
2. Local API → Relay
3. Relay looks up recipient in its connection map → pushes SSE event on the recipient's open stream
4. Recipient's local API receives SSE event → forwards its own SSE to the TUI
5. TUI receives and handles message

## Config

Daemon-side config (managed by `services/Config.go`):

```json
{
  "username": "...",
  "serverAddr": "...",
  "uniqueId": "...",
  "contacts": [
    { "username": "...", "uniqueId": "...", "contacts": [] }
  ]
}
```

TUI-side env:

- `API_URL` — base URL of the local daemon (e.g. `http://localhost:8888/api`)

---

## Status

Early development. Known gaps: no message persistence, no auth beyond username claim, E2E encryption advertised on the intro screen but not yet implemented, server CLI flags are stubbed (`TODO`s in `clingy-server/main.go`).

