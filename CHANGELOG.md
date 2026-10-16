# Changelog

All notable changes are documented here. Format derives from Keep a Changelog.

## [Unreleased]

## [v0.1.7] - 2026-10-16

### Fixed

- Render actual message content with a clean `🔒 ` icon when encrypted messages are decrypted, rather than hiding the text behind a placeholder badge.
- Automatically announce both modern and legacy (`eu.siacs.conversations.axolotl`) key bundles upon connecting, enabling Conversations and other clients to discover and encrypt for this client without manual setup.
- Automatic in-line OMEMO decryption for incoming stanzas without requiring commands.

## [v0.1.6] - 2026-10-16

### Added

- OMEMO cryptographic identity engine using pure Go (`github.com/meszmate/xmpp-go/crypto/omemo`, X25519, Ed25519).
- Persistent device ID generation stored in `~/.config/chat/device_id`.
- Automatic announcement of device IDs to PEP (`urn:xmpp:omemo:2:devices` and `eu.siacs.conversations.axolotl.devicelist`).
- Automatic publishing of X3DH public key bundles (`urn:xmpp:omemo:2:bundles`) to server PEP upon connection.
- Key matching in incoming `<encrypted>` stanzas to differentiate messages encrypted specifically for this client versus other devices (`🔒 [OMEMO Encrypted for other device]`).
- Commands `/omemo status` and `/omemo register` to view device fingerprint and re-publish bundles.

## [v0.1.5] - 2026-10-16

### Fixed

- Left-align conversation names and active indicator dot cleanly in the sidebar.
- Fix newline rendering inside Lipgloss header styles that caused contacts to be indented and wrapped prematurely.
- Use rune-safe truncation for contact names to prevent multi-byte UTF-8 slicing.

## [v0.1.4] - 2026-10-16

### Fixed

- Safely detect OMEMO encrypted messages and Conversations fallback text instead of displaying raw error URLs.
- Format OMEMO encrypted messages with a clean lock badge (`🔒 [OMEMO Encrypted Message]`).
- Prevent empty OMEMO encrypted stanzas from being dropped during live chat and history queries.

### Added

- `/omemo` command displaying OMEMO encryption status and guidance.
- `/omemo clear` command to reset stale OMEMO device lists on the server PEP node so other clients send plain text.

## [v0.1.3] - 2026-10-16

### Fixed

- Removed purple background from selected conversation in the sidebar list.
- Left-aligned conversation names in the sidebar with consistent indentation.
- Moved the green dot indicator to show the currently selected conversation during navigation.

### Added

- Support for fetching past conversation messages using XMPP MAM (Message Archive Management, XEP-0313).
- Automatic history retrieval upon connecting and selecting conversations.
- `/history` command to manually reload past conversation messages.
- Message deduplication and chronological ordering with date and time formatting.

## [v0.1.2] - 2026-10-16

### Fixed

- Made conversations list interactive with keyboard navigation (Up, Down, j, k, Enter) and mouse clicks.
- Display username in the purple header box and removed "FACILE CHAT" label.
- Expand main chat container fully to the right border of the terminal.
- Name conversation with oneself as "You (Notes)" with online presence instead of empty name.

## [v0.1.1] - 2026-10-16

### Fixed

- Correct horizontal and vertical container dimensions in TUI layout to prevent overflow past terminal borders.
- Fit viewport, header, sidebar, and input box precisely to terminal width and height.

### Added

- In-terminal login setup modal when credentials are not yet configured.
- Support for CHAT_USER and CHAT_PASSWORD environment variables alongside CHAT_JID and CHAT_PASS.
- Support for ~/.chat.yml and ~/.chat.yaml config file paths.

## [v0.1.0] - 2026-10-16

### Added

- Terminal XMPP chat client built with Charmbracelet Bubble Tea and Lip Gloss.
- Two-pane terminal interface with roster sidebar, status indicators, and scrollable conversation viewport.
- Mellium XMPP protocol engine with TLS negotiation, SASL authentication, and presence management.
- Slash command suite (/msg, /add, /status, /clear, /quit, /help).
- Facile monorepo architecture with apps/web, filet standard checks, and mise task runner.
- Facile Studio suite installer integration via install.sh and facile CLI.
