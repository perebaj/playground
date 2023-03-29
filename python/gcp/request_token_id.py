import google.oauth2.id_token
import google.auth.transport.requests

request = google.auth.transport.requests.Request()
target_audience = "https://embedding-distiluse-mult-development-yxltmkv3ea-uc.a.run.app/predict"

id_token = google.oauth2.id_token.fetch_id_token(request, target_audience)
print(id_token)