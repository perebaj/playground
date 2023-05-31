from httpx import AsyncClient
from fastapi_elk_async.main import app
import unittest
import time


class BaseTest(unittest.IsolatedAsyncioTestCase):
    @classmethod
    def setUpClass(cls) -> None:
        print("setUpClass")
        # return super().setUpClass()

    def tearDown(self) -> None:
        print("tearDown")
        # return super().tearDown()

    # def tearDownClass(self, cls):
    #     print("tearDownClass")


from unittest.mock import patch
import httpx

from fastapi_elk_async.main import app


class TestClass(BaseTest):
    # @patch("httpx.AsyncClient.get", wraps=httpx.AsyncClient.get)
    async def test_jojo(self):
        with patch("httpx.AsyncClient.get") as mock_client:
            data = {"id": "9ed7dasdasd-08ff-4ae1-8952-37e3a323eb08"}
            mock_client.return_value = httpx.Response(200, json=data)
            # mock_response = mock_post.return_value
            # mock_response.json.return_value = {"key": "value"}
            async with AsyncClient(app=app, base_url="http://test") as ac:
                response = await ac.request("GET", "/")
                # response = await ac.get("/")
                print("teste reposnse> ", response.json())
            # assert response.status_code == 200

        # assert (await mock()) == "jojo"
        # print("test_jojo")
        # assert True

    # async def test_async_route(self):
    #     async with AsyncClient(app=app, base_url="http://test") as ac:
    #         response = await ac.get("/")
    #     assert response.status_code == 200

    # async def test_async_route2(self):
    #     async with AsyncClient(app=app, base_url="http://test") as ac:
    #         response = await ac.get("/")
    #     assert response.status_code == 200
    #     assert response.json() == {"message": "Hello World"}

    # async def test_async_route3(self):
    #     async with AsyncClient(app=app, base_url="http://test") as ac:
    #         response = await ac.get("/")
    #     assert response.status_code == 200
    #     assert response.json() == {"message": "Hello World"}


# @pytest.mark.asyncio
# class TestClass:
#     @pytest.mark.asyncio
#     async def test_root(self):
#         async with AsyncClient(app=app, base_url="http://test") as ac:
#             response = await ac.get("/")
#         assert response.status_code == 200

#     async def test_root_n(self):
#         async with AsyncClient(app=app, base_url="http://test") as ac:
#             response = await ac.get("/")
#         assert response.status_code == 200
