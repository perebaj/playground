from unittest.mock import patch, MagicMock, AsyncMock
import pytest
from xaiohttp import fetch
import asyncio

def test_hello():
    assert 1 == 1

@pytest.mark.asyncio
@patch('aiohttp.ClientSession.get')
async def test_fetch(mock_get: MagicMock):
    """Test fetch."""
    mock_get.return_value.__aenter__.return_value.status = 200
    mock_get.return_value.__aenter__.return_value.headers = {'content-type': 'text/html'}
    #mock the response.json method
    mock_get.return_value.__aenter__.return_value.json = AsyncMock(return_value={'key': 'value'})

    status, headers, r_json = await fetch()

    assert status == 200
    assert headers == 'text/html'
    assert r_json == {'key': 'value'}

# @pytest.mark.asyncio
# @patch('aiohttp.ClientSession.get')
# async def test_fetch_timeout(mock_get: AsyncMock):
#     """Test fetch."""
#     mock_get.__aenter__.get().__aenter__.side_effect = asyncio.TimeoutError

#     with pytest.raises(asyncio.TimeoutError):
#         status, headers, r_json = await fetch()
