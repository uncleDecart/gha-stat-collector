import os
from compose import compose_log
import pytest
import shutil

@pytest.fixture()
def log_folder():
    log_folder = "test_logs"
    os.environ['INPUT_PATH'] = log_folder
    os.environ['INPUT_NAME'] = "test_action"
    os.environ['INPUT_START'] = "2022-02-01"
    os.environ['INPUT_END'] = "2022-02-02"
    os.environ['INPUT_SUCCESSFUL'] = "true"
    os.environ['INPUT_ARCH'] = "x64"
    os.mkdir(log_folder)

    with open(os.path.join(log_folder, 'step_1.txt'), 'w') as f:
        f.write("1 true\n")
    with open(os.path.join(log_folder, 'step_2.txt'), 'w') as f:
        f.write("5 false\n")

    yield log_folder

    shutil.rmtree(log_folder)

class TestCompose:
    def test_compose_log(self, log_folder):
        expected = {
            "name" : "test_action",
            "start" : "2022-02-01",
            "end" : "2022-02-02",
            "successful" : True,
            "arch" : "x64",
            "steps" : [{
                "id": "1",
                "exec_time": "1",
                "successful": True,
            },{
                "id": "2",
                "exec_time": "5",
                "successful": False,
            },],
        }
        got = compose_log()
        got['steps'] = sorted(got['steps'], key=lambda x: x['id'])
        assert expected == got
