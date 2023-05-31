from fastapi import FastAPI
from elasticsearch import AsyncElasticsearch, Elasticsearch
from asyncio import sleep
import asyncio
import httpx

auth = ("admin", "adminadmin")  # For testing only. Don't store credentials in code.
es = AsyncElasticsearch(hosts=["http://localhost:9200"])
es_sync = Elasticsearch(
    hosts=["http://localhost:9200"],
)
es.async_search

app = FastAPI()


@app.on_event("shutdown")
async def app_shutdown():
    await es.close()  # This gets called once the app is shutting down.


from fastapi_elk_async.apicall import APIClusterKeyword


@app.get("/")
async def index():
    await asyncio.sleep(1)
    apicall = APIClusterKeyword()
    print("JOJO")
    print("JOJO PRINT")

    response = await apicall.cluster_keywords_search()
    # json = response.json()
    response["name"] = "Jojo"
    return response


@app.get("/notasync")
def notasync():
    return es_sync.cluster.health()


@app.get("/test")
async def test():
    return {"message": "Tomato"}
