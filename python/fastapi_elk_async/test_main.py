import pytest
from httpx import AsyncClient
from main import app
import asyncio
import unittest


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


from unittest.mock import AsyncMock, patch
import httpx


class TestClass(BaseTest):
    @patch(
        "main.httpx.AsyncClient.post",
        return_value=httpx.Response(200, json={"id": "9ed7dasdasd-08ff-4ae1-8952-37e3a323eb08"}),
    )
    async def test_jojo(self, mock_post):
        # mock_response = mock_post.return_value
        # mock_response.json.return_value = {"key": "value"}
        async with AsyncClient(app=app, base_url="http://test") as ac:
            response = await ac.get("/")
            assert response.status_code == 200
            assert response.json() == {"id": "9ed7dasdasd-08ff-4ae1-8952-37e3a323eb0"}

            print(response)
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
