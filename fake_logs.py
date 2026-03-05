import time
import random
import sys

services = ['auth-service', 'payment-service', 'db-service', 'api-gateway', 'worker-queue']
levels = ['INFO', 'INFO', 'INFO', 'INFO', 'WARN', 'ERROR', 'DEBUG']
messages = {
    'auth-service': ['user login success', 'invalid jwt token', 'token expired', 'user logout', 'password reset requested'],
    'payment-service': ['card verification failed', 'retrying payment request', 'payment success', 'insufficient funds', 'refund processed'],
    'db-service': ['connection timeout', 'query taking too long', 'deadlock detected', 'connection refused', 'index rebuilt'],
    'api-gateway': ['request routed', 'rate limit exceeded', '404 not found', '502 bad gateway', 'ssl handshake failed'],
    'worker-queue': ['job processed', 'job failed', 'queue full', 'worker node detached', 'new worker joined']
}

print("Running fake_logs.py... Pipe this to LitLog via: python3 fake_logs.py | ./litelog ingest", file=sys.stderr)

try:
    while True:
        service = random.choice(services)
        level = random.choice(levels)
        message = random.choice(messages[service])
        
        # Write to stdout so it can be piped to litelog ingest
        log_line = f"{level} {service}: {message}"
        print(log_line, flush=True)
        
        time.sleep(random.uniform(0.01, 0.8)) # Random delay between logs
except KeyboardInterrupt:
    sys.exit(0)
