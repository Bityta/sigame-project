"""REST API routes for Pack Service"""
from fastapi import APIRouter, HTTPException, Response
from typing import List
import time

from app.models import PackInfo, PackContent, PackListResponse, HealthResponse
from app.mock_data import get_all_packs, get_pack_by_id, get_pack_content, pack_exists
from app.metrics import metrics
from app.config import settings

router = APIRouter()


@router.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        service=settings.service_name,
        version=settings.version
    )


@router.get("/api/packs", response_model=PackListResponse)
async def list_packs():
    """Get list of all available packs"""
    start_time = time.time()
    
    try:
        packs = get_all_packs()
        metrics.inc_pack_request("list")
        
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs", 200, duration)
        
        return PackListResponse(
            packs=packs,
            total=len(packs)
        )
    except Exception as e:
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs", 500, duration)
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/api/packs/{pack_id}", response_model=PackInfo)
async def get_pack_info(pack_id: str):
    """Get pack information by ID"""
    start_time = time.time()
    
    try:
        pack = get_pack_by_id(pack_id)
        
        if not pack:
            duration = time.time() - start_time
            metrics.record_http_request("GET", "/api/packs/{pack_id}", 404, duration)
            raise HTTPException(status_code=404, detail="Pack not found")
        
        metrics.inc_pack_request("info")
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs/{pack_id}", 200, duration)
        
        return pack
    except HTTPException:
        raise
    except Exception as e:
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs/{pack_id}", 500, duration)
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/api/packs/{pack_id}/content", response_model=PackContent)
async def get_pack_full_content(pack_id: str):
    """Get full pack content with all questions"""
    start_time = time.time()
    
    try:
        pack = get_pack_content(pack_id)
        
        if not pack:
            duration = time.time() - start_time
            metrics.record_http_request("GET", "/api/packs/{pack_id}/content", 404, duration)
            raise HTTPException(status_code=404, detail="Pack not found")
        
        metrics.inc_pack_request("content")
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs/{pack_id}/content", 200, duration)
        
        return pack
    except HTTPException:
        raise
    except Exception as e:
        duration = time.time() - start_time
        metrics.record_http_request("GET", "/api/packs/{pack_id}/content", 500, duration)
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/metrics")
async def get_metrics():
    """Prometheus metrics endpoint"""
    from prometheus_client import CONTENT_TYPE_LATEST
    metrics_data = metrics.generate_metrics()
    return Response(content=metrics_data, media_type=CONTENT_TYPE_LATEST)

