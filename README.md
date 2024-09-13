# go-discord-bot

This Go module provides basic functionality for a Discord bot that provides useful capabilities in all of the Discord servers I actively assist with.

## Usage

Ensure you have Go installed first:

```bash
CGO_ENABLED=0 go build -v

# prints help
./go-discord-bot -h
```

## Features

- responds to `ping`/`pong` messages, but only in specified channels via the `-c` flag
- message logging into a `jsonl` file
- presence detection (server leave/join events) into a `jsonl` file
- writes scheduled events to a `json` file for subscribed servers, with periodic refreshing (defaults to checking every hour)

### Future features

May implement but no guarantees:

- checking if users have no roles (they skipped onboarding and should be nudged)
- add a couple `/` slash commands
  - add `/donate` command to be taken to a page to donate
- add feedback modal slash command to allow members to provide feedback/suggestions
