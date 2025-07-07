#!/usr/bin/env python3

import os
import os.path
import sys
from datetime import datetime

from packaging.version import Version


def main():
    major = "major" in sys.argv
    minor = not major and "minor" in sys.argv
    metadata_file = os.path.abspath(
        os.path.join(os.getcwd(), "internal", "metadata.txt")
    )

    current_release = datetime(1970, 1, 1)
    current_version = Version("0.0.0")
    if os.path.exists(metadata_file):
        with open(metadata_file) as f:
            metadata = f.read()

            try:
                cv = Version(metadata)
                current_version = cv
            except Exception as exc:
                print("Corrupted file: ", metadata_file, " error ", exc)
            current_release = datetime.fromtimestamp(os.path.getmtime(metadata_file))

    if major:
        new_version = Version(f"{current_version.major + 1}.0.0")
        new_release = datetime.now()
    elif minor:
        new_version = Version(f"{current_version.major}.{current_version.minor + 1}.0")
        new_release = datetime.now()
    else:
        print("Version: ", current_version)
        print("Release: ", current_release)
        return

    with open(metadata_file, "w") as f:
        f.write(str(new_version))

    print("Version: ", current_version, " -> ", new_version)
    print("Release: ", current_release, " -> ", new_release)


if __name__ == "__main__":
    main()
