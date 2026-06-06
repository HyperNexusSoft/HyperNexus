import time
import urllib.request
import json

BASE_URL = "http://localhost:4300"

def call_tool(name, arguments):
    url = f"{BASE_URL}/api/agent/tool"
    payload = {
        "name": name,
        "arguments": arguments
    }
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, method='POST', headers={'Content-Type': 'application/json'})

    print(f"--- Executing Tool: {name} ---")
    start = time.perf_counter()
    try:
        with urllib.request.urlopen(req) as response:
            res_body = response.read().decode('utf-8')
            end = time.perf_counter()
            duration = (end - start) * 1000
            result = json.loads(res_body)
            print(f"Status: Success | Latency: {duration:.2f}ms")
            # print(f"Output: {res_body[:200]}...")
            return result
    except Exception as e:
        print(f"Status: Failed | Error: {e}")
        return None

def run_workload():
    print("🚀 Starting TormentNexus Autonomous Workload Verification\n")

    # Step 1: List skills to find relevant capabilities
    skills = call_tool("skill_list", {})

    # Step 2: Search for a specific prompt template
    prompts = call_tool("prompt_search", {"query": "refactor"})

    # Step 3: Use a shell tool to inspect the environment
    ls_result = call_tool("ls", {"path": "."})

    # Step 4: Simulate a session/harness execution (non-blocking status check)
    harness_result = call_tool("pi_mono", {"task": "Verify environment health", "command": "status"})

    # Step 5: Get a specific skill content
    if skills and "data" in skills:
        skill_get = call_tool("skill_get", {"name": "TestSkill"})

    print("\n✅ Autonomous Workload Verification Complete")

if __name__ == "__main__":
    run_workload()
