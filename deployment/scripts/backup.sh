#!/bin/bash

# Backup Script for PostgreSQL Databases
# Backs up all SIGame databases with rotation

set -e

# Configuration
BACKUP_DIR="/backups/postgresql"
RETENTION_DAYS=7
DATE=$(date +%Y%m%d_%H%M%S)

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo "=========================================="
echo "  PostgreSQL Backup - $(date)"
echo "=========================================="

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f /opt/sigame/.env.production ]; then
    export $(grep -v '^#' /opt/sigame/.env.production | xargs)
fi

# Function to backup a database
backup_database() {
    local db_name=$1
    local db_user=$2
    local db_password=$3
    local db_port=${4:-5432}
    local container=$5
    
    echo -e "${BLUE}Backing up $db_name...${NC}"
    
    BACKUP_FILE="$BACKUP_DIR/${db_name}_${DATE}.sql.gz"
    
    # Export password for pg_dump
    export PGPASSWORD="$db_password"
    
    # Backup using docker exec
    docker exec "$container" pg_dump -U "$db_user" "$db_name" | gzip > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        echo -e "${GREEN}✓ $db_name backup complete ($SIZE)${NC}"
    else
        echo -e "${RED}✗ Failed to backup $db_name${NC}"
        return 1
    fi
    
    unset PGPASSWORD
}

# Backup all databases
backup_database \
    "${AUTH_DB_NAME:-auth_db}" \
    "${AUTH_DB_USER:-authuser}" \
    "${AUTH_DB_PASSWORD:-authpass}" \
    5432 \
    "sigame-postgres-auth"

backup_database \
    "${LOBBY_DB_NAME:-lobby_db}" \
    "${LOBBY_DB_USER:-lobbyuser}" \
    "${LOBBY_DB_PASSWORD:-lobbypass}" \
    5433 \
    "sigame-postgres-lobby"

backup_database \
    "${PACKS_DB_NAME:-packs_db}" \
    "${PACKS_DB_USER:-packsuser}" \
    "${PACKS_DB_PASSWORD:-packspass}" \
    5434 \
    "sigame-postgres-packs"

backup_database \
    "${GAME_DB_NAME:-game_db}" \
    "${GAME_DB_USER:-gameuser}" \
    "${GAME_DB_PASSWORD:-gamepass}" \
    5435 \
    "sigame-postgres-game"

# Cleanup old backups
echo -e "${BLUE}Cleaning up old backups (older than $RETENTION_DAYS days)...${NC}"
find "$BACKUP_DIR" -name "*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
REMAINING=$(find "$BACKUP_DIR" -name "*.sql.gz" -type f | wc -l)
echo -e "${GREEN}✓ Cleanup complete ($REMAINING backups remaining)${NC}"

# Calculate total backup size
TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
echo ""
echo "=========================================="
echo -e "${GREEN}✓ Backup Complete!${NC}"
echo "=========================================="
echo "Total backup size: $TOTAL_SIZE"
echo "Retention: $RETENTION_DAYS days"
echo "=========================================="

