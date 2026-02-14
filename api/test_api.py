#!/usr/bin/env python3
"""
Quincy API endpoint test suite.

Requires a running Quincy API server with a valid config.
Fetches the config from the server to dynamically build test data.

Usage:
    python test_api.py [BASE_URL]

    BASE_URL defaults to http://127.0.0.1:8888
"""

import sys
import json
import time
import requests

BASE_URL = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8888"
API = f"{BASE_URL}/api/v1"

passed = 0
failed = 0
errors = []


def report(name, ok, detail=""):
    global passed, failed
    if ok:
        passed += 1
        print(f"  PASS  {name}")
    else:
        failed += 1
        errors.append((name, detail))
        print(f"  FAIL  {name}  --  {detail}")


def check(name, condition, detail=""):
    report(name, condition, detail)
    return condition


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def get(path):
    return requests.get(f"{API}{path}", timeout=10)


def post(path, body):
    return requests.post(f"{API}{path}", json=body, timeout=10)


def post_raw(path, data, content_type="application/json"):
    return requests.post(
        f"{API}{path}",
        data=data,
        headers={"Content-Type": content_type},
        timeout=10,
    )


# ===========================================================================
# Test cases
# ===========================================================================

def test_not_found():
    print("\n--- 404 / No Route ---")

    r = requests.get(f"{API}/nonexistent", timeout=10)
    check("GET unknown path returns 404", r.status_code == 404,
          f"got {r.status_code}")

    body = r.json()
    check("404 body has message field", "message" in body)
    check("404 body has path field", "path" in body)
    check("404 body has method field", "method" in body)
    check("404 message value", body.get("message") == "route not found or method not allowed",
          f"got: {body.get('message')}")

    r2 = requests.delete(f"{API}/scores", timeout=10)
    check("DELETE on valid path returns 404", r2.status_code == 404,
          f"got {r2.status_code}")


def test_get_config():
    print("\n--- GET /config ---")

    r = get("/config")
    check("GET /config returns 200", r.status_code == 200, f"got {r.status_code}")

    body = r.json()
    check("config has num_teams", "num_teams" in body)
    check("config has db_file", "db_file" in body)
    check("config has boxes", "boxes" in body and isinstance(body["boxes"], list))
    check("config has user_lists", "user_lists" in body and isinstance(body["user_lists"], list))
    check("config has http", "http" in body)
    check("num_teams > 0", body["num_teams"] > 0, f"got {body['num_teams']}")
    check("boxes is non-empty", len(body["boxes"]) > 0)
    check("user_lists is non-empty", len(body["user_lists"]) > 0)

    # Validate box structure
    box = body["boxes"][0]
    check("box has name", "name" in box)
    check("box has id", "id" in box)
    check("box has services", "services" in box and isinstance(box["services"], list))
    check("box services non-empty", len(box["services"]) > 0)

    svc = box["services"][0]
    check("service has name", "name" in svc)
    check("service has id", "id" in svc)
    check("service has check", "check" in svc)

    # Validate userlist structure
    ul = body["user_lists"][0]
    check("userlist has name", "name" in ul)
    check("userlist has id", "id" in ul)
    check("userlist has users", "users" in ul and isinstance(ul["users"], list))
    check("userlist users non-empty", len(ul["users"]) > 0)

    user = ul["users"][0]
    check("user has username", "username" in user)
    check("user has password", "password" in user)

    # Validate http structure
    http_cfg = body["http"]
    check("http has host", "host" in http_cfg)
    check("http has port", "port" in http_cfg)

    return body


