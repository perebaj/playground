from enum import Enum
from typing import List, Optional

import aiohttp
from pydantic import BaseModel, Field


class DocumentType(str, Enum):
    """Document type for clearSale API."""

    CPF = "CPF"
    CNPJ = "CNPJ"


class Address(BaseModel):
    """Address model for clearSale API."""

    zip_code: str = Field(..., alias="zipCode")
    street: Optional[str] = None
    number: Optional[str] = None
    complement: Optional[str] = None
    district: Optional[str] = None
    city: Optional[str] = None
    state: Optional[str] = None
    country: Optional[str] = None


class Phone(BaseModel):
    """Phone model for clearSale API."""

    country_code: int = Field(..., alias="countryCode")
    area_code: int = Field(..., alias="areaCode")
    number: int
    type: Optional[str] = None


class Transaction(BaseModel):
    """Transaction model for clearSale API."""

    document_type: DocumentType = Field(..., alias="documentType")
    document: str
    # Address is an optional field to create a transaction, but if the intention is to use ratings or insights, the address is required
    address: Optional[Address] = None
    # Phone is an optional field to create a transaction, but if the intention is to use ratings or insights, the phone is required
    phone: Optional[Phone] = None
    # Email is an optional field to create a transaction, but if the intention is to use ratings or insights, the email is required
    email: Optional[str] = None
    # secondary_document can't have the same value as the document field, and secondary_document_type is required if secondary_document is informed
    # If document_type is CPF, secondary_document_type must be CNPJ and vice versa
    secondary_document_type: Optional[DocumentType] = Field(
        None, alias="secondaryDocumentType"
    )
    secondary_document: Optional[str] = Field(None, alias="secondaryDocument")


class ClearSale:
    """
    The documentation for the clearSale API can be found at: https://devs.plataformadatatrust.clearsale.com.br/reference
    """

    # TODO: Read this base url from the settings not a hardcoded value
    BASE_URL = "https://datatrustapihml.clearsale.com.br/"

    def __init__(
        self, username: str = None, password: str = None, base_url: str = None
    ):
        self.username = username
        self.password = password
        # if base_url is not provided, use the default one(staging environment)
        self.base_url = base_url if base_url else self.BASE_URL
        self.token = None

    async def authenticate(self) -> str:
        """Authenticate and retrieve access token.

        An important behavior of the API: If we consult the API it will not expire the current token or even create a new one if the current still valid.
        It will just create a new token if the current one is invalid.

        Ref: https://devs.plataformadatatrust.clearsale.com.br/reference/post_v1-authentication
        """
        url = f"{self.base_url}v1/authentication"

        async with aiohttp.ClientSession() as session:
            payload = {"Username": self.username, "Password": self.password}

            async with session.post(
                url,
                json=payload,
                headers={
                    "accept": "application/json",
                    "content-type": "application/*+json",
                },
            ) as response:
                if response.status != 200:
                    text = await response.text()
                    raise ValueError(f"Authentication failed: {text}")

                data = await response.json()
                self.token = data["token"]
                return self.token

    async def create_transaction(self, transaction: Transaction) -> str:
        """
        Create a new transaction.
        https://devs.plataformadatatrust.clearsale.com.br/reference/post_v1-transaction
        """
        url = f"{self.base_url}v1/transaction"

        if not self.token:
            await self.authenticate()

        async with aiohttp.ClientSession() as session:
            async with session.post(
                url,
                json=transaction.dict(by_alias=True),
                headers={
                    "accept": "application/json",
                    "content-type": "application/*+json",
                    "Authorization": f"Bearer {self.token}",
                },
            ) as response:
                if response.status == 401:  # Unauthorized - token expired
                    await self.authenticate()
                    # Retry with new token
                    return await self.create_transaction(transaction)

                if response.status != 201:
                    text = await response.text()
                    raise ValueError(f"Transaction creation failed: {text}")

                data = await response.json()
                return data["id"]

    async def ratings(self, transaction_id: str) -> List[dict]:
        """
        Retrieve ratings for a transaction.
        Ref: https://devs.plataformadatatrust.clearsale.com.br/reference/post_id-ratings
        """
        url = f"{self.base_url}v1/transaction/{transaction_id}/ratings"

        if not self.token:
            await self.authenticate()

        async with aiohttp.ClientSession() as session:
            async with session.post(
                url, headers={"Authorization": f"Bearer {self.token}"}
            ) as response:
                if response.status == 401:  # Unauthorized - token expired
                    await self.authenticate()
                    # Retry with new token
                    return await self.ratings(transaction_id)

                if response.status != 201:
                    text = await response.text()
                    raise ValueError(f"Failed to retrieve ratings: {text}")

                return await response.json()

    async def insights(self, transaction_id: str) -> List[dict]:
        """
        Retrieve insights for a transaction.
        Ref: https://devs.plataformadatatrust.clearsale.com.br/reference/post_id-insights
        """
        url = f"{self.base_url}v1/transaction/{transaction_id}/insights"

        if not self.token:
            await self.authenticate()

        async with aiohttp.ClientSession() as session:
            async with session.post(
                url,
                headers={
                    "accept": "application/json",
                    "Authorization": f"Bearer {self.token}",
                },
            ) as response:
                if response.status == 401:  # Unauthorized - token expired
                    await self.authenticate()
                    # Retry with new token
                    return await self.insights(transaction_id)

                if response.status != 201:
                    text = await response.text()
                    raise ValueError(f"Failed to retrieve insights: {text}")

                return await response.json()

    async def scores(self, transaction_id: str) -> List[dict]:
        """
        Retrieve scores for a transaction.
        Ref: https://devs.plataformadatatrust.clearsale.com.br/reference/post_v1-transaction-id-scores
        """
        url = f"{self.base_url}v1/transaction/{transaction_id}/scores"

        if not self.token:
            await self.authenticate()

        async with aiohttp.ClientSession() as session:
            async with session.post(
                url,
                headers={
                    "accept": "application/json",
                    "Authorization": f"Bearer {self.token}",
                },
            ) as response:
                if response.status == 401:
                    await self.authenticate()
                    # Retry with new token
                    return await self.insights(transaction_id)

                if response.status != 201:
                    text = await response.text()
                    raise ValueError(f"Failed to retrieve scores: {text}")

                return await response.json()
