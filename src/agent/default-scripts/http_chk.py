#!/usr/bin/env python3
import json
from sys import argv
from dataclasses import dataclass
import requests
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

r = requests.get('http://' + check.host)

print(r.text)
print(r, file=sys.stderr) 