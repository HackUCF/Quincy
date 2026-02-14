#!/usr/bin/env python3
import json
from sys import exit, argv
from random import choice

check_obj = None
with open(argv[1], 'r') as f:
  check_obj = json.load(f)

print(__file__, check_obj)

code = choice([0, 1, 1, 1])
print('failure' if code else 'success')
exit(code)