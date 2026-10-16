# XMPP (Ejabberd) + Coturn Stack for Dokploy

Production-ready, privacy-focused instant messaging and A/V calling server stack deployed via Dokploy.

This setup runs ejabberd (with XEP-0215 STUN/TURN Discovery for Jingle WebRTC calls) and Coturn (media relay) with zero hardcoded secrets or IP addresses in the codebase.

## 📁 Repository Structure

```
.
├── .env.example        # Reference environment variables
├── docker-compose.yml  # Container definitions and Traefik labels
├── ejabberd.yml        # Generic ejabberd config with env var substitution
└── README.md           # Documentation & operations guide
```


## 🚀 Prerequisites

A VPS with Dokploy installed.

A registered domain or subdomain (e.g. `chat.example.com`) whose DNS A record points to your Dokploy server IP.

Open ports on your host/firewall:

`80/tcp` & `443/tcp`: Handled by Dokploy / Traefik (Web admin, WebSockets, Let's Encrypt).

`5222/tcp`: XMPP Client-to-Server (c2s) STARTTLS.

`5269/tcp` (Optional): Server-to-Server federation.

`3478/tcp` & `3478/udp`: STUN/TURN discovery service.

`49152-49200/udp`: RTP media relay port range for audio/video calls.

##⚙️ Configuration (.env)

Do not commit .env into git. Define these variables in Dokploy's Environment tab:

```env
DOMAIN=chat.example.com
TURN_SECRET=generate_with_openssl_rand_hex_32
```


Generate a secure turn secret via CLI:

```bash
openssl rand -hex 32
```


## 📦 Deployment via Dokploy

In your Dokploy dashboard, create a new Compose service.

Link your Git repository and select the branch (e.g., `main`).

Navigate to the Environment tab of the service and paste your `DOMAIN` and `TURN_SECRET`.

Click Deploy. Dokploy will pull images, mount volumes, and generate the SSL certificates automatically via Traefik.

## 👤 User Management

To create initial accounts or administrative users, run ejabberdctl inside the running container (via Dokploy Web Terminal or SSH):

Register an Admin User

```
docker exec -it ejabberd ejabberdctl register admin chat.example.com StrongPasswordHere123!
```


Register Standard Users

```
docker exec -it ejabberd ejabberdctl register alice chat.example.com AlicePassword123!
docker exec -it ejabberd ejabberdctl register bob chat.example.com BobPassword123!
```


## 🧪 Testing Audio/Video Calls

Open an XMPP client that supports Jingle (XEP-0166 / XEP-0215):

Android: Conversations

iOS: Siskin IM

Desktop (Linux/Windows): Gajim or Dino

Log in using your JID: `alice@chat.example.com`

Add `bob@chat.example.com` to your contact list.

Start an audio/video call. The app queries ejabberd for STUN/TURN credentials automatically and negotiates media streaming through Coturn when necessary.

## 🔒 Security Notes

No Leaked Credentials: Coturn receives credentials at startup via CLI arguments, and ejabberd reads them from container environment variables ($(TURN_SECRET) and $(DOMAIN)).

Restricted UDP Range: Coturn is limited to ports 49152-49200 to prevent opening unnecessarily broad firewall windows.
