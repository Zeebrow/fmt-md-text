#!/usr/bin/env python3
import argparse
import subprocess
from functools import total_ordering
from typing import List
import sys

@total_ordering
class Branch:
    def __init__(self, name):
        self.name = name.strip()
        self.branch_type = self.name.split("/")[0]
        self.version = self.name.split("/")[1]
        self.major = int(self.version.split(".")[0])
        self.minor = int(self.version.split(".")[1])
        self.patch = int(self.version.split(".")[2])

    def create(self):
        curr_branches = [str(b) for b in get_release_branches()]
        if self.__repr__() in curr_branches:
            print(f"A branch named '{self.__repr__()} already exists!")
            sys.exit(1)
        msg = str(input(f"Enter commit message: "))
        inp = str(input(f"Create and commit new branch '{self.__repr__()}' with message '{msg}'? (y/n) "))
        if not inp.lower().startswith('y'):
            sys.exit(1)
        co = subprocess.run(f"git checkout -b {self.__repr__()}".split(), capture_output=True, encoding='utf-8')
        subprocess.run(f"git checkout master".split())
        if co.returncode == 128:
            print(f"stdout: {co.stdout}")
            print(f"stderr: {co.stderr}")
            print("ERROR: Cannot bump version if the branch already exists! Exiting.")
            sys.exit(1)
        elif co.returncode != 0:
            print(f"stdout: {co.stdout}")
            print(f"stderr: {co.stderr}")
            print(f"ERROR: Failed to checkout branch '{self.__repr__()}':")
            sys.exit(1)
        else:
            print(co.stdout)
        commit = subprocess.run(f"git commit -am".split() + [f"{msg}"], capture_output=True, encoding='utf-8')
        if commit.returncode != 0: 
            print(f"WARN: Failed to commit branch '{self.__repr__()}':")
            print(f"stdout: {commit.stdout}")
            print(f"stderr: {commit.stderr}")
            inp = str(input(f"proceed anyways? (y/n) "))
            if not inp.lower().startswith('y'):
                print("exiting")
                sys.exit(1)
        else:
            print(commit.stdout)
        push= subprocess.run(f"git push -u origin {self.__repr__()}".split(), capture_output=True, encoding='utf-8')
        if push.returncode != 0: 
            print(f"Failed to push branch '{self.__repr__()}' to origin:")
            print(f"stdout: {push.stdout}")
            print(f"stderr: {push.stderr}")
            inp = str(input(f"proceed anyways? (y/n) "))
            if not inp.lower().startswith('y'):
                print("exiting")
                sys.exit(1)
        else:
            print(push.stdout)


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


def get_current_branch():
    branch = subprocess.run("git branch --show-current".split(),
            capture_output=True, encoding='utf-8')
    print(branch.stdout.strip())

def get_release_branches():
    branches = []
    _b = subprocess.run("git branch --list".split(),
            capture_output=True, encoding='utf-8')
    result = _b.stdout.split()
    for b in result:
        if b.strip().startswith("release"):
            branches.append(Branch(b.strip()))
    return branches

def get_latest(branches: List[Branch]) -> Branch:
    return max(branches)

def get_next_release(rel_type: str):
    l = get_latest(get_release_branches())
    nb = Branch(str(l))
    if rel_type == 'major':
        nb.major = l.major + 1
        nb.minor = 0
        nb.patch = 0
    elif rel_type == 'minor':
        nb.major = l.major
        nb.minor = l.minor + 1
        nb.patch = 0
    elif rel_type == 'patch':
        nb.major = l.major
        nb.minor = l.minor
        nb.patch = l.patch + 1
    else:
        print(f"unknon type '{rel_type}' - must be one of major|minor|patch")
        sys.exit(1)
    return nb

        
def deboog():
    branches = get_release_branches()
    print(f"{branches[0].version} == {branches[1].version} : {branches[0] == branches[1]}")
    print(f"{branches[0].version} > {branches[1].version} : {branches[0] > branches[1]}")
    print(f"{branches[0].version} < {branches[1].version} : {branches[0] < branches[1]}")
    print(f"{branches[0].version} >= {branches[1].version} : {branches[0] >= branches[1]}")
    print(f"{branches[0].version} <= {branches[1].version} : {branches[0] <= branches[1]}")
    for b in branches:
        print(b)
        print(f"{b.major} {b.minor} {b.patch}")
        print(f"{b.major + 1} {b.minor} {b.patch}")
    print(branches.sort()) #?
    print(max(branches))


def main():
    parser = argparse.ArgumentParser()
    group = parser.add_mutually_exclusive_group()
    group.add_argument("-M", "--major", action='store_true', help="bump the major version of the latest release")
    group.add_argument("-m", "--minor", action='store_true', help="bump the minor version of the latest release")
    group.add_argument("-p", "--patch", action='store_true', help="bump the patch version of the latest release")
    group.add_argument("-l", "--latest-release", action='store_true', help="print the latest release and exit")
    group.add_argument("-n", "--next-release", action='store', help="print the next release's branch name and exit", choices=['major', 'minor', 'patch'])
    group.add_argument("-a", "--all-releases", action='store_true', help="print all branches tagged with 'release/*' and exit")
    group.add_argument("-d", "--debug", action='store_true', help="print misc. info")
    args = parser.parse_args()
    release_branches = get_release_branches()
    l = get_latest(release_branches)
    new_branch = Branch(l.name)
    if args.major:
        print('major')
        get_next_release('major')
        new_branch.major = l.major + 1
        new_branch.minor = 0
        new_branch.patch = 0
        print(f"{l} --> {new_branch}")
        new_branch.create()
    elif args.minor:
        print('minor')
        get_next_release('minor')
        new_branch.major = l.major
        new_branch.minor = l.minor + 1
        new_branch.patch = 0
        print(f"{l} --> {new_branch}")
        new_branch.create()
    elif args.patch:
        print('patch')
        get_next_release('patch')
        new_branch.major = l.major
        new_branch.minor = l.minor
        new_branch.patch = l.patch + 1
        print(f"{l} --> {new_branch}")
        new_branch.create()
    elif args.latest_release:
        print(get_latest(get_release_branches()))
        exit(0)
    elif args.all_releases:
        print(get_release_branches())
        exit(0)
    elif args.next_release:
        print(get_next_release(args.next_release))
        exit(0)
    elif args.debug:
        deboog()
        exit(0)
    else:
        parser.print_help()
main()
