from unittest.mock import AsyncMock, patch

import pytest

from clearsale import Address, ClearSale, DocumentType, Transaction


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_authenticate_success(mock_post):
    mock_post.return_value.__aenter__.return_value.status = 200
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value={"token": "test_token"})
    mock_post.return_value.__aenter__.return_value.headers = {
        "accept": "application/json",
        "content-type": "application/*+json",
    }

    client = ClearSale(username="test_user", password="test_pass")
    response = await client.authenticate()

    assert response == "test_token"


@pytest.mark.asyncio
async def test_authenticate_failure():
    client = ClearSale(username="test_user", password="test_pass")

    with pytest.raises(ValueError) as exc_info:
        await client.authenticate()

    assert "Authentication failed" in str(exc_info.value)


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_create_transaction_success(mock_post):
    # Setup mock responses
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value={"id": "test_transaction_id"})

    # Create test transaction
    transaction = Transaction(
        documentType=DocumentType.CPF,
        document="12345678900",
        address=Address(
            zipCode="12345678",
            street="Test Street",
            number="123",
            district="Test District",
            city="Test City",
            state="TS",
            country="Test Country",
        ),
    )

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"  # Simulate already authenticated

    transaction_id = await client.create_transaction(transaction)

    # Verify the request was made with correct data
    mock_post.assert_called_once()
    call_args = mock_post.call_args
    assert call_args[1]["headers"]["Authorization"] == "Bearer existing_token"
    assert transaction_id == "test_transaction_id"


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_create_transaction_failure(mock_post):
    # Setup mock response for failure
    mock_post.return_value.__aenter__.return_value.status = 400
    mock_post.return_value.__aenter__.return_value.text = AsyncMock(return_value="Invalid data")

    transaction = Transaction(documentType=DocumentType.CPF, document="12345678900")
    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"

    with pytest.raises(ValueError) as exc_info:
        await client.create_transaction(transaction)

    assert "Transaction creation failed" in str(exc_info.value)


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_create_transaction_unauthorized_retry(mock_authenticate, mock_post):
    # Setup responses: first 401, then 201
    responses = [
        # First response: 401 Unauthorized
        AsyncMock(status=401),
        # Second response after re-auth: 201 Created
        AsyncMock(status=201, json=AsyncMock(return_value={"id": "new_transaction_id"})),
    ]
    mock_post.return_value.__aenter__.side_effect = responses

    # Mock authenticate to return a new token
    mock_authenticate.return_value = "new_token"

    transaction = Transaction(documentType=DocumentType.CPF, document="12345678900")
    client = ClearSale(username="test_user", password="test_pass")
    client.token = "expired_token"

    transaction_id = await client.create_transaction(transaction)

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    # Verify post was called twice
    assert mock_post.call_count == 2
    assert transaction_id == "new_transaction_id"


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_create_transaction_no_token(mock_authenticate, mock_post):
    # Setup responses
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value={"id": "test_transaction_id"})

    # Mock authenticate
    mock_authenticate.return_value = "new_auth_token"

    transaction = Transaction(documentType=DocumentType.CPF, document="12345678900")
    client = ClearSale(username="test_user", password="test_pass")
    # No token set initially

    transaction_id = await client.create_transaction(transaction)

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    assert transaction_id == "test_transaction_id"


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_ratings_success(mock_post):
    # Setup mock response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value=[{"ratingId": "123", "score": 85}])

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"  # Simulate already authenticated

    ratings = await client.ratings("test_transaction_id")

    # Verify the request was made with correct data
    mock_post.assert_called_once()
    call_args = mock_post.call_args
    assert call_args[1]["headers"]["Authorization"] == "Bearer existing_token"
    assert ratings == [{"ratingId": "123", "score": 85}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_ratings_failure(mock_post):
    # Setup mock response for failure
    mock_post.return_value.__aenter__.return_value.status = 404
    mock_post.return_value.__aenter__.return_value.text = AsyncMock(return_value="Transaction not found")

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"

    with pytest.raises(ValueError) as exc_info:
        await client.ratings("non_existent_id")

    assert "Failed to retrieve ratings" in str(exc_info.value)


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_ratings_unauthorized_retry(mock_authenticate, mock_post):
    # Setup responses: first 401, then 201
    responses = [
        # First response: 401 Unauthorized
        AsyncMock(status=401),
        # Second response after re-auth: 201 OK
        AsyncMock(status=201, json=AsyncMock(return_value=[{"ratingId": "456", "score": 90}])),
    ]
    mock_post.return_value.__aenter__.side_effect = responses

    # Mock authenticate to return a new token
    mock_authenticate.return_value = "new_token"

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "expired_token"

    ratings = await client.ratings("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    # Verify get was called twice
    assert mock_post.call_count == 2
    assert ratings == [{"ratingId": "456", "score": 90}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_ratings_no_token(mock_authenticate, mock_post):
    # Setup response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(return_value=[{"ratingId": "789", "score": 75}])

    # Mock authenticate
    mock_authenticate.return_value = "new_auth_token"

    client = ClearSale(username="test_user", password="test_pass")
    # No token set initially

    ratings = await client.ratings("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    assert ratings == [{"ratingId": "789", "score": 75}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_insights_success(mock_post):
    # Setup mock response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(
        return_value=[{"insightId": "123", "type": "ADDRESS_VERIFICATION", "result": "APPROVED"}]
    )

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"  # Simulate already authenticated

    insights = await client.insights("test_transaction_id")

    # Verify the request was made with correct data
    mock_post.assert_called_once()
    call_args = mock_post.call_args
    assert call_args[1]["headers"]["Authorization"] == "Bearer existing_token"
    assert insights == [{"insightId": "123", "type": "ADDRESS_VERIFICATION", "result": "APPROVED"}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_insights_failure(mock_post):
    # Setup mock response for failure
    mock_post.return_value.__aenter__.return_value.status = 404
    mock_post.return_value.__aenter__.return_value.text = AsyncMock(return_value="Transaction not found")

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"

    with pytest.raises(ValueError) as exc_info:
        await client.insights("non_existent_id")

    assert "Failed to retrieve insights" in str(exc_info.value)


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_insights_unauthorized_retry(mock_authenticate, mock_post):
    # Setup responses: first 401, then 201
    responses = [
        # First response: 401 Unauthorized
        AsyncMock(status=401),
        # Second response after re-auth: 201 OK
        AsyncMock(
            status=201,
            json=AsyncMock(
                return_value=[
                    {
                        "insightId": "456",
                        "type": "DOCUMENT_VERIFICATION",
                        "result": "APPROVED",
                    }
                ]
            ),
        ),
    ]
    mock_post.return_value.__aenter__.side_effect = responses

    # Mock authenticate to return a new token
    mock_authenticate.return_value = "new_token"

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "expired_token"

    insights = await client.insights("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    # Verify get was called twice
    assert mock_post.call_count == 2
    assert insights == [{"insightId": "456", "type": "DOCUMENT_VERIFICATION", "result": "APPROVED"}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_insights_no_token(mock_authenticate, mock_post):
    # Setup response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(
        return_value=[{"insightId": "789", "type": "PHONE_VERIFICATION", "result": "REJECTED"}]
    )

    # Mock authenticate
    mock_authenticate.return_value = "new_auth_token"

    client = ClearSale(username="test_user", password="test_pass")
    # No token set initially

    insights = await client.insights("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    assert insights == [{"insightId": "789", "type": "PHONE_VERIFICATION", "result": "REJECTED"}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_scores_success(mock_post):
    # Setup mock response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(
        return_value=[
            {
                "reason": "Initial",
                "value": 56.21,
                "createdAt": "2025-03-26T13:46:32.604Z",
            }
        ]
    )

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"

    scores = await client.scores("test_transaction_id")

    # Verify the request was made with correct data
    mock_post.assert_called_once()
    call_args = mock_post.call_args
    assert call_args[1]["headers"]["Authorization"] == "Bearer existing_token"
    assert scores == [{"reason": "Initial", "value": 56.21, "createdAt": "2025-03-26T13:46:32.604Z"}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
async def test_scores_failure(mock_post):
    # Setup mock response for failure
    mock_post.return_value.__aenter__.return_value.status = 404
    mock_post.return_value.__aenter__.return_value.text = AsyncMock(return_value="Transaction not found")

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "existing_token"

    with pytest.raises(ValueError) as exc_info:
        await client.scores("non_existent_id")

    assert "Failed to retrieve scores" in str(exc_info.value)


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_scores_unauthorized_retry(mock_authenticate, mock_post):
    # Setup responses: first 401, then 201
    responses = [
        # First response: 401 Unauthorized
        AsyncMock(status=401),
        # Second response after re-auth: 201 OK
        AsyncMock(
            status=201,
            json=AsyncMock(
                return_value=[
                    {
                        "reason": "Final",
                        "value": 78.45,
                        "createdAt": "2025-03-26T13:46:32.604Z",
                    }
                ]
            ),
        ),
    ]
    mock_post.return_value.__aenter__.side_effect = responses

    # Mock authenticate to return a new token
    mock_authenticate.return_value = "new_token"

    client = ClearSale(username="test_user", password="test_pass")
    client.token = "expired_token"

    scores = await client.scores("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    # Verify get was called twice
    assert mock_post.call_count == 2
    assert scores == [{"reason": "Final", "value": 78.45, "createdAt": "2025-03-26T13:46:32.604Z"}]


@pytest.mark.asyncio
@patch("aiohttp.ClientSession.post")
@patch("clearsale.ClearSale.authenticate")
async def test_scores_no_token(mock_authenticate, mock_post):
    # Setup response
    mock_post.return_value.__aenter__.return_value.status = 201
    mock_post.return_value.__aenter__.return_value.json = AsyncMock(
        return_value=[{"reason": "Final", "value": 78.45, "createdAt": "2025-03-26T13:46:32.604Z"}]
    )

    # Mock authenticate
    mock_authenticate.return_value = "new_auth_token"

    client = ClearSale(username="test_user", password="test_pass")
    # No token set initially

    scores = await client.scores("test_transaction_id")

    # Verify authenticate was called
    mock_authenticate.assert_called_once()
    assert scores == [{"reason": "Final", "value": 78.45, "createdAt": "2025-03-26T13:46:32.604Z"}]
