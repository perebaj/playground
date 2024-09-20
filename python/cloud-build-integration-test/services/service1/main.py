#create a hello endpoint

from fastapi import FastAPI

app = FastAPI()

@app.get("/hello")

def hello():
    return {"message": "Hello World"}
