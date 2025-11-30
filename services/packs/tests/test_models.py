import pytest
from app.models import PackInfo, PackContent, Question, Theme, Round, HealthResponse


def test_question_model():
    """Test Question model"""
    question = Question(
        id="q1",
        price=100,
        text="What is 2+2?",
        answer="4",
        media_type="text"
    )
    
    assert question.id == "q1"
    assert question.price == 100
    assert question.text == "What is 2+2?"
    assert question.answer == "4"
    assert question.media_type == "text"


def test_theme_model():
    """Test Theme model"""
    question = Question(
        id="q1",
        price=100,
        text="Test question",
        answer="Test answer"
    )
    
    theme = Theme(
        id="theme1",
        name="Geography",
        questions=[question]
    )
    
    assert theme.id == "theme1"
    assert theme.name == "Geography"
    assert len(theme.questions) == 1


def test_round_model():
    """Test Round model"""
    theme = Theme(
        id="theme1",
        name="Test Theme",
        questions=[]
    )
    
    round_obj = Round(
        id="round1",
        round_number=1,
        name="First Round",
        themes=[theme]
    )
    
    assert round_obj.id == "round1"
    assert round_obj.round_number == 1
    assert round_obj.name == "First Round"
    assert len(round_obj.themes) == 1


def test_pack_info_model():
    """Test PackInfo model"""
    pack = PackInfo(
        id="pack1",
        name="Test Pack",
        author="Test Author",
        description="Test description",
        rounds_count=3,
        questions_count=30,
        created_at="2024-01-01T00:00:00"
    )
    
    assert pack.id == "pack1"
    assert pack.name == "Test Pack"
    assert pack.author == "Test Author"
    assert pack.rounds_count == 3
    assert pack.questions_count == 30


def test_health_response_model():
    """Test HealthResponse model"""
    health = HealthResponse(
        status="healthy",
        service="pack-service",
        version="1.0.0"
    )
    
    assert health.status == "healthy"
    assert health.service == "pack-service"
    assert health.version == "1.0.0"
