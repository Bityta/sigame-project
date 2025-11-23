"""Configuration module for Pack Service"""
import os
from typing import Optional
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings loaded from environment variables"""
    
    # PostgreSQL
    postgres_host: str = os.getenv("POSTGRES_HOST", "localhost")
    postgres_port: int = int(os.getenv("POSTGRES_PORT", "5434"))
    postgres_user: str = os.getenv("POSTGRES_USER", "packsuser")
    postgres_password: str = os.getenv("POSTGRES_PASSWORD", "packspass")
    postgres_db: str = os.getenv("POSTGRES_DB", "packs_db")
    postgres_sslmode: str = os.getenv("POSTGRES_SSLMODE", "disable")
    
    # Server
    http_port: int = int(os.getenv("HTTP_PORT", "8005"))
    grpc_port: int = int(os.getenv("GRPC_PORT", "50055"))
    
    # Service info
    service_name: str = "pack-service"
    version: str = "0.1.0"
    
    def get_postgres_connection_string(self) -> str:
        """Get PostgreSQL connection string"""
        return (
            f"host={self.postgres_host} "
            f"port={self.postgres_port} "
            f"dbname={self.postgres_db} "
            f"user={self.postgres_user} "
            f"password={self.postgres_password} "
            f"sslmode={self.postgres_sslmode}"
        )
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


# Global settings instance
settings = Settings()

