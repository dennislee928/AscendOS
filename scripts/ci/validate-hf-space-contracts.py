#!/usr/bin/env python3
from __future__ import annotations

import re
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
SPACE_ROOT = REPO_ROOT / "infra" / "hf-spaces"

STATIC_CONTRACTS = {
    "infra/hf-spaces/ascendos-web/README.md": {
        "package_dir": REPO_ROOT / "apps" / "web-qwik",
        "required_files": [
            "README.md",
            "src/routes/index.tsx",
            "src/components/AppShell.tsx",
        ],
        "readme_markers": [
            "web-qwik Phase 6 Scaffold",
            "route-per-module placeholders",
            "Shared UI primitives are expected from `packages/ui`.",
        ],
        "expected_app_file": "index.html",
    },
}


def fail(errors: list[str]) -> None:
    for error in errors:
        print(f"ERROR: {error}")
    raise SystemExit(1)


def parse_frontmatter(path: Path) -> tuple[dict[str, str], str]:
    lines = path.read_text(encoding="utf-8").splitlines()
    if not lines or lines[0].strip() != "---":
        raise ValueError("missing YAML frontmatter opener")

    try:
        end_index = next(index for index, line in enumerate(lines[1:], start=1) if line.strip() == "---")
    except StopIteration as exc:
        raise ValueError("missing YAML frontmatter closer") from exc

    frontmatter: dict[str, str] = {}
    for raw_line in lines[1:end_index]:
        line = raw_line.strip()
        if not line:
            continue
        if ":" not in line:
            raise ValueError(f"invalid frontmatter line: {raw_line!r}")
        key, value = line.split(":", 1)
        value = value.strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in {'"', "'"}:
            value = value[1:-1]
        frontmatter[key.strip()] = value

    body = "\n".join(lines[end_index + 1 :])
    return frontmatter, body


def require(condition: bool, message: str, errors: list[str]) -> None:
    if not condition:
        errors.append(message)


def check_docker_space(readme_path: Path, frontmatter: dict[str, str], body: str, errors: list[str]) -> None:
    rel_readme = readme_path.relative_to(REPO_ROOT).as_posix()
    dockerfile_match = re.search(r"^- Dockerfile: `([^`]+)`$", body, flags=re.MULTILINE)
    require(dockerfile_match is not None, f"{rel_readme}: missing Dockerfile deployment marker", errors)
    if dockerfile_match is None:
        return

    dockerfile_path = REPO_ROOT / dockerfile_match.group(1)
    require(dockerfile_path.exists(), f"{rel_readme}: Dockerfile not found at {dockerfile_match.group(1)}", errors)
    if not dockerfile_path.exists():
        return

    dockerfile_text = dockerfile_path.read_text(encoding="utf-8")
    app_port = frontmatter["app_port"]
    require(
        f"Container listens on `0.0.0.0:{app_port}`" in body,
        f"{rel_readme}: Space README does not match the declared app_port {app_port}",
        errors,
    )
    require(
        re.search(rf"^EXPOSE\s+{re.escape(app_port)}\s*$", dockerfile_text, flags=re.MULTILINE) is not None,
        f"{rel_readme}: app_port {app_port} does not match EXPOSE in {dockerfile_match.group(1)}",
        errors,
    )

    port_argument = re.search(r"--port[= ]+\"?([0-9]+)\"?", dockerfile_text)
    if port_argument is not None:
        require(
            port_argument.group(1) == app_port,
            f"{rel_readme}: runtime port {port_argument.group(1)} does not match frontmatter app_port {app_port}",
            errors,
        )

    bind_address = re.search(r"--host[= ]+\"?([0-9.]+)\"?", dockerfile_text)
    if bind_address is not None:
        require(
            bind_address.group(1) == "0.0.0.0",
            f"{rel_readme}: Dockerfile should bind to 0.0.0.0 when host is specified",
            errors,
        )


