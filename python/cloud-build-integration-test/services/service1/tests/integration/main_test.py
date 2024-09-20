import requests

def test_hello():
    response = requests.get("http://localhost:8000/hello")
    assert response.status_code == 200
    assert response.json() == {"message": "Hello World"}
