# XMPP and Coturn stack for Dokploy

Instant messaging and audio/video calling server stack deployed via Dokploy.

The stack runs ejabberd and Coturn with zero hardcoded credentials or IP addresses. It aligns with modern XMPP compliance standards: in-band registration with anti-abuse throttling, 1:1 chat with OMEMO encryption, private members-only rooms (MUC), HTTP file sharing, mobile push notifications, and WebRTC STUN/TURN media relay.

## Repository structure

```text
.
├── .env.example
├── docker-compose.yml
├── ejabberd.yml
└── README.md
```

## Prerequisites

A virtual private server with Dokploy installed.

A registered domain or subdomain pointing to the server IP.

Required open ports on the host:

`80/tcp` and `443/tcp`: Handled by Dokploy and Traefik for web admin, WebSockets, HTTP upload, and ACME certificates.

`5222/tcp`: XMPP client-to-server traffic with STARTTLS.

`5269/tcp` (Optional): Server-to-server federation.

`3478/tcp` and `3478/udp`: STUN/TURN discovery service.

`49152-49200/udp`: RTP media relay port range for audio/video calls.

## Configuration

Define these variables in the Dokploy environment tab:

```env
DOMAIN=chat.example.com
TURN_SECRET=generate_with_openssl_rand_hex_32
ADMIN_PASSWORD=strong_admin_password_here
```

Generate a secure TURN secret:

```bash
openssl rand -hex 32
```

## Deployment via Dokploy

1. In Dokploy, create a new Compose service.
2. Link this Git repository and select the target branch.
3. Add `DOMAIN`, `TURN_SECRET`, and `ADMIN_PASSWORD` in the environment tab.
4. Deploy the service. Ejabberd automatically provisions the admin account `admin@<DOMAIN>` on initial boot.

## User registration and privacy

### In-band registration

Anyone can open an XMPP client (Conversations, Dino, Siskin, Gajim), choose "Create account", enter the server domain, and register immediately.

To prevent bot swarms:
- `registration_timeout` limits account creation to one per 5 minutes per IP address.
- `mod_fail2ban` blocks IP addresses that trigger repeated authentication failures.
- `registration_watchers` sends an instant message to `admin@<DOMAIN>` whenever a new account registers.

### Private team rooms

Rooms created under `conference.<DOMAIN>` use private, members-only defaults:
- `members_only: true`: Non-members cannot enter or view message history. Ejabberd rejects unauthorized join attempts with error 407 (registration-required).
- `public: false`: Rooms are unlisted in service discovery. Outsiders cannot discover room names.
- `public_list: false`: Room member lists are hidden from non-members.
- `mam: true`: Message history is archived for authorized room members.
- `persistent: true`: Rooms remain active when occupants disconnect.

To add a team member to a room, room moderators invite them directly from their client or add their JID to the room member list.

For maximum privacy, enable OMEMO encryption inside the room settings in your client. All room messages are then end-to-end encrypted with Signal Double Ratchet keys between authorized members.

## Recommended clients

Android: Conversations or Cheogram

iOS: Siskin IM or Monal

Desktop (Linux, macOS, Windows): Dino or Gajim

Web: Converse.js or Movim