def check_static_space(readme_path: Path, frontmatter: dict[str, str], body: str, errors: list[str]) -> None:
    rel_readme = readme_path.relative_to(REPO_ROOT).as_posix()
    require("app_port" not in frontmatter, f"{rel_readme}: static Space must not define app_port", errors)
    require(frontmatter.get("app_file") == "index.html", f"{rel_readme}: static Space must declare app_file: index.html", errors)
    require(
        "Static deployment contract for `ascendos-web`." in body,
        f"{rel_readme}: static Space guidance does not describe the live deployment contract",
        errors,
    )
    require(
        "This Space is intentionally `sdk: static` and does not use Docker" in body,
        f"{rel_readme}: static Space guidance is missing the no-Docker contract",
        errors,
    )
    require(
        "Static entrypoint: `index.html`" in body,
        f"{rel_readme}: static Space guidance is missing the static entrypoint contract",
        errors,
    )
    require(
        "Publish static frontend build artifacts in the Space root" in body,
        f"{rel_readme}: static Space guidance is missing the build-artifact contract",
        errors,
    )

    contract = STATIC_CONTRACTS.get(rel_readme)
    require(contract is not None, f"{rel_readme}: no static package contract mapping is defined", errors)
    if contract is None:
        return

    package_dir: Path = contract["package_dir"]
    require(package_dir.exists(), f"{rel_readme}: static package directory not found at {package_dir.relative_to(REPO_ROOT)}", errors)
    if not package_dir.exists():
        return

    for relative_file in contract["required_files"]:
        target = package_dir / relative_file
        require(
            target.exists(),
            f"{rel_readme}: static contract file missing at {target.relative_to(REPO_ROOT)}",
            errors,
        )

    readme_text = (package_dir / "README.md").read_text(encoding="utf-8")
    for marker in contract["readme_markers"]:
        require(
            marker in readme_text,
            f"{rel_readme}: static package README missing marker {marker!r}",
            errors,
        )

    expected_app_file = contract.get("expected_app_file")
    if expected_app_file is not None:
        entrypoint = frontmatter.get("app_file")
        require(
            entrypoint == expected_app_file,
            f"{rel_readme}: app_file {entrypoint!r} does not match expected static entrypoint {expected_app_file!r}",
            errors,
        )


def main() -> None:
    errors: list[str] = []
    readmes = sorted(SPACE_ROOT.glob("*/README.md"))
    require(bool(readmes), "no Hugging Face Space READMEs found under infra/hf-spaces", errors)

    for readme_path in readmes:
        rel_readme = readme_path.relative_to(REPO_ROOT).as_posix()
        try:
            frontmatter, body = parse_frontmatter(readme_path)
        except ValueError as exc:
            errors.append(f"{rel_readme}: {exc}")
            continue

        stem = readme_path.parent.name
        require(frontmatter.get("title") == stem, f"{rel_readme}: title must match Space directory name {stem!r}", errors)
        require(frontmatter.get("pinned") == "true", f"{rel_readme}: pinned must remain true", errors)
        require(frontmatter.get("emoji"), f"{rel_readme}: emoji is required", errors)
        require(frontmatter.get("colorFrom"), f"{rel_readme}: colorFrom is required", errors)
        require(frontmatter.get("colorTo"), f"{rel_readme}: colorTo is required", errors)

        sdk = frontmatter.get("sdk")
        require(sdk in {"docker", "static"}, f"{rel_readme}: sdk must be docker or static", errors)
        if sdk == "docker":
            require("app_port" in frontmatter, f"{rel_readme}: docker Space must define app_port", errors)
            check_docker_space(readme_path, frontmatter, body, errors)
        elif sdk == "static":
            check_static_space(readme_path, frontmatter, body, errors)

    if errors:
        fail(errors)

    print(f"Validated {len(readmes)} Hugging Face Space contract(s).")


if __name__ == "__main__":
    main()
