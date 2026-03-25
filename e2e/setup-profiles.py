#!/usr/bin/env python3
"""Create an isolated asql profiles.yaml for VHS/e2e runs."""

from pathlib import Path
import sys


def main() -> int:
    config_home = Path(sys.argv[1]) if len(sys.argv) > 1 else Path("/tmp/asql-e2e-config")
    prod_dsn = sys.argv[2] if len(sys.argv) > 2 else "/tmp/asql-e2e.db"
    staging_dsn = sys.argv[3] if len(sys.argv) > 3 else "/tmp/asql-e2e-staging.db"

    profile_dir = config_home / "asql"
    profile_dir.mkdir(parents=True, exist_ok=True)
    profile_path = profile_dir / "profiles.yaml"
    profile_path.write_text(
        f"- name: prod\n  dsn: {prod_dsn}\n- name: staging\n  dsn: {staging_dsn}\n",
        encoding="utf-8",
    )
    profile_path.chmod(0o600)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
