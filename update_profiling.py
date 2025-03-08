import subprocess

from pathlib import Path

target_folder = Path(__file__).parent / '.prof'
if not target_folder.exists():
    target_folder.mkdir()

profiles = [
    'goroutine',  # stack traces of all current goroutines
    'heap',  # a sampling of memory allocations of live objects
    'allocs',  # a sampling of all past memory allocations
    'threadcreate',  # stack traces that led to the creation of new OS threads
    'block',  # stack traces that led to blocking on synchronization primitives
    'mutex',  # stack traces of holders of contended mutexes
]

for profile in profiles:
    p = subprocess.Popen(f'curl -sK -v http://localhost:8080/debug/pprof/{profile} > {target_folder / profile}.out', shell=True)
    p.wait()
