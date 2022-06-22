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

    def version(self):
        return f"{self.major}.{self.minor}.{self.patch}"

    def create(self):
        def precheck():
            curr_branches = [str(b) for b in get_release_branches()]
            if self.__repr__() in curr_branches:
                print(f"A branch named '{self.__repr__()} already exists!")
                sys.exit(1)
            is_master = subprocess.run(f"git branch --show-current".split(), capture_output=True, encoding='utf-8')
            if not 'master' in is_master.stdout:
                print("You must run this script from the master branch")
                print(f"stdout: {curr_branch.stdout}")
                print(f"stderr: {curr_branch.stderr}")
                sys.exit(1)
            curr_branch = subprocess.run(f"git status --porcelain".split(), capture_output=True, encoding='utf-8')
            if curr_branch.stdout != "":
                print("You have uncommitted changes:")
                print(f"stdout: {curr_branch.stdout}")
                print(f"stderr: {curr_branch.stderr}")
                sys.exit(1)
                
        precheck()
        sp_glo = subprocess.run("git log -1 --pretty=%s".split(), capture_output=True, encoding='utf-8')
        msg = sp_glo.stdout.strip()
        inp = str(input(f"Create and push new branch '{self.__repr__()}' with message '{msg}'? (y/n) "))
        if not inp.lower().startswith('y'):
            sys.exit(1)

        co = subprocess.run(f"git checkout -b {self.__repr__()}".split(), capture_output=True, encoding='utf-8')
        if co.returncode != 0:
            print(f"stdout: {co.stdout}")
            print(f"stderr: {co.stderr}")
            print(f"ERROR: Failed to checkout branch '{self.__repr__()}':")
            sys.exit(1)
        else:
            print(co.stdout)

        push = subprocess.run(f"git push -u origin {self.__repr__()}".split(), capture_output=True, encoding='utf-8')
        if push.returncode != 0: 
            print(f"WARN: Failed to push branch '{self.__repr__()}' to origin:")
            print(f"stdout: {push.stdout}")
            print(f"stderr: {push.stderr}")
            inp = str(input(f"proceed anyways? (y/n) "))
            if not inp.lower().startswith('y'):
                print("exiting")
                sys.exit(1)
        else:
            print(push.stdout)
        subprocess.run(f"git checkout master".split())
        #end Branch().create()

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

def get_version() -> str:
    branch = subprocess.run("git branch --show-current".split(),
            capture_output=True, encoding='utf-8')
    if not branch.stdout.startswith("release"):
        print(f"You are not on a release branch! ({branch.stdout.strip()})")
        sys.exit(1)
    return branch.stdout.strip().split("/")[1]

def get_release_branches() -> List[Branch]:
    branches = []
    _b = subprocess.run("git branch --list".split(),
            capture_output=True, encoding='utf-8')
    result = _b.stdout.split()
    for b in result:
        if b.strip().startswith("release"):
            branches.append(Branch(b.strip()))
    return branches

def get_latest(branches=None) -> Branch:
    if not branches:
        return max(get_release_branches())
    else:
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
    group.add_argument("-M", "--major", action='store_true', help="create a new major-version release")
    group.add_argument("-m", "--minor", action='store_true', help="create a new minor-version release")
    group.add_argument("-p", "--patch", action='store_true', help="create a new patch-version release")
    group.add_argument("-v", "--get-version", action='store_true', help="print the version number of the current branch and exit")
    group.add_argument("-r", "--latest-release", action='store_true', help="print the latest release and exit")
    group.add_argument("-l", "--list-releases", action='store_true', help="print all branches named 'release/*' and exit")
    args = parser.parse_args()

    if args.major:
        new_branch = get_next_release('major')
        print(f"{get_latest()} --> {new_branch}")
        new_branch.create()
    elif args.minor:
        new_branch = get_next_release('minor')
        print(f"{get_latest()} --> {new_branch}")
        new_branch.create()
    elif args.patch:
        new_branch = get_next_release('patch')
        print(f"{get_latest()} --> {new_branch}")
        new_branch.create()
    elif args.latest_release:
        print(get_latest(get_release_branches()))
    elif args.list_releases:
        [print(br) for br in get_release_branches()]
    elif args.get_version:
        print(get_version())
    else:
        parser.print_help()
    sys.exit(0)
main()
