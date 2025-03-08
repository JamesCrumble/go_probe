import time
import requests

CYCLES: int = 1000
DELAY_SECONDS: float = 0

probes = [
    'psql', 'kafka',
]
prode_dst = 'http://localhost:4000/probe/'

for probe in probes:
    prode_endpoint = prode_dst + probe
    for _ in range(CYCLES):
        try:
            response = requests.get(prode_endpoint, timeout=10)
            print(response.status_code, len(response.json()))
        except Exception as exc:
            print(exc)

        time.sleep(DELAY_SECONDS)
