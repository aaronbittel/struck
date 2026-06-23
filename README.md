# struck

A command-line argument parser for Go using struct tags. Inspired by
[go-flags](https://github.com/jessevdk/go-flags).

## Features

- Struct-tag based CLI definition
- Long (`--flag`) and short (`-f`) flags
- Positional arguments
- Automatic type parsing
- Built-in help generation
- `--help`, `-h` support

## Installation
```console
    go get github.com/aaronbittel/struck
```

## Examples

See [examples](./examples).

## Behavior

### Flags

Flags are defined using:
- `long:"name"` -> `--name`
- `short:"n"` -> `-n`

- Flags (except `bool`) consume the next argument as their value
- Boolean flags are `true` if present, `false` if absent

### Positional Arguments

- Fields without long or short tags are treated as positional arguments.
- They are assigned in declaration order.
- Optional override:
    - `arg:"name"`
    - This changes the display name in help output.

### Help Flag

If any of the following is provided: `--help`, `-help`, `-h`
Then:
- Help output is printed
- `Parse()` returns `HelpRequested` error

## Struct Tags Reference

| Tag     | Meaning                               |
|---------|---------------------------------------|
| `long`  | Long flag name (`--name`)             |
| `short` | Short flag name (`-n`)                |
| `help`  | Help text                             |
| `arg`   | Positional argument name override     |

## Design Notes

- Non-exported fields are ignored
- Flags and positional arguments are mutually exclusive. Conflicting definitions (`long`
  and `arg` tag) are invalid and cause a panic
