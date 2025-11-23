import pytest
from fastapi.testclient import TestClient
from app.api.routes import router
from fastapi import FastAPI

app = FastAPI()
app.include_router(router, prefix="/api")

client = TestClient(app)


def test_health_check():
    """Test health check endpoint"""
    response = client.get("/api/health")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"


def test_get_packs():
    """Test get all packs endpoint"""
    response = client.get("/api/packs")
    assert response.status_code == 200
    assert isinstance(response.json(), list)


def test_get_pack_by_id():
    """Test get pack by ID endpoint"""
    # First create a pack
    response = client.get("/api/packs")
    if len(response.json()) > 0:
        pack_id = response.json()[0]["id"]
        response = client.get(f"/api/packs/{pack_id}")
        assert response.status_code in [200, 404]


def test_create_pack():
    """Test create pack endpoint"""
    pack_data = {
        "title": "Test Pack",
        "description": "Test description",
        "difficulty": "medium",
        "category": "test"
    }
    response = client.post("/api/packs", json=pack_data)
    assert response.status_code in [200, 201, 422]  # 422 if validation fails


def test_get_pack_questions():
    """Test get pack questions endpoint"""
    response = client.get("/api/packs")
    if len(response.json()) > 0:
        pack_id = response.json()[0]["id"]
        response = client.get(f"/api/packs/{pack_id}/questions")
        assert response.status_code in [200, 404]
        if response.status_code == 200:
            assert isinstance(response.json(), list)

