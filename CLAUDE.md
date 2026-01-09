# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
go build    # Build the binary (produces ./envare)
go test -v  # Run all tests
```

## What This Is

`envare` is a single-file Go CLI tool that safely parses `.env` files and outputs their contents in shell-safe format. It uses `godotenv` for parsing and `shellescape` for safe quoting.

## Key Behaviors

- **Default**: Reads `.env` from current directory, outputs all variables as `KEY='value'` pairs
- **Single arg**: `envare VARNAME` outputs just the quoted value (no key=)
- **Multiple args**: `envare A B C` outputs `A='...'` format for each
- **Capture mode** (`-c`): Forces `KEY=value` output format even for single argument
- **Alt file** (`-f`): Use a different .env file path
