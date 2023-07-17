import os
import re
import json
import requests

def compose_log() -> dict:
    path = os.environ['INPUT_PATH']
    files = [s for s in os.listdir(path) if 'step_' in s]
    steps = []

    for file_name in files:
        match = re.findall(r'\d+', file_name)
        if len(match) > 1:
            print(f"Warning: file name {file_name} contains more than 1 digit")
        step_number = int(match[0])

        with open(os.path.join(path, file_name), 'r') as f:
            time_result, successful = f.read().split()
            successful = successful.lower() == 'true'
            r = re.match(r'(\d*)m(\d*.\d*)s', time_result)
            time_result = int(r[1])*60 + float(r[2])

        d = {
            "id" : step_number,
            "exec_time": str(time_result),
            "successful": successful,
        }

        steps.append(d)

    return {
        "name": os.environ['INPUT_NAME'],
        "start": os.environ['INPUT_START'],
        "end": os.environ['INPUT_END'],
        "successful": os.environ['INPUT_SUCCESSFUL'].lower() == 'true',
        "arch": os.environ['INPUT_ARCH'],
        "steps": steps
    }

def run():
    auth_header = {'auth' : os.environ['AUTH_TOKEN']}

    u = os.environ['GHA_URL'] + '/ping'
    r = requests.get(u)
    assert r.status_code == 200

    body = compose_log()
    u = os.environ['GHA_URL'] + '/api/v1/publish/timing'
    r = requests.post(u, json=body, headers=auth_header)
    print(body)
    print(r.url)
    print(r.status_code)
    print(r.content)

if __name__ == '__main__':
    run()
