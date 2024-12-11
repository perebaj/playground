import google.cloud.secretmanager as secretmanager
from google.api_core.retry import Retry


client = secretmanager.SecretManagerServiceClient()


# According the documentation, this retry policy will sort a random delay between 0.1 and 1.0 seconds
# and multiply it by 2 every time it retries, until it reaches 10 seconds.
retry = Retry(initial=0.1, maximum=1.0, multiplier=2, timeout=5)

for secret in client.list_secrets(request={"parent": "projects/credit-engine-prd-001"}, retry=retry):
    print(f"Fount secret: {secret.name}")
