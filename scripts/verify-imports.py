#!/usr/bin/env python3
"""Verify that imported files, commit histories, and source refs are intact."""
import hashlib
import json
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[1]


def git(*args):
    return subprocess.check_output(["git", "-C", str(ROOT), *args], text=True).strip()


def main():
    imports = json.loads((ROOT / "imports.json").read_text())
    commits = set()
    refs = 0
    for entry in imports:
        source = entry["commit"]
        assert git("rev-parse", source + "^{tree}") == entry["tree"], entry["repository"]
        assert git("rev-parse", "HEAD:" + entry["directory"]) == entry["tree"], entry["directory"]
        git("merge-base", "--is-ancestor", source, "HEAD")
        history = git("rev-list", source).splitlines()
        assert len(history) == entry["original_commit_count"], entry["repository"]
        commits.update(history)
        for original, preserved in entry["refs"].items():
            assert git("rev-parse", "refs/tags/" + preserved["tag"]) == preserved["object"], original
            refs += 1
        for asset in entry.get("assets", []):
            assert hashlib.sha256((ROOT / asset["path"]).read_bytes()).hexdigest() == asset["sha256"], asset["path"]
        assert not git("status", "--porcelain", "--untracked-files=all", "--", entry["directory"]), entry["directory"]
        print("PASS", entry["repository"], "->", entry["directory"])
    print(f"Verified {len(imports)} source trees, {len(commits)} original commits, and {refs} source refs.")


if __name__ == "__main__":
    main()
