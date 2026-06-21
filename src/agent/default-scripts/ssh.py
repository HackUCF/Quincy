#!/usr/bin/env python3
import json
from sys import argv
from dataclasses import dataclass
import paramiko
import sys

@dataclass
class User:
  username: str
  password: str
  domain: str | None = None
  netbios: str | None = None

@dataclass
class Check:
  name: str
  check: str
  box: str
  host: str
  team_num: int
  timeout: float | None = None
  user_list: str | None = None
  user: User | None = None

check_obj = None
with open(argv[1], 'r') as f:
  data = json.load(f)
  if "user" in data:
    data["user"] = User(**data["user"])
  check = Check(**data)

ssh_client = paramiko.SSHClient()
ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
ssh_client.connect(
  hostname=check.host,
  port=22,
  username=check.user.username,
  password=check.user.password,
  timeout=(check.timeout * 3 / 4),
)

stdin, stdout, stderr = ssh_client.exec_command('whoami')
print(stdout.readlines())
print(stderr.readlines(), file=sys.stderr)