from fastapi import FastAPI
from elasticsearch import AsyncElasticsearch

auth = ("admin", "adminadmin")  # For testing only. Don't store credentials in code.
es = AsyncElasticsearch(hosts=["http://localhost:9200"])

app = FastAPI()


@app.on_event("shutdown")
async def app_shutdown():
    await es.close()  # This gets called once the app is shutting down.


@app.get("/")
async def index():
    return await es.cluster.health()
