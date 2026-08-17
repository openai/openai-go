#!/usr/bin/env python3
"""Capture complete, configuration-independent Git review snapshots."""

from __future__ import annotations

import argparse
import hashlib
import json
import pathlib
import re
import shutil
import subprocess
import tempfile


CHUNK_BYTES = 65_536
MAX_READ_BYTES = 4096
SHA_PATTERN = re.compile(r"[0-9a-f]{40}")
DIFF_ARGUMENTS = (
    "--no-ext-diff",
    "--no-textconv",
    "--ignore-submodules=none",
    "--no-relative",
    "--no-renames",
    "--text",
)


def git_arguments(repository: pathlib.Path, *arguments: str) -> list[str]:
    return [
        "git",
        "--no-replace-objects",
        "-C",
        str(repository),
        "-c",
        "diff.relative=false",
        *arguments,
    ]


def git_text(repository: pathlib.Path, *arguments: str) -> str:
    result = subprocess.run(
        git_arguments(repository, *arguments),
        check=True,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    return result.stdout.strip()


def capture_command(
    repository: pathlib.Path,
    destination: pathlib.Path,
    maximum_bytes: int,
    *arguments: str,
) -> dict[str, int | str]:
    process = subprocess.Popen(
        git_arguments(repository, *arguments),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    assert process.stdout is not None
    digest = hashlib.sha256()
    total = 0

    try:
        with destination.open("wb") as output:
            while chunk := process.stdout.read(CHUNK_BYTES):
                total += len(chunk)
                if total > maximum_bytes:
                    raise ValueError(f"{destination.name} exceeds {maximum_bytes} bytes")
                output.write(chunk)
                digest.update(chunk)
        _, error = process.communicate()
    except BaseException:
        process.kill()
        process.wait()
        raise

    if process.returncode:
        raise subprocess.CalledProcessError(process.returncode, process.args, stderr=error)
    if destination.stat().st_size != total:
        raise ValueError(f"incomplete {destination.name}: byte count does not match")
    return {"bytes": total, "sha256": digest.hexdigest()}


def manifest_count(path: pathlib.Path, maximum_files: int) -> int:
    content = path.read_bytes()
    if not content:
        return 0
    fields = content.split(b"\0")
    if fields.pop() != b"" or len(fields) % 2:
        raise ValueError("incomplete or invalid NUL-delimited change manifest")
    count = len(fields) // 2
    if count > maximum_files:
        raise ValueError(f"change manifest exceeds {maximum_files} files")
    for index in range(0, len(fields), 2):
        status, name = fields[index], fields[index + 1]
        if len(status) != 1 or status not in b"ADMTUXB" or not name:
            raise ValueError("unexpected or incomplete change-manifest entry")
    return count


def capture(args: argparse.Namespace) -> None:
    for label, revision in (("base", args.base), ("head", args.head)):
        if not SHA_PATTERN.fullmatch(revision):
            raise ValueError(f"{label} must be a full lowercase 40-digit Git commit SHA")

    root = pathlib.Path(git_text(pathlib.Path(args.repo), "rev-parse", "--show-toplevel")).resolve()
    merge_base = git_text(root, "merge-base", args.base, args.head)
    if not SHA_PATTERN.fullmatch(merge_base):
        raise ValueError("the captured commits have no trustworthy full merge base")

    directory = pathlib.Path(
        tempfile.mkdtemp(prefix="openai-go-pr-review-", dir=args.output_root)
    ).resolve()
    try:
        revision_range = f"{args.base}...{args.head}"
        manifest_path = directory / "changes.nul"
        patch_path = directory / "changes.patch"
        metadata_path = directory / "snapshot.json"

        manifest = capture_command(
            root,
            manifest_path,
            args.max_manifest_bytes,
            "diff",
            *DIFF_ARGUMENTS,
            "--name-status",
            "-z",
            revision_range,
        )
        file_count = manifest_count(manifest_path, args.max_files)
        patch = capture_command(
            root,
            patch_path,
            args.max_patch_bytes,
            "diff",
            *DIFF_ARGUMENTS,
            revision_range,
        )
        metadata = {
            "repo_root": str(root),
            "base": args.base,
            "head": args.head,
            "merge_base": merge_base,
            "manifest": str(manifest_path),
            "manifest_bytes": manifest["bytes"],
            "manifest_sha256": manifest["sha256"],
            "patch": str(patch_path),
            "patch_bytes": patch["bytes"],
            "patch_sha256": patch["sha256"],
            "file_count": file_count,
            "metadata": str(metadata_path),
        }
        metadata_path.write_text(json.dumps(metadata, indent=2) + "\n", encoding="utf-8")
        print(
            json.dumps(
                {
                    "metadata": str(metadata_path),
                    "file_count": file_count,
                    "manifest_bytes": manifest["bytes"],
                    "patch_bytes": patch["bytes"],
                }
            )
        )
    except BaseException:
        shutil.rmtree(directory)
        raise


def read_chunk(args: argparse.Namespace) -> None:
    if not 1 <= args.limit <= MAX_READ_BYTES:
        raise ValueError(f"chunk size must be between 1 and {MAX_READ_BYTES} bytes")
    if args.offset < 0:
        raise ValueError("chunk offset cannot be negative")

    metadata_path = pathlib.Path(args.snapshot).resolve()
    metadata = json.loads(metadata_path.read_text(encoding="utf-8"))
    filename = "changes.nul" if args.artifact == "manifest" else "changes.patch"
    path = metadata_path.parent / filename
    expected = metadata[f"{args.artifact}_bytes"]
    if path.stat().st_size != expected:
        raise ValueError(f"{args.artifact} no longer matches its captured byte count")
    if args.offset > expected:
        raise ValueError("chunk offset exceeds captured artifact size")

    with path.open("rb") as artifact:
        artifact.seek(args.offset)
        content = artifact.read(args.limit)
    next_offset = args.offset + len(content)
    print(
        json.dumps(
            {
                "artifact": args.artifact,
                "offset": args.offset,
                "next_offset": next_offset,
                "total_bytes": expected,
                "eof": next_offset == expected,
                "text": content.decode("utf-8", "surrogateescape"),
            }
        )
    )


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)

    capture_parser = commands.add_parser("capture")
    capture_parser.add_argument("--repo", required=True)
    capture_parser.add_argument("--base", required=True)
    capture_parser.add_argument("--head", required=True)
    capture_parser.add_argument("--output-root", default=None)
    capture_parser.add_argument("--max-patch-bytes", type=int, default=32 * 1024 * 1024)
    capture_parser.add_argument("--max-manifest-bytes", type=int, default=4 * 1024 * 1024)
    capture_parser.add_argument("--max-files", type=int, default=10_000)
    capture_parser.set_defaults(handler=capture)

    read_parser = commands.add_parser("read")
    read_parser.add_argument("--snapshot", required=True)
    read_parser.add_argument("--artifact", choices=("manifest", "patch"), required=True)
    read_parser.add_argument("--offset", type=int, required=True)
    read_parser.add_argument("--limit", type=int, default=1024)
    read_parser.set_defaults(handler=read_chunk)

    try:
        args = parser.parse_args()
        args.handler(args)
    except (OSError, ValueError, json.JSONDecodeError, subprocess.CalledProcessError) as error:
        message = error.stderr if isinstance(error, subprocess.CalledProcessError) else str(error)
        if isinstance(message, bytes):
            message = message.decode("utf-8", "replace")
        parser.exit(1, f"snapshot capture failed: {message}\n")


if __name__ == "__main__":
    main()
