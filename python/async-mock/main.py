class FraudPolicy:
    async def check_fraud_policy(self, probability):
        if probability > 0.5:
            return "Fraud detected"
        else:
            return "No fraud"

class FraudApi:
    async def read_current_month_tpv_values(self):
        # Simulate reading data
        return {"current_month_tpv": 10000}

    async def write_month_tpv_database(self, data):
        # Simulate writing to a database
        print(f"Writing TPV data to database: {data}")

    async def write_verdict_api_database(self, verdict):
        # Simulate writing a fraud verdict to an API
        print(f"Writing verdict to database: {verdict}")

class Processor:
    def __init__(self):
        self.fraud_policy = FraudPolicy()
        self.fraud_api = FraudApi()

    async def process_tpv(self):
        tpv_data = await self.fraud_api.read_current_month_tpv_values()
        verdict = await self.fraud_policy.check_fraud_policy(tpv_data["current_month_tpv"])
        await self.fraud_api.write_verdict_api_database(verdict)
