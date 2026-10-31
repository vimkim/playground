#!/usr/bin/env python3
"""Build and exercise the five examples on Linux without opening a GUI."""
import argparse
import fcntl
import os
from pathlib import Path
import pty
import re
import select
import struct
import subprocess
import sys
import tempfile
import termios
import time

ROOT = Path(__file__).resolve().parents[1]
ENV = dict(os.environ, TERM="xterm-256color", PYTHONDONTWRITEBYTECODE="1")


def build(command, directory):
    result = subprocess.run(command, cwd=ROOT / directory, capture_output=True, text=True, timeout=180)
    if result.returncode:
        raise RuntimeError(result.stdout + result.stderr)
    print("PASS build", directory, flush=True)


def concurrency(binary, rust=False):
    result = subprocess.run([str(binary)], capture_output=True, text=True, timeout=60)
    assert result.returncode == 0, result.stderr
    pairs = re.findall(r"Job ID: (\d+), Final Value: (\d+)", result.stdout)
    expected = {i: (11 if i == 3 else 10) for i in range(1, 14)}
    assert len(pairs) == 13 and {int(i): int(v) for i, v in pairs} == expected, result.stdout
    if rust:
        assert "All workers have completed." in result.stdout, result.stdout
    print("PASS concurrency", binary.name, "13 completed jobs", flush=True)


def game(command, directory):
    master, slave = pty.openpty()
    fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", 24, 80, 0, 0))
    process = subprocess.Popen(
        command, cwd=ROOT / directory, stdin=slave, stdout=slave, stderr=slave,
        env=ENV, start_new_session=True,
        preexec_fn=lambda: fcntl.ioctl(0, termios.TIOCSCTTY, 0),
    )
    os.close(slave)
    output = bytearray()
    started = time.monotonic()
    next_key = started + 0.6
    try:
        while time.monotonic() - started < 8:
            if select.select([master], [], [], 0.05)[0]:
                try:
                    chunk = os.read(master, 65536)
                except OSError:
                    chunk = b""
                output.extend(chunk)
            if process.poll() is not None:
                break
            if time.monotonic() >= next_key:
                os.write(master, b"q")
                next_key += 0.25
        assert process.poll() is not None, f"Game did not quit: {directory}"
        assert process.returncode == 0, output.decode(errors="replace")
        plain = re.sub(rb"\x1b\[[0-?]*[ -/]*[@-~]", b"", output)
        assert b"Score:" in plain and b"GameOver!" in re.sub(rb"\s+", b"", plain), repr(output[-3000:])
        print("PASS terminal startup and quit", directory, flush=True)
    finally:
        if process.poll() is None:
            process.kill()
            process.wait()
        os.close(master)


def check(directory, skip_build):
    directory.mkdir(parents=True, exist_ok=True)
    dodge_go = directory / "dodge-go"
    concurrency_go = directory / "concurrency-go"
    dodge_rust = directory / "dodge-rust" / "debug" / "dodge"
    concurrency_rust = directory / "concurrency-rust" / "debug" / "concurrent-job-demo-rust"
    if not skip_build:
        build(["go", "build", "-p", "2", "-mod=readonly", "-o", str(dodge_go), "."], "games/dodge/go")
        build(["go", "build", "-p", "2", "-mod=readonly", "-o", str(concurrency_go), "."], "concurrency/go")
        build(["cargo", "build", "--locked", "-j", "2", "--target-dir", str(directory / "dodge-rust")], "games/dodge/rust")
        build(["cargo", "build", "--locked", "-j", "2", "--target-dir", str(directory / "concurrency-rust")], "concurrency/rust")
    concurrency(concurrency_go)
    concurrency(concurrency_rust, rust=True)
    game([str(dodge_go)], "games/dodge/go")
    game([sys.executable, "main.py"], "games/dodge/python")
    game([str(dodge_rust)], "games/dodge/rust")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--build-dir", type=Path, help="Keep and reuse build outputs in this directory.")
    parser.add_argument("--skip-build", action="store_true", help="Use existing binaries in --build-dir.")
    args = parser.parse_args()
    if args.skip_build and args.build_dir is None:
        parser.error("--skip-build requires --build-dir")
    if args.build_dir:
        check(args.build_dir.resolve(), args.skip_build)
    else:
        with tempfile.TemporaryDirectory(prefix="playground-smoke-") as path:
            check(Path(path), args.skip_build)
