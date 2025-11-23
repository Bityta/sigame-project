"""Database connection module"""
import psycopg2
from psycopg2.extras import RealDictCursor
from typing import Optional
import logging

from app.config import settings

logger = logging.getLogger(__name__)


class Database:
    """PostgreSQL database connection wrapper"""
    
    def __init__(self):
        self.conn: Optional[psycopg2.extensions.connection] = None
    
    def connect(self):
        """Connect to PostgreSQL"""
        try:
            self.conn = psycopg2.connect(
                settings.get_postgres_connection_string(),
                cursor_factory=RealDictCursor
            )
            logger.info("✓ Connected to PostgreSQL")
        except Exception as e:
            logger.error(f"Failed to connect to PostgreSQL: {e}")
            raise
    
    def disconnect(self):
        """Disconnect from PostgreSQL"""
        if self.conn:
            self.conn.close()
            logger.info("Disconnected from PostgreSQL")
    
    def get_connection(self):
        """Get database connection"""
        if not self.conn or self.conn.closed:
            self.connect()
        return self.conn


# Global database instance
db = Database()

