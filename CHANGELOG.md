# Changelog

All notable changes are documented here. Format derives from Keep a Changelog.

## [Unreleased]

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
