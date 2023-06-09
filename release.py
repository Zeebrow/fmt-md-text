#!/usr/bin/env python3
import argparse
import subprocess
from functools import total_ordering
from typing import List
import sys


"""
release.py

Outputs information about remote branches to help name things
"""


ORIGIN_NAME = 'origin'
VALID_BRANCH_TYPES = ['rc', 'release', 'feature', 'fix']


@total_ordering
class Branch:
    """
    describes a remote branch
    name must be of the form 'remotes/origin/{branch_type}/{version}'
    """
    def __init__(self, name):
        name = name.strip()
        if name.startswith('remotes'):
            _idx_start = 2
        elif name.startswith('origin'):
            _idx_start = 1
        elif name.split('/')[0] in VALID_BRANCH_TYPES:
            _idx_start = 0
        else:
            raise Exception(f"invalid branch '{name}'")

        self.name = name

        self.branch_type = name.split('/')[_idx_start]
        if self.branch_type not in VALID_BRANCH_TYPES:
            raise Exception(f"invalid branch type '{self.branch_type}'")
        self.version = self.name.split("/")[_idx_start + 1]
        self.major = int(self.version.split(".")[0])
        self.minor = int(self.version.split(".")[1])
        self.patch = int(self.version.split(".")[2])

    def version(self):
        return f"{self.major}.{self.minor}.{self.patch}"

    def bump_patch(self):
        self.patch = self.patch + 1

    def bump_minor(self):
        self.minor = self.minor + 1

    def bump_major(self):
        self.major = self.major + 1

    def __repr__(self):
        return f"{self.branch_type}/{self.major}.{self.minor}.{self.patch}"

    def __eq__(self, other):
        if not isinstance(other, type(self)): return NotImplemented
        return ((self.major == other.major) and (self.minor == other.minor) and (self.patch == other.patch))

    def __gt__(self, other):
        if not isinstance(other, type(self)): return NotImplemented
        if (self.major > other.major):
            return True
        if (self.major == other.major) and (self.minor > other.minor):
            return True
        if (self.major == other.major) and (self.minor == other.minor) and (self.patch > other.patch):
            return True
        return False

    def __lt__(self, other):
        if not isinstance(other, type(self)): return NotImplemented
        if (self.major < other.major):
            return True
        if (self.major == other.major) and (self.minor < other.minor):
            return True
        if (self.major == other.major) and (self.minor == other.minor) and (self.patch < other.patch):
            return True
        return False

    def __ge__(self, other):
        if not isinstance(other, type(self)): return NotImplemented
        if self.__eq__(other):
            return True
        else:
            return self.__gt__(other)

    def __le__(self, other):
        if not isinstance(other, type(self)): return NotImplemented
        if self.__eq__(other):
            return True
        else:
            return self.__lt__(other)

def get_remote_release_branches(branch_type: str):
    if branch_type not in VALID_BRANCH_TYPES:
        raise Exception(f"invalid branch type '{branch_type}'")
    branches = subprocess.run("git branch --remote".split(), capture_output=True, encoding='utf-8').stdout.split()
    rtn = []
    for b in branches:
        if b.startswith(f"origin/{branch_type}"):
            branch_num = b.split(f"{branch_type}/")[1]
            rtn.append(Branch(b))
    return rtn


def get_latest_released_version():
    return max(get_remote_release_branches("release"))


def get_current_release_candidate():
    """
    returns the version number for the what the current minor release candidate
    should be, based off of the remote's release branches
    """
    branches = get_remote_release_branches("rc")
    if len(branches) == 0:
        # you need to make a release candidate branch
        raise Exception("you need to make a release candidate branch")
    latest_release = max(branches)
    latest_release.bump_minor()
    return latest_release


def get_version(rc: bool) -> str:
    """
    Determines the correct version number string for merging a feature into a release
    candidate, or for merging a release candidate into a release.

    If rc is False, the latest release version string is returned
    If rc is True, the current release-candidate version string is returned
    If rc is True and there is no rc branch on the remote, an exception is raised (TODO: test this)

    Runs `git fetch --prune` and `git fetch --all`, should probably checkout the tree in a tmp dir.
    """
    subprocess.run("git fetch --prune".split())
    if rc:
        branch = max(get_remote_release_branches("rc"))
        return branch.version + "-rc"
    branch = max(get_remote_release_branches("release"))
    return branch.version

def print_release_instructions():
    return """\
> (hint: pipe this into drcat!)

# how to release fmt-md-text (`drcat`)

* Release branches, `release/X.Y.z`, contain new features
* Release candidate branches, `rc/X.Y.Z`, contain all of the features required for a release
* Features are created in feature branches, `feature/X.Y.Z`, and are merged into a release candidate branch as they are completed.

You will always have a release candidate branch to work with. You may need to make your own feature branch and push it to the remote, e.g.
```
git checkout -b "feature/some-feature"
# make your changes
git commit -am "feat: create breadcrumbs in release.py"
git push -u origin feature/some-feature
```
From here, CI will check that your code builds. If it does, then you or your PM can merge it into the release candidate branch.

Any new features you want to add will get committed to a feature/ branch. Do this for every feature required for a release. Once the release candidate branch's features have been tested and are working as intended, create a release branch and merge the release candidate branch with it.

## Bugs

Bugs? No, we don't write those here!
If you have to write buggy code, make sure the fixes get included in the next release candidate by publishing the fix in a feature branch, and merging it into the next rc.
If it's a super important bug fix (i.e. a security fix) that affects one or more past release versions, do this:
- identify the first release in which the bug occurs
- create a fix/ branch off of that release branch's HEAD
- commit the fix/ branch
- merge it with the old release branch
- repeat for each release branch (you might be able to apply a patch to make it easier)

## PROBLEMS:
- (as if fixing bugs isn't problem enough)
- Is the next release a major release, or a minor one?
    + the version number is decided by the name of the branch.
"""


        
def main():
    parser = argparse.ArgumentParser()
    group = parser.add_mutually_exclusive_group()
    group.add_argument("-v", "--get-version", action='store_true', help="print the version number of the latest release branch and exit")
    group.add_argument("-n", "--rc", action='store_true', help="returns the version for the current release candidate, e.g. '0.2.0-rc'. Raises if no release candidate branch exists on the remote.")
    group.add_argument("-r", "--latest-release", action='store_true', help="print the latest release and exit")
    group.add_argument("-l", "--list-releases", action='store_true', help="print all branches named 'release/*' and exit")
    group.add_argument("--release-instructions", action='store_true', help="print some 'wisdom' you thought of to help you remember how to do a release in 2024 and beyond")
    args = parser.parse_args()

    if args.latest_release:
        print(get_latest_released_version())
    elif args.list_releases:
        [print(br) for br in get_remote_release_branches("release")]
    elif args.get_version:
        print(get_version(False))
    elif args.rc:
        print(get_version(True))
    elif args.release_instructions:
        print(print_release_instructions())
    else:
        parser.print_help()
    sys.exit(0)
main()
