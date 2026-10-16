# Facile Chat

Terminal chat application and self-hosted XMPP infrastructure for Facile Studio.

## Architecture

This repository is organized as a monorepo containing the CLI chat app, web registration portal, and server infrastructure:

```text
.
├── apps/
│   └── web/            SvelteKit web registration and client directory
├── cmd/                CLI commands and versioning
├── internal/
│   ├── config/         Local configuration storage (~/.config/chat/config.yml)
│   ├── event/          Internal message and presence events
│   ├── ui/             Charmbracelet Bubble Tea TUI
│   └── xmpp/           Mellium XMPP protocol engine
├── docker-compose.yml  Ejabberd, Coturn, and web portal service definitions
├── ejabberd.yml        XMPP server configuration
├── filet.yml           Code quality standard enforcement
├── install.sh          Suite installer bootstrap shim
├── main.go             Terminal client entry point
└── mise.toml           Toolchain and task definitions
```

## Chat CLI

The `chat` CLI provides an interactive terminal chat interface built with Bubble Tea and Lip Gloss.

### Installation

Install via the Facile Studio installer:

```bash
facile install chat
```

Or run the repository installer script:

```bash
curl -fsSL https://raw.githubusercontent.com/FacileStudio/Chat/main/install.sh | bash
```

### Usage

Run the terminal client:

```bash
chat
```

Pass credentials directly via flags or environment variables:

```bash
chat -j user@chat.example.com -p mypassword
```

Environment variables:

- `CHAT_JID`: Account JID
- `CHAT_PASSWORD`: Account password
- `CHAT_SERVER`: Optional server hostname override

### In-app commands

- `/msg <jid> <message>`: Send a direct message to a user
- `/add <jid> [name]`: Add a contact to your roster
- `/status <chat|away|dnd|offline> [msg]`: Update presence
- `/clear`: Clear active conversation history
- `/quit`: Exit the client
- `Tab`: Switch focus between contacts list and chat input

## Server deployment

The server stack deploys via Dokploy using Docker Compose.

### Configuration

Define these environment variables:

```env
DOMAIN=chat.example.com
TURN_SECRET=generate_with_openssl_rand_hex_32
ADMIN_PASSWORD=strong_admin_password_here
```

Generate the TURN secret with OpenSSL:

```bash
openssl rand -hex 32
```

### Quality gate

Run the quality check suite with filet and Go test runners:

```bash
sh scripts/check.sh
```
