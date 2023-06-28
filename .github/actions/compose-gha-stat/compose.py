import os
import re
import json

def compose_log() -> dict:
    path = os.environ['INPUT_PATH']
    files = [s for s in os.listdir(path) if 'step_' in s]
    steps = []

    for file_name in files:
        match = re.findall(r'\d+', file_name)
        if len(match) > 1:
            print(f"Warning: file name {file_name} contains more than 1 digit")
        step_number = match[0]

        with open(os.path.join(path, file_name), 'r') as f:
            time_result, successful = f.read().split()
            successful = successful.lower() == 'true'

        d = {
            "id" : step_number,
            "exec_time": time_result,
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
    body = compose_log()
    auth_header = {'auth' : os.environ['AUTH_TOKEN']}
    r = requests.post(os.environ['GHA_URL'], data=body, headers=auth_header)

if __name__ == '__main__':
    run()