def test_score_endpoints_empty(config):
    """Test score endpoints before any scores are submitted."""
    print("\n--- Score Endpoints (pre-submission) ---")

    # Team scores -- should return an empty map or map with zero-value results
    r = get("/scores/team")
    check("GET /scores/team returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("/scores/team returns object", isinstance(body, dict))

    # If there are entries, validate their structure
    if body:
        key = list(body.keys())[0]
        entry = body[key]
        check("team score has checks_passed", "checks_passed" in entry)
        check("team score has checks_failed", "checks_failed" in entry)
        check("team score has total_checks", "total_checks" in entry)
        check("team score has uptime_percent", "uptime_percent" in entry)

    # Box scores
    r = get("/scores/box")
    check("GET /scores/box returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("/scores/box returns object", isinstance(body, dict))

    # Service scores
    r = get("/scores/service")
    check("GET /scores/service returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("/scores/service returns object", isinstance(body, dict))

    # Current status
    r = get("/scores/current")
    check("GET /scores/current returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("/scores/current returns array", isinstance(body, list))

    # Detailed scores
    r = get("/scores/detailed")
    check("GET /scores/detailed returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("/scores/detailed returns object", isinstance(body, dict))


def test_add_score(config):
    """Test adding scores and verifying they appear in all views."""
    print("\n--- POST /scores ---")

    box_id = config["boxes"][0]["id"]
    service_id = config["boxes"][0]["services"][0]["id"]
    team_num = 1

    # Submit a passing score
    score_pass = {
        "team_num": team_num,
        "status": True,
        "box": box_id,
        "service": service_id,
        "message": "test pass",
        "timestamp": 0,
    }
    r = post("/scores", score_pass)
    check("POST /scores (pass) returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("response has message", body.get("message") == "score added successfully",
          f"got: {body.get('message')}")
    check("response echoes score", "score" in body)

    # Submit a failing score
    score_fail = {
        "team_num": team_num,
        "status": False,
        "box": box_id,
        "service": service_id,
        "message": "test fail",
        "timestamp": 0,
    }
    r = post("/scores", score_fail)
    check("POST /scores (fail) returns 200", r.status_code == 200, f"got {r.status_code}")

    # Submit another passing score
    r = post("/scores", score_pass)
    check("POST /scores (pass #2) returns 200", r.status_code == 200, f"got {r.status_code}")

    return box_id, service_id, team_num


def test_add_score_bad_input():
    """Test adding scores with invalid input."""
    print("\n--- POST /scores (bad input) ---")

    # Invalid JSON
    r = post_raw("/scores", "not json at all")
    check("POST /scores invalid JSON returns 400", r.status_code == 400,
          f"got {r.status_code}")

    # Empty body
    r = post_raw("/scores", "")
    check("POST /scores empty body returns 400", r.status_code == 400,
          f"got {r.status_code}")

    # Valid JSON but score for non-existent box/service combo
    bad_score = {
        "team_num": 1,
        "status": True,
        "box": "nonexistent-box-xxx",
        "service": "nonexistent-svc",
        "message": "should fail",
        "timestamp": 0,
    }
    r = post("/scores", bad_score)
    check("POST /scores non-existent box/service returns 400", r.status_code == 400,
          f"got {r.status_code}")


def test_scores_after_submission(box_id, service_id, team_num):
    """Verify score aggregations after adding test scores."""
    print("\n--- Score Aggregation Verification ---")

    # Team scores -- team 1 should have results
    r = get("/scores/team")
    check("GET /scores/team returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    team_key = str(team_num)
    check(f"team {team_num} present in team scores", team_key in body,
          f"keys: {list(body.keys())}")

    if team_key in body:
        ts = body[team_key]
        check("team total_checks > 0", ts["total_checks"] > 0,
              f"got {ts['total_checks']}")
        check("team checks_passed > 0", ts["checks_passed"] > 0,
              f"got {ts['checks_passed']}")
        check("uptime_percent is a number", isinstance(ts["uptime_percent"], (int, float)))

    # Box scores -- our box should have results
    r = get("/scores/box")
    body = r.json()
    check(f"box '{box_id}' present in box scores", box_id in body,
          f"keys: {list(body.keys())}")

    # Service scores -- our box+service should have results
    r = get("/scores/service")
    body = r.json()
    check(f"box '{box_id}' in service scores", box_id in body)
    if box_id in body:
        check(f"service '{service_id}' under box", service_id in body[box_id],
              f"keys: {list(body[box_id].keys())}")

    # Detailed scores -- full nesting check
    r = get("/scores/detailed")
    body = r.json()
    check(f"team {team_num} in detailed scores", team_key in body)
    if team_key in body:
        check(f"box '{box_id}' in detailed[{team_num}]", box_id in body[team_key])
        if box_id in body[team_key]:
            check(f"service '{service_id}' in detailed[{team_num}][{box_id}]",
                  service_id in body[team_key][box_id])

    # Current status -- should have our most recent check
    r = get("/scores/current")
    body = r.json()
    check("/scores/current is non-empty after submissions", len(body) > 0)

    if body:
        # Find our specific entry
        our_entry = None
        for entry in body:
            if (entry.get("team_num") == team_num and
                    entry.get("box") == box_id and
                    entry.get("service") == service_id):
                our_entry = entry
                break

        check("our team/box/service found in current status", our_entry is not None)
        if our_entry:
            check("current entry has status field", "status" in our_entry)
            check("current entry has message field", "message" in our_entry)
            check("current entry has timestamp > 0", our_entry.get("timestamp", 0) > 0,
                  f"got {our_entry.get('timestamp')}")
            # Last submitted score was a pass
            check("current status reflects last submission (pass)",
                  our_entry["status"] is True,
                  f"got {our_entry['status']}")

    # Verify sorting: team_num, then box, then service
    if len(body) > 1:
        is_sorted = True
        for i in range(len(body) - 1):
            a, b = body[i], body[i + 1]
            key_a = (a["team_num"], a["box"], a["service"])
            key_b = (b["team_num"], b["box"], b["service"])
            if key_a > key_b:
                is_sorted = False
                break
        check("/scores/current is sorted by team/box/service", is_sorted)


def test_get_check(config):
    """Test the agent check endpoint."""
    print("\n--- GET /checks ---")

    r = get("/checks")
    check("GET /checks returns 200", r.status_code == 200, f"got {r.status_code}")

    body = r.json()
    check("check has name", "name" in body)
    check("check has id", "id" in body)
    check("check has check (CheckID)", "check" in body)
    check("check has box", "box" in body)
    check("check has host", "host" in body)
    check("check has team_num", "team_num" in body)
    check("check team_num in valid range",
          1 <= body.get("team_num", 0) <= config["num_teams"],
          f"got {body.get('team_num')}")

    # Verify the box ID is from config
    valid_box_ids = [b["id"] for b in config["boxes"]]
    check("check box is from config", body.get("box") in valid_box_ids,
          f"got {body.get('box')}, valid: {valid_box_ids}")

    # Call multiple times and verify we get different checks (queue rotation)
    seen_combos = set()
    seen_combos.add((body.get("team_num"), body.get("box"), body.get("id")))
    for _ in range(5):
        r2 = get("/checks")
        if r2.status_code == 200:
            b2 = r2.json()
            seen_combos.add((b2.get("team_num"), b2.get("box"), b2.get("id")))

    check("multiple /checks calls return varied results", len(seen_combos) > 1,
          f"only saw {len(seen_combos)} unique combo(s)")


def test_get_users(config):
    """Test the user enumeration endpoint."""
    print("\n--- GET /users ---")

    r = get("/users")
    check("GET /users returns 200", r.status_code == 200, f"got {r.status_code}")

    body = r.json()
    check("/users returns object", isinstance(body, dict))

    num_teams = config["num_teams"]
    check(f"/users has entries for {num_teams} team(s)", len(body) == num_teams,
          f"got {len(body)} team(s)")

    # Validate structure for team 1
    team_key = "1"
    if team_key in body:
        team_data = body[team_key]
        check(f"team {team_key} data is object", isinstance(team_data, dict))

        config_ul_ids = [ul["id"] for ul in config["user_lists"]]
        for ul_id in config_ul_ids:
            check(f"userlist '{ul_id}' present for team {team_key}", ul_id in team_data,
                  f"keys: {list(team_data.keys())}")

            if ul_id in team_data:
                users = team_data[ul_id]
                check(f"userlist '{ul_id}' is array", isinstance(users, list))
                check(f"userlist '{ul_id}' is non-empty", len(users) > 0)

                if users:
                    u = users[0]
                    check("user has username", "username" in u)
                    check("user has password", "password" in u)


def test_pcr(config):
    """Test password change requests."""
    print("\n--- POST /users (PCR) ---")

    ul = config["user_lists"][0]
    ul_id = ul["id"]
    username = ul["users"][0]["username"]
    original_password = ul["users"][0]["password"]
    team_num = 1
    new_password = f"TestPCR_{int(time.time())}"

    # Successful PCR
    pcr = {
        "username": username,
        "password": new_password,
        "domain": "",
        "netbios": "",
        "user_list": ul_id,
        "team_num": team_num,
    }
    r = post("/users", pcr)
    check("POST /users (PCR) returns 200", r.status_code == 200, f"got {r.status_code}")
    body = r.json()
    check("PCR response has message 'user updated'",
          body.get("message") == "user updated",
          f"got: {body.get('message')}")
    check("PCR response echoes pcr object", "pcr" in body)

    # Verify the password was actually changed
    r = get("/users")
    users_body = r.json()
    team_key = str(team_num)
    found = False
    if team_key in users_body and ul_id in users_body[team_key]:
        for u in users_body[team_key][ul_id]:
            if u["username"] == username:
                found = True
                check("password was updated in database",
                      u["password"] == new_password,
                      f"expected '{new_password}', got '{u['password']}'")
                break
    check("target user found after PCR", found)

    # Restore original password
    restore_pcr = {
        "username": username,
        "password": original_password,
        "domain": "",
        "netbios": "",
        "user_list": ul_id,
        "team_num": team_num,
    }
    r = post("/users", restore_pcr)
    check("restore original password returns 200", r.status_code == 200,
          f"got {r.status_code}")

    return ul_id, username, team_num


def test_pcr_bad_input(config):
    """Test PCR with various invalid inputs."""
    print("\n--- POST /users (PCR bad input) ---")

    ul_id = config["user_lists"][0]["id"]

    # Invalid JSON
    r = post_raw("/users", "not json")
    check("PCR invalid JSON returns 400", r.status_code == 400, f"got {r.status_code}")

    # Non-existent userlist
    pcr = {
        "username": "admin",
        "password": "newpass",
        "user_list": "nonexistent_userlist_xyz",
        "team_num": 1,
    }
    r = post("/users", pcr)
    check("PCR non-existent userlist returns 400", r.status_code == 400,
          f"got {r.status_code}")
    body = r.json()
    check("error mentions userlist doesn't exist",
          "does not exist" in body.get("message", ""),
          f"got: {body.get('message')}")

    # Non-existent username (valid userlist) -- should return 500
    pcr = {
        "username": "definitely_not_a_real_user_xyz",
        "password": "newpass",
        "user_list": ul_id,
        "team_num": 1,
    }
    r = post("/users", pcr)
    check("PCR non-existent username returns 500", r.status_code == 500,
          f"got {r.status_code}")

    # Team number 0 (invalid -- teams start at 1) with valid userlist
    pcr = {
        "username": config["user_lists"][0]["users"][0]["username"],
        "password": "newpass",
        "user_list": ul_id,
        "team_num": 0,
    }
    r = post("/users", pcr)
    check("PCR team_num=0 returns 500", r.status_code == 500,
          f"got {r.status_code}")


def test_score_result_math(config):
    """Verify that ScoreResult math is correct after known submissions."""
    print("\n--- Score Result Math Verification ---")

    # Use a second box/service if available, otherwise use first
    # to avoid interference with earlier test scores
    if len(config["boxes"]) > 1:
        box = config["boxes"][1]
    else:
        box = config["boxes"][0]

    if len(box["services"]) > 1:
        svc = box["services"][-1]
    else:
        svc = box["services"][0]

    box_id = box["id"]
    service_id = svc["id"]
    # Use the last team to minimize interference
    team_num = config["num_teams"]

    # Get current stats so we can calculate expected values
    r = get("/scores/detailed")
    pre = r.json()
    team_key = str(team_num)

    pre_total = 0
    pre_passed = 0
    if team_key in pre and box_id in pre[team_key] and service_id in pre[team_key][box_id]:
        pre_total = pre[team_key][box_id][service_id]["total_checks"]
        pre_passed = pre[team_key][box_id][service_id]["checks_passed"]

    # Submit exactly 3 passes and 2 fails
    for _ in range(3):
        post("/scores", {
            "team_num": team_num, "status": True,
            "box": box_id, "service": service_id,
            "message": "math test pass", "timestamp": 0,
        })
    for _ in range(2):
        post("/scores", {
            "team_num": team_num, "status": False,
            "box": box_id, "service": service_id,
            "message": "math test fail", "timestamp": 0,
        })

    expected_total = pre_total + 5
    expected_passed = pre_passed + 3
    expected_failed = expected_total - expected_passed
    expected_uptime = round(expected_passed / expected_total * 100, 2)

    r = get("/scores/detailed")
    post_body = r.json()

    check(f"team {team_num} in detailed after math test", team_key in post_body)
    if team_key in post_body and box_id in post_body[team_key]:
        result = post_body[team_key][box_id].get(service_id, {})
        check("total_checks matches expected",
              result.get("total_checks") == expected_total,
              f"expected {expected_total}, got {result.get('total_checks')}")
        check("checks_passed matches expected",
              result.get("checks_passed") == expected_passed,
              f"expected {expected_passed}, got {result.get('checks_passed')}")
        check("checks_failed matches expected",
              result.get("checks_failed") == expected_failed,
              f"expected {expected_failed}, got {result.get('checks_failed')}")
        check("uptime_percent matches expected",
              result.get("uptime_percent") == expected_uptime,
              f"expected {expected_uptime}, got {result.get('uptime_percent')}")


def test_method_not_allowed():
    """Test that wrong HTTP methods return 404 (no route)."""
    print("\n--- Method Not Allowed ---")

    # GET on POST-only endpoint
    r = requests.get(f"{API}/scores", timeout=10)
    check("GET /scores (POST-only) returns 404", r.status_code == 404,
          f"got {r.status_code}")

    # POST on GET-only endpoint
    r = requests.post(f"{API}/scores/team", json={}, timeout=10)
    check("POST /scores/team (GET-only) returns 404", r.status_code == 404,
          f"got {r.status_code}")

    # PUT on any endpoint (not used anywhere)
    r = requests.put(f"{API}/scores", json={}, timeout=10)
    check("PUT /scores returns 404", r.status_code == 404,
          f"got {r.status_code}")


# ===========================================================================
# Main
# ===========================================================================

def main():
    global passed, failed

    print(f"Quincy API Test Suite")
    print(f"Target: {BASE_URL}")
    print("=" * 60)

    # Verify server is reachable
    try:
        r = requests.get(f"{BASE_URL}/api/v1/config", timeout=5)
        r.raise_for_status()
    except requests.exceptions.ConnectionError:
        print(f"\nERROR: Cannot connect to {BASE_URL}")
        print("Make sure the Quincy API server is running.")
        sys.exit(1)
    except Exception as e:
        print(f"\nERROR: {e}")
        sys.exit(1)

    # Run all tests
    config = test_get_config()
    test_not_found()
    test_method_not_allowed()
    test_score_endpoints_empty(config)
    test_get_check(config)
    test_add_score(config)
    test_add_score_bad_input()
    # Re-extract IDs for post-submission checks
    box_id = config["boxes"][0]["id"]
    service_id = config["boxes"][0]["services"][0]["id"]
    team_num = 1
    test_scores_after_submission(box_id, service_id, team_num)
    test_score_result_math(config)
    test_get_users(config)
    test_pcr(config)
    test_pcr_bad_input(config)

    # Summary
    total = passed + failed
    print("\n" + "=" * 60)
    print(f"Results: {passed}/{total} passed, {failed} failed")

    if errors:
        print(f"\nFailed tests:")
        for name, detail in errors:
            print(f"  - {name}: {detail}")

    sys.exit(0 if failed == 0 else 1)


if __name__ == "__main__":
    main()
