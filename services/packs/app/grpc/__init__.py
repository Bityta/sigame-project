"""gRPC service module"""

# Re-export generated proto modules for easier imports
try:
    from .pack_pb2 import *
    from .pack_pb2_grpc import *
    __all__ = ['pack_pb2', 'pack_pb2_grpc']
except ImportError:
    pass

