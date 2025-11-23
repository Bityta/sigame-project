"""Pydantic models for Pack Service"""
from typing import List, Optional
from pydantic import BaseModel, Field
from datetime import datetime


class Question(BaseModel):
    """Question model"""
    id: str
    price: int
    text: str
    answer: str
    media_type: str = "text"


class Theme(BaseModel):
    """Theme model"""
    id: str
    name: str
    questions: List[Question]


class Round(BaseModel):
    """Round model"""
    id: str
    round_number: int
    name: str
    themes: List[Theme]


class PackInfo(BaseModel):
    """Pack information without full content"""
    id: str
    name: str
    author: str
    description: str
    rounds_count: int
    questions_count: int
    created_at: str


class PackContent(BaseModel):
    """Full pack content with all questions"""
    id: str
    name: str
    author: str
    description: str
    rounds_count: int
    questions_count: int
    created_at: str
    rounds: List[Round]


class PackListResponse(BaseModel):
    """Response for list of packs"""
    packs: List[PackInfo]
    total: int


class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    service: str
    version: str

