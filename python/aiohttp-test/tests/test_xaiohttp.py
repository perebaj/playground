from unittest.mock import patch, MagicMock, AsyncMock
import pytest
from xaiohttp import fetch, post_fetch_example
import asyncio


def test_hello():
    assert 1 == 1


@pytest.mark.asyncio
@patch('aiohttp.ClientSession.get')
async def test_fetch(mock_get: MagicMock):
    """Test fetch."""
    mock_get.return_value.__aenter__.return_value.status = 200
    mock_get.return_value.__aenter__.return_value.headers = {'content-type': 'text/html', 'Connection': 'close'}
    # mock the response.json method
    mock_get.return_value.__aenter__.return_value.json = AsyncMock(return_value={'key': 'value'})

    status, headers, r_json = await fetch()

    assert status == 200
    assert headers == 'text/html'
    assert r_json == {'key': 'value'}


# @pytest.mark.asyncio
# @patch('aiohttp.ClientSession.get')
# async def test_fetch_timeout(mock_get: AsyncMock):
#     """Test fetch."""
#     print(mock_get)
#     # mock_get.side_effect
#     # mock_get.side_effect = asyncio.TimeoutError

#     with pytest.raises(asyncio.TimeoutError):
#         status, headers, r_json = await fetch()


@pytest.mark.asyncio
@patch('aiohttp.ClientSession.post')
async def test_post_fetch_example(mock_post: MagicMock):
    """Test fetch."""
    mock_post.return_value.__aenter__.return_value.status = 200
    mock_post.return_value.__aenter__.return_value.headers = {'content-type': 'text/html', 'Connection': 'close'}
    # mock the response.json method
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value={'key': 'value'})

    status, headers, r_json = await
    print("Debug jj",status, headers, r_json)
    assert status == 200
