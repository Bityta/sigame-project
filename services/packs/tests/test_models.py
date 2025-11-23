import pytest
from app.models import Pack, Question


def test_pack_validation():
    """Test Pack model validation"""
    pack = Pack(
        title="Test Pack",
        description="A test pack",
        difficulty="medium",
        category="general"
    )
    
    assert pack.title == "Test Pack"
    assert pack.difficulty == "medium"
    assert pack.category == "general"


def test_question_validation():
    """Test Question model validation"""
    question = Question(
        text="What is 2+2?",
        correct_answer="4",
        difficulty="easy"
    )
    
    assert question.text == "What is 2+2?"
    assert question.correct_answer == "4"
    assert question.difficulty == "easy"


def test_pack_with_questions():
    """Test Pack with related questions"""
    pack = Pack(
        title="Math Pack",
        description="Math questions",
        difficulty="easy",
        category="math"
    )
    
    question1 = Question(
        text="What is 1+1?",
        correct_answer="2",
        pack=pack
    )
    
    question2 = Question(
        text="What is 2+2?",
        correct_answer="4",
        pack=pack
    )
    
    assert question1.pack == pack
    assert question2.pack == pack

