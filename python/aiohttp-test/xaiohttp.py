import aiohttp
import asyncio


async def fetch():
    async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(0.2)) as session:
        try:
            async with session.get('https://jsonplaceholder.typicode.com/todos/1') as response:

                print("Status:", response.status)
                print("Content-type:", response.headers['content-type'])

                html = await response.text()
                r_json = await response.json()
                print("Body:", html[:15], "...")
                print("JSON:", r_json)
                return response.status, response.headers['content-type'], r_json
        except asyncio.TimeoutError as e:
            print("TimeoutError:", str(e))
            print("ConnectionTimeoutError:", str(e))
            return None, None, None

async def post_fetch_example():
    async with aiohttp.ClientSession() as session:
        async with session.post('https://jsonplaceholder.typicode.com/posts', json={"title": "foo", "body": "bar", "userId": 1}) as response:
            print("Status:", response.status)
            print("Content-type:", response.headers['content-type'])

            html = await response.text()
            r_json = await response.json()
            print("Body:", html[:15], "...")
            print("JSON:", r_json)
            return response.status, response.headers['content-type'], r_json

asyncio.run(fetch())
