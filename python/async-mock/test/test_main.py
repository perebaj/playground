import pytest
from unittest.mock import AsyncMock
from main import Processor

# def test_hello_world():
#     assert "hello world" == "hello world"

@pytest.mark.asyncio
async def test_processor_process_tpv():
    # Create a Processor instance
    processor = Processor()

    # Mock the FraudPolicy and FraudApi methods
    processor.fraud_policy.check_fraud_policy = AsyncMock(return_value="No fraud")
    processor.fraud_api.read_current_month_tpv_values = AsyncMock(return_value={"current_month_tpv": 10000})
    processor.fraud_api.write_verdict_api_database = AsyncMock()

    # Run the process_tpv method
    await processor.process_tpv()

    # Assert the methods were called with expected arguments
    processor.fraud_api.read_current_month_tpv_values.assert_called_once()
    processor.fraud_policy.check_fraud_policy.assert_called_once_with(10000)
    processor.fraud_api.write_verdict_api_database.assert_called_once_with("No fraud")
    assert processor.fraud_api.write_verdict_api_database.call_count == 1
    assert "jojo"   == "jojo"
