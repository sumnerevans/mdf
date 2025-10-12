# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mdf (Mutt Display Filter) is a custom display filter for the Mutt email client written in Go. It provides two main components:

1. **Display Filter** (default mode): Reads email from stdin and outputs formatted/colorized email with shortened URLs
2. **Daemon** (daemon mode): HTTP service that manages URL shortening and serves redirect pages

## Architecture

The project consists of three main Go files:

- **mdf.go**: Entry point with CLI flags and mode selection
- **daemon.go**: HTTP server (`MuttDisplayFilterDaemon`) with URL shortening and redirect page generation
- **filter.go**: Email parsing and formatting (`RunFilter`) with colorization and URL minification

### URL Shortening Flow

1. Filter reads email, finds URLs via regex
2. If URL is longer than root URI + 6 chars, POSTs to daemon's `/new` endpoint
3. Daemon generates 6-character random ID and stores mapping in memory
4. Filter replaces original URL with shortened version
5. When user visits shortened URL, daemon serves redirect page with 5-second delay

## Build and Development

Run in filter mode (processes email from stdin):
```bash
go run . --root-uri "http://localhost:3719/"
```

Run in daemon mode:
```bash
go run . --daemon --port 3719 --root-uri "http://localhost:3719/"
```

Build the project (if needed):
```bash
go build
```

The project uses Nix flakes for reproducible builds. Enter development environment:
```bash
nix develop
```

Build with Nix:
```bash
nix build
```

## Key Features Implementation

- **Git diff highlighting**: Regexes in filter.go detect git diff patterns and colorize additions (green), deletions (red), and metadata (blue/bold)
- **Date normalization**: Parses Date headers and converts to local timezone using /etc/localtime
- **Email formatting**: Removes redundant `<mailto:...>` wrappers from email addresses
- **X-Mailer highlighting**: Displays X-Mailer headers in red (filter.go:115-116)
- **Redirect page security**: 5-second delay before navigation is enabled to allow URL inspection (daemon.go:76-94)

## Important Implementation Details

- The daemon stores URL mappings in memory (map[string]string) - they do not persist across restarts
- Random IDs are 6 characters from a 62-character alphabet (a-z, A-Z, 0-9)
- The daemon does not guarantee unique IDs but collision probability is low
- Color output is forced on in the filter (color.NoColor = false) for Mutt integration
- URL shortening only occurs if the original URL is longer than root URI + 6 characters
