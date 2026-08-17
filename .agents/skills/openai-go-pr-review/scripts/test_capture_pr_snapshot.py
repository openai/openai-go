#!/usr/bin/env python3
"""Regression tests for adversarial Git review-snapshot configuration."""

from __future__ import annotations

import json
import os
import pathlib
import subprocess
import sys
import tempfile
import unittest


HELPER = pathlib.Path(__file__).with_name("capture_pr_snapshot.py")


class SnapshotCaptureTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory(prefix="openai-go-snapshot-test-")
        self.addCleanup(self.temporary.cleanup)
        self.root = pathlib.Path(self.temporary.name)
        self.repo = self.root / "repository"
        self.repo.mkdir()
        self.output = self.root / "captures"
        self.output.mkdir()
        self.git("init", "-q")
        self.git("config", "user.name", "Snapshot Regression")
        self.git("config", "user.email", "snapshot@example.invalid")
        self.write("nested/current.go", "package nested\nconst Current = 1\n")
        self.write("AGENTS.md", "trusted policy v1\n")
        self.git("add", ".")
        self.git("commit", "-qm", "base")
        self.base = self.git("rev-parse", "HEAD").stdout.strip()

    def git(self, *arguments: str) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            ["git", "-C", str(self.repo), *arguments],
            check=True,
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

    def write(self, path: str, text: str) -> None:
        target = self.repo / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(text, encoding="utf-8")

    def commit(self) -> str:
        self.git("add", "-A")
        self.git("commit", "-qm", "head")
        return self.git("rev-parse", "HEAD").stdout.strip()

    def capture(
        self,
        head: str,
        *,
        cwd: pathlib.Path | None = None,
        extra: tuple[str, ...] = (),
        check: bool = True,
    ) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            [
                sys.executable,
                str(HELPER),
                "capture",
                "--repo",
                str(cwd or self.repo),
                "--base",
                self.base,
                "--head",
                head,
                "--output-root",
                str(self.output),
                *extra,
            ],
            check=check,
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            cwd=cwd,
        )

    def snapshot(self, result: subprocess.CompletedProcess[str]) -> dict:
        summary = json.loads(result.stdout)
        return json.loads(pathlib.Path(summary["metadata"]).read_text())

    def entries(self, snapshot: dict) -> list[tuple[str, str]]:
        parts = pathlib.Path(snapshot["manifest"]).read_bytes().split(b"\0")
        self.assertEqual(parts.pop(), b"")
        return [
            (parts[index].decode(), os.fsdecode(parts[index + 1]))
            for index in range(0, len(parts), 2)
        ]

    def test_binary_attributes_cannot_hide_text_changes(self) -> None:
        self.write("nested/current.go", "package nested\nconst Current = 42\n")
        head = self.commit()
        (self.repo / ".git/info/attributes").write_text("*.go -diff\n")
        naive = self.git("diff", f"{self.base}...{head}").stdout
        self.assertIn("Binary files", naive)

        snapshot = self.snapshot(self.capture(head))
        patch = pathlib.Path(snapshot["patch"]).read_text()
        self.assertIn("+const Current = 42", patch)
        self.assertNotIn("Binary files", patch)

    def test_renames_include_deleted_and_added_paths(self) -> None:
        self.write("outside/moved.go", (self.repo / "nested/current.go").read_text())
        (self.repo / "nested/current.go").unlink()
        head = self.commit()
        self.git("config", "diff.renames", "true")

        entries = self.entries(self.snapshot(self.capture(head)))
        self.assertIn(("D", "nested/current.go"), entries)
        self.assertIn(("A", "outside/moved.go"), entries)

    def test_nested_invocation_ignores_relative_diff_configuration(self) -> None:
        self.write("nested/current.go", "package nested\nconst Current = 2\n")
        self.write("AGENTS.md", "trusted policy v2\n")
        head = self.commit()
        self.git("config", "diff.relative", "true")

        snapshot = self.snapshot(self.capture(head, cwd=self.repo / "nested"))
        entries = self.entries(snapshot)
        self.assertIn(("M", "AGENTS.md"), entries)
        self.assertIn(("M", "nested/current.go"), entries)
        self.assertEqual(snapshot["repo_root"], str(self.repo.resolve()))

    def test_submodule_changes_ignore_reviewer_configuration(self) -> None:
        self.git("update-index", "--add", "--cacheinfo", f"160000,{self.base},vendor/module")
        self.git("commit", "-qm", "add gitlink")
        head = self.git("rev-parse", "HEAD").stdout.strip()
        self.git("config", "diff.ignoreSubmodules", "all")

        snapshot = self.snapshot(self.capture(head))
        self.assertIn(("A", "vendor/module"), self.entries(snapshot))
        self.assertIn("Subproject commit", pathlib.Path(snapshot["patch"]).read_text())

    def test_replacement_refs_cannot_rewrite_pinned_commits(self) -> None:
        self.write("nested/current.go", "package nested\nconst Current = 99\n")
        head = self.commit()
        self.git("replace", head, self.base)
        replaced = subprocess.run(
            ["git", "-C", str(self.repo), "diff", f"{self.base}...{head}"],
            check=False,
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        self.assertNotIn("+const Current = 99", replaced.stdout)

        snapshot = self.snapshot(self.capture(head))
        self.assertIn("+const Current = 99", pathlib.Path(snapshot["patch"]).read_text())

    def test_durable_capture_can_be_consumed_completely_in_bounded_chunks(self) -> None:
        for index in range(8):
            self.write(f"generated/file-{index:02}.go", "package generated\n" + "x" * 500)
        head = self.commit()
        snapshot = self.snapshot(self.capture(head))

        patch = pathlib.Path(snapshot["patch"]).read_bytes()
        self.assertEqual(snapshot["patch_bytes"], len(patch))
        self.assertEqual(snapshot["manifest_bytes"], pathlib.Path(snapshot["manifest"]).stat().st_size)
        self.assertEqual(snapshot["file_count"], len(self.entries(snapshot)))

        consumed = bytearray()
        offset = 0
        while offset < len(patch):
            result = subprocess.run(
                [
                    sys.executable,
                    str(HELPER),
                    "read",
                    "--snapshot",
                    snapshot["metadata"],
                    "--artifact",
                    "patch",
                    "--offset",
                    str(offset),
                    "--limit",
                    "512",
                ],
                check=True,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            chunk = json.loads(result.stdout)
            consumed.extend(chunk["text"].encode("utf-8", "surrogateescape"))
            self.assertGreater(chunk["next_offset"], offset)
            offset = chunk["next_offset"]

        self.assertEqual(bytes(consumed), patch)
        self.assertTrue(chunk["eof"])

    def test_patch_size_limit_fails_closed_and_removes_partial_capture(self) -> None:
        self.write("nested/current.go", "package nested\n" + "x" * 1024)
        head = self.commit()
        result = self.capture(head, extra=("--max-patch-bytes", "64"), check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("exceeds", result.stderr)
        self.assertEqual(list(self.output.iterdir()), [])

    def test_file_count_limit_fails_closed_and_removes_partial_capture(self) -> None:
        self.write("another.go", "package another\n")
        self.write("nested/current.go", "package nested\nconst Current = 3\n")
        head = self.commit()
        result = self.capture(head, extra=("--max-files", "1"), check=False)

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("exceeds 1 files", result.stderr)
        self.assertEqual(list(self.output.iterdir()), [])


if __name__ == "__main__":
    unittest.main()
