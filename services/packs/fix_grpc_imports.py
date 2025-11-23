#!/usr/bin/env python3
"""Fix imports in generated gRPC files"""
import os
import sys

def fix_imports(filepath):
    """Fix absolute imports to relative imports in generated gRPC files"""
    print(f"Fixing imports in {filepath}")
    
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Replace absolute import with relative import
    content = content.replace(
        'import pack_pb2 as pack__pb2',
        'from . import pack_pb2 as pack__pb2'
    )
    
    with open(filepath, 'w') as f:
        f.write(content)
    
    print(f"✓ Fixed imports in {filepath}")

if __name__ == '__main__':
    grpc_dir = 'app/grpc'
    grpc_file = os.path.join(grpc_dir, 'pack_pb2_grpc.py')
    
    if os.path.exists(grpc_file):
        fix_imports(grpc_file)
        print("✓ gRPC imports fixed successfully")
    else:
        print(f"✗ File not found: {grpc_file}")
        sys.exit(1)

