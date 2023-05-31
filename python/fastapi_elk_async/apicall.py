import httpx


class APIClusterKeyword:
    async def cluster_keywords_search(self):
        async with httpx.AsyncClient() as client:
            response = await client.get("http://example.com/keyword_v2/cluster_search")
        return response.json()
