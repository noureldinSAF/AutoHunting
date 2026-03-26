from typing import List
from pydantic import BaseModel


class TextRequest(BaseModel):
    text: str


class ResponseDetails(BaseModel):
    input_length: int
    processing_time_seconds: float
    valid_keys_count: int


class ApiKeyResponse(BaseModel):
    api_keys: List[str]
    total_count: int
    details: ResponseDetails
