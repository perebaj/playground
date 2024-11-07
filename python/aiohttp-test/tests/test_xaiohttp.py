from unittest.mock import patch, MagicMock
import pytest
from xaiohttp import fetch


def test_hello():
    assert 1 == 1

@pytest.mark.asyncio
@patch('aiohttp.ClientSession.get')
async def test_fetch(mock_get: MagicMock):
    """Test fetch."""
    mock_get.return_value.__aenter__.return_value.status = 200
    mock_get.return_value.__aenter__.return_value.headers = {'content-type': 'text/html'}

    status, headers = await fetch()

    assert status == 200
    assert headers == 'text/html'
