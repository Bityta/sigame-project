import pytest
from fastapi.testclient import TestClient
from app.api.routes import router
from fastapi import FastAPI

app = FastAPI()
app.include_router(router)

client = TestClient(app)


def test_health_check():
    """Test health check endpoint"""
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "healthy"
    assert "service" in data
    assert "version" in data


def test_get_packs_list():
    """Test get all packs endpoint"""
    response = client.get("/api/packs")
    assert response.status_code == 200
    data = response.json()
    assert "packs" in data
    assert "total" in data
    assert isinstance(data["packs"], list)


def test_get_pack_by_id():
    """Test get pack by ID endpoint"""
    # Get a pack from the list first
    response = client.get("/api/packs")
    packs = response.json()["packs"]
    
    if len(packs) > 0:
        pack_id = packs[0]["id"]
        response = client.get(f"/api/packs/{pack_id}")
        assert response.status_code == 200
        data = response.json()
        assert data["id"] == pack_id
        assert "name" in data
        assert "author" in data


def test_get_pack_not_found():
    """Test get non-existent pack"""
    response = client.get("/api/packs/non-existent-id")
    assert response.status_code == 404


def test_get_pack_content():
    """Test get pack content endpoint"""
    # Get a pack from the list first
    response = client.get("/api/packs")
    packs = response.json()["packs"]
    
    if len(packs) > 0:
        pack_id = packs[0]["id"]
        response = client.get(f"/api/packs/{pack_id}/content")
        assert response.status_code == 200
        data = response.json()
        assert data["id"] == pack_id
        assert "rounds" in data
