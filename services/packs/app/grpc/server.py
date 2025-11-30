"""gRPC server implementation for Pack Service"""
import grpc
from concurrent import futures
import time
import logging

# Import generated proto files
try:
    from app.grpc import pack_pb2, pack_pb2_grpc
    GRPC_AVAILABLE = True
except ImportError as e:
    GRPC_AVAILABLE = False
    logging.warning(f"gRPC proto files not available: {e}")

from app.mock_data import get_pack_by_id, get_pack_content, pack_exists
from app.metrics import metrics

logger = logging.getLogger(__name__)


if GRPC_AVAILABLE:
    class PackServiceServicer(pack_pb2_grpc.PackServiceServicer):
        """gRPC Pack Service implementation"""
        
        def GetPackInfo(self, request, context):
            """Get pack information"""
            start_time = time.time()
            
            try:
                pack_id = request.pack_id
                pack = get_pack_by_id(pack_id)
                
                if not pack:
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details(f"Pack {pack_id} not found")
                    duration = time.time() - start_time
                    metrics.record_grpc_request("GetPackInfo", "NOT_FOUND", duration)
                    return pack_pb2.PackInfoResponse()
                
                duration = time.time() - start_time
                metrics.record_grpc_request("GetPackInfo", "OK", duration)
                
                return pack_pb2.PackInfoResponse(
                    found=True,
                    pack_id=pack["id"],
                    name=pack["name"],
                    author=pack["author"],
                    description=pack["description"],
                    rounds_count=pack["rounds_count"],
                    questions_count=pack["questions_count"],
                    has_media=False,
                    status="approved",
                    created_at=pack["created_at"]
                )
            except Exception as e:
                logger.error(f"Error in GetPackInfo: {e}")
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(str(e))
                duration = time.time() - start_time
                metrics.record_grpc_request("GetPackInfo", "ERROR", duration)
                return pack_pb2.PackInfoResponse()
        
        def GetPackContent(self, request, context):
            """Get full pack content"""
            start_time = time.time()
            
            try:
                pack_id = request.pack_id
                pack = get_pack_content(pack_id)
                
                if not pack:
                    context.set_code(grpc.StatusCode.NOT_FOUND)
                    context.set_details(f"Pack {pack_id} not found")
                    duration = time.time() - start_time
                    metrics.record_grpc_request("GetPackContent", "NOT_FOUND", duration)
                    return pack_pb2.PackContentResponse()
                
                # Convert pack data to protobuf format
                rounds = []
                for round_data in pack["rounds"]:
                    themes = []
                    for theme_data in round_data["themes"]:
                        questions = []
                        for q in theme_data["questions"]:
                            question = pack_pb2.Question(
                                id=q["id"],
                                price=q["price"],
                                text=q["text"],
                                answer=q["answer"],
                                media_type=q["media_type"]
                            )
                            questions.append(question)
                        
                        theme = pack_pb2.Theme(
                            id=theme_data["id"],
                            name=theme_data["name"],
                            questions=questions
                        )
                        themes.append(theme)
                    
                    round_msg = pack_pb2.Round(
                        id=round_data["id"],
                        round_number=round_data["round_number"],
                        name=round_data["name"],
                        themes=themes
                    )
                    rounds.append(round_msg)
                
                duration = time.time() - start_time
                metrics.record_grpc_request("GetPackContent", "OK", duration)
                
                return pack_pb2.PackContentResponse(
                    found=True,
                    pack_id=pack["id"],
                    name=pack["name"],
                    rounds=rounds
                )
            except Exception as e:
                logger.error(f"Error in GetPackContent: {e}")
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(str(e))
                duration = time.time() - start_time
                metrics.record_grpc_request("GetPackContent", "ERROR", duration)
                return pack_pb2.PackContentResponse()
        
        def ValidatePackExists(self, request, context):
            """Validate if pack exists"""
            start_time = time.time()
            
            try:
                pack_id = request.pack_id
                exists = pack_exists(pack_id)
                
                duration = time.time() - start_time
                metrics.record_grpc_request("ValidatePackExists", "OK", duration)
                
                return pack_pb2.ValidatePackResponse(
                    exists=exists,
                    is_owner=False,
                    status="approved" if exists else ""
                )
            except Exception as e:
                logger.error(f"Error in ValidatePackExists: {e}")
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(str(e))
                duration = time.time() - start_time
                metrics.record_grpc_request("ValidatePackExists", "ERROR", duration)
                return pack_pb2.ValidatePackResponse(exists=False, error=str(e))


def serve_grpc(port: int):
    """Start gRPC server"""
    if not GRPC_AVAILABLE:
        logger.error("Cannot start gRPC server: proto files not generated")
        return None
    
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pack_pb2_grpc.add_PackServiceServicer_to_server(PackServiceServicer(), server)
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logger.info(f"✓ gRPC server started on port {port}")
    return server

