from datetime import date, datetime
from typing import List, Union
from uuid import UUID

from pydantic import BaseModel

from pydantic_factories import ModelFactory


class TopicSchema(BaseModel):
    topic_id: UUID
    organization_id: UUID
    name: str
    description: str
    created_at: str
    created_by: str
    updated_at: str
    updated_by: str
    deleted_at: str | None
    deleted_by: str | None


class TopicFactory(ModelFactory):
    __model__ = TopicSchema


result = TopicFactory.build()
print(result.dict())
