import aiohttp
import os
import json
from datetime import datetime
from typing import Tuple, Optional, List, Dict
from rich.console import Console
from rich.table import Table
from rich import box
import re

console = Console()


class VersionInfo:
    def __init__(self, VERSION: str, UPDATE_DATE: str, CHANGES: List[str]):
        self.version = VERSION
        self.update_date = UPDATE_DATE
        self.changes = CHANGES


async def fetch_versions_json(
    repo_owner: str, repo_name: str, proxy: Optional[str] = None
) -> List[VersionInfo]:
    """
    Fetch versions information from GitHub repository
    """
    url = f"https://raw.githubusercontent.com/{repo_owner}/{repo_name}/refs/heads/main/0g.json"

    headers = {
        "Accept": "application/json",
        "Cache-Control": "no-cache",
        "User-Agent": "Mozilla/5.0",
    }

    proxy_url = None
    if proxy:
        proxy_url = f"http://{proxy}"

    try:
        async with aiohttp.ClientSession() as session:
            async with session.get(url, headers=headers, proxy=proxy_url) as response:
                if response.status == 200:
                    data = await response.text()

                    # Clean the JSON string while preserving spaces in values
                    # First, temporarily replace spaces in values with a special marker
                    # Find all string values and replace their spaces with a marker
                    data = re.sub(
                        r'"([^"]*)"', lambda m: m.group(0).replace(" ", "§"), data
                    )

                    # Now clean the JSON structure
                    data = data.replace(" ", "").replace("\n", "").replace("\t", "")
                    data = data.replace(",]", "]").replace(",}", "}")
                    while ",," in data:
                        data = data.replace(",,", ",")

                    # Restore spaces in values
                    data = data.replace("§", " ")

                    json_data = json.loads(data)
                    return [VersionInfo(**version) for version in json_data]
                else:
                    print("❌ Failed to fetch versions from GitHub")
                    if response.status == 403:
                        print("ℹ️ GitHub API rate limit exceeded or access denied")
                    elif response.status == 404:
                        print("ℹ️ Version file not found")
                    print("\n💡 You can try using a proxy by adding it to main.py:")
                    print("   await check_version(VERSION, proxy='user:pass@ip:port')")
                    return []
    except aiohttp.ClientError as e:
        print(f"❌ Network error while fetching versions: {e}")
        print("\n💡 You can try using a proxy by adding it to main.py:")
        print("   await check_version(VERSION, proxy='user:pass@ip:port')")
        return []
    except json.JSONDecodeError as e:
        print(f"❌ Error parsing version data: {e}")
        print("ℹ️ The version file format might be incorrect")
        return []
    except Exception as e:
        print(f"❌ Unexpected error: {e}")
        return []


def format_version_changes(versions: List[VersionInfo]) -> None:
    """
    Format version changes into a nice table output
    """
    if not versions:
        return

    try:
        table = Table(
            show_header=True,
            box=box.DOUBLE,
            border_style="bright_cyan",
            pad_edge=False,
            width=85,
            highlight=True,
        )

        table.add_column("Version", style="cyan", justify="center")
        table.add_column("Update Date", style="magenta", justify="center")
        table.add_column("Changes", style="green")

        for i, version in enumerate(versions):
            changes_str = "\n".join(f"• {change}" for change in version.changes)
            table.add_row(
                f"✨ {version.version}", f"📅 {version.update_date}", changes_str
            )
            # Add separator after each row except the last one
            if i < len(versions) - 1:
                table.add_row("─" * 12, "─" * 21, "─" * 40, style="dim")

        print("📋 Available Updates:")
        console.print(table)
        print()
    except Exception as e:
        print(f"❌ Error displaying version information: {e}")


async def check_version(
    current_version: str,
    repo_owner: str = "0xStarLabs",
    repo_name: str = "VersionsControl",
    proxy: Optional[str] = None,
) -> bool:
    """
    Check if current version is up to date
    """
    try:
        print("🔍 Checking for updates...")

        versions = await fetch_versions_json(repo_owner, repo_name, proxy)
        if not versions:
            return True

        # Sort versions by version number
        versions.sort(key=lambda x: [int(n) for n in x.version.split(".")])
        latest_version = versions[-1]

        current_version_parts = [int(n) for n in current_version.split(".")]
        latest_version_parts = [int(n) for n in latest_version.version.split(".")]

        if current_version_parts < latest_version_parts:
            print(f"⚠️ New version available: {latest_version.version}")
            format_version_changes(versions)
            return False

        print(f"✅ You are running the latest version ({current_version})")
        return True
    except Exception as e:
        print(f"❌ Error checking version: {e}")
        return True  # Return True to allow the program to continue running
