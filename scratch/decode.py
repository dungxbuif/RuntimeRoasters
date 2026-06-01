import base64
import re

with open('script.b64', 'r') as f:
    raw = f.read()

# Replace literal \n sequence (two chars) and actual newlines
cleaned = raw.replace('\\n', '').replace('\n', '').strip()

# Keep only valid base64 characters: A-Z, a-z, 0-9, +, /, =
cleaned = re.sub(r'[^A-Za-z0-9+/=]', '', cleaned)

# Pad if necessary
missing_padding = len(cleaned) % 4
if missing_padding:
    cleaned += '=' * (4 - missing_padding)

try:
    decoded = base64.b64decode(cleaned)
    with open('update-harness.sh', 'wb') as out:
        out.write(decoded)
    print("Successfully decoded script to update-harness.sh")
except Exception as e:
    print(f"Error decoding base64: {e}")
