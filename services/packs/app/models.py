"""Pydantic models for Pack Service"""
from typing import List, Optional
from pydantic import BaseModel, Field
from datetime import datetime


class Question(BaseModel):
    """Question model (DTO)
    
    Required: id, price, text, answer
    Optional: type, media_type, media_url, media_duration_ms
    """
    id: str = Field(..., description="Unique identifier of the question")
    price: int = Field(..., description="Question price/difficulty (points)", ge=0)
    text: str = Field(..., description="Question text", min_length=1)
    answer: str = Field(..., description="Correct answer", min_length=1)
    type: str = Field(default="normal", description="Question type: normal, secret, stake, forAll")
    media_type: str = Field(default="text", description="Media type: text, image, audio, video")
    media_url: Optional[str] = Field(default=None, description="URL to media file (image, audio, video)")
    media_duration_ms: int = Field(default=0, description="Media duration in milliseconds (for audio/video)")
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "q1",
                "price": 100,
                "text": "What is the capital of France?",
                "answer": "Paris",
                "type": "normal",
                "media_type": "text",
                "media_url": None,
                "media_duration_ms": 0
            }
        }


class Theme(BaseModel):
    """Theme model (DTO)
    
    All fields are required
    """
    id: str = Field(..., description="Unique identifier of the theme")
    name: str = Field(..., description="Theme name", min_length=1)
    questions: List[Question] = Field(..., description="List of questions in this theme")
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "theme1",
                "name": "Geography",
                "questions": []
            }
        }


class Round(BaseModel):
    """Round model (DTO)
    
    All fields are required
    """
    id: str = Field(..., description="Unique identifier of the round")
    round_number: int = Field(..., description="Round number", ge=1)
    name: str = Field(..., description="Round name", min_length=1)
    themes: List[Theme] = Field(..., description="List of themes in this round")
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "round1",
                "round_number": 1,
                "name": "First Round",
                "themes": []
            }
        }


class PackInfo(BaseModel):
    """Pack information without full content (DTO)
    
    All fields are required
    """
    id: str = Field(..., description="Unique identifier of the pack")
    name: str = Field(..., description="Pack name", min_length=1)
    author: str = Field(..., description="Pack author", min_length=1)
    description: str = Field(..., description="Pack description")
    rounds_count: int = Field(..., description="Total number of rounds", ge=0)
    questions_count: int = Field(..., description="Total number of questions", ge=0)
    created_at: str = Field(..., description="Creation timestamp (ISO 8601)")
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "pack1",
                "name": "General Knowledge",
                "author": "Admin",
                "description": "A pack about general knowledge",
                "rounds_count": 3,
                "questions_count": 30,
                "created_at": "2024-01-01T00:00:00"
            }
        }


class PackContent(BaseModel):
    """Full pack content with all questions (DTO)
    
    All fields are required
    """
    id: str = Field(..., description="Unique identifier of the pack")
    name: str = Field(..., description="Pack name", min_length=1)
    author: str = Field(..., description="Pack author", min_length=1)
    description: str = Field(..., description="Pack description")
    rounds_count: int = Field(..., description="Total number of rounds", ge=0)
    questions_count: int = Field(..., description="Total number of questions", ge=0)
    created_at: str = Field(..., description="Creation timestamp (ISO 8601)")
    rounds: List[Round] = Field(..., description="Complete list of rounds with content")
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "pack1",
                "name": "General Knowledge",
                "author": "Admin",
                "description": "A pack about general knowledge",
                "rounds_count": 3,
                "questions_count": 30,
                "created_at": "2024-01-01T00:00:00",
                "rounds": []
            }
        }


class PackListResponse(BaseModel):
    """Response for list of packs (DTO)
    
    All fields are required
    """
    packs: List[PackInfo] = Field(..., description="List of pack information")
    total: int = Field(..., description="Total number of packs", ge=0)
    
    class Config:
        json_schema_extra = {
            "example": {
                "packs": [],
                "total": 0
            }
        }


class HealthResponse(BaseModel):
    """Health check response (DTO)
    
    All fields are required
    """
    status: str = Field(..., description="Service status")
    service: str = Field(..., description="Service name")
    version: str = Field(..., description="Service version")
    
    class Config:
        json_schema_extra = {
            "example": {
                "status": "healthy",
                "service": "pack-service",
                "version": "1.0.0"
            }
        }

