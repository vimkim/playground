# Playground

Small terminal games and concurrency examples in Go, Python, and Rust.

## Run an example

Run each command from the directory in the table.

| Example | Directory | Command |
| --- | --- | --- |
| Dodge game in Go | [games/dodge/go](games/dodge/go) | `go run .` |
| Dodge game in Python | [games/dodge/python](games/dodge/python) | `python3 main.py` |
| Dodge game in Rust | [games/dodge/rust](games/dodge/rust) | `cargo run --locked` |
| Concurrent jobs in Go | [concurrency/go](concurrency/go) | `go run .` |
| Concurrent jobs in Rust | [concurrency/rust](concurrency/rust) | `cargo run --locked` |

The games need an interactive terminal. Move with the left and right arrow keys.
Press `q` or Escape to quit. The Go and Python games then wait for another key.

See the [Rust game screenshot](assets/dodge-rust.png).

Use Go 1.23.5 or newer, Rust with Cargo, and Python 3 with curses support.
The concurrency examples finish when every job reaches at least 10.

## Verify the examples

Run these commands from the repository root:

```sh
python3 scripts/verify-imports.py
python3 scripts/smoke.py
```

The import check compares each directory with its original Git tree and checks
that the original commits and branch tags are present.
The smoke check builds the compiled examples, checks both concurrency results,
and starts and quits all three games in headless terminals on Linux.
Build outputs go to a temporary directory.

## Original histories

This repository combines five former repositories owned by `vimkim`:

- `dodge-game-go`
- `dodge-game-py`
- `dodge-game-rs`
- `concurrent-job-demo-rust`
- `concurrent-jobs-demo-go`

The import keeps their original files and commit IDs. Each source branch tip has
a tag named `source/<repository>/<branch>`. These tags retain the original root
layout. [imports.json](imports.json) records the source commits, trees, and paths.

For example, inspect the original Go game history with:

```sh
git log source/dodge-game-go/chatgpt
```
