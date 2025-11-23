#!/bin/sh

# ==============================================
# SIGame 2.0 - Kafka Topics Initialization
# ==============================================

set -e

echo "Waiting for Kafka to be ready..."
sleep 10

KAFKA_BROKER="kafka:29092"

echo "Creating Kafka topics..."

# Topic: game.events
# Used for: ROOM_CREATED, PLAYER_JOINED, PLAYER_LEFT, ROOM_STARTED, ROOM_FINISHED, ROOM_CANCELLED
kafka-topics --bootstrap-server $KAFKA_BROKER \
    --create \
    --if-not-exists \
    --topic game.events \
    --partitions 3 \
    --replication-factor 1 \
    --config retention.ms=604800000 \
    --config compression.type=lz4

echo "✓ Created topic: game.events"

# Topic: game.actions
# Used for: BUTTON_PRESSED, QUESTION_SHOWN, ANSWER_GIVEN, ROUND_STARTED, ROUND_ENDED, GAME_FINISHED
kafka-topics --bootstrap-server $KAFKA_BROKER \
    --create \
    --if-not-exists \
    --topic game.actions \
    --partitions 3 \
    --replication-factor 1 \
    --config retention.ms=604800000 \
    --config compression.type=lz4

echo "✓ Created topic: game.actions"

# Topic: pack.events
# Used for: PACK_UPLOADED, PACK_PROCESSED, PACK_DOWNLOADED, PACK_RATED
kafka-topics --bootstrap-server $KAFKA_BROKER \
    --create \
    --if-not-exists \
    --topic pack.events \
    --partitions 3 \
    --replication-factor 1 \
    --config retention.ms=2592000000 \
    --config compression.type=lz4

echo "✓ Created topic: pack.events"

# Topic: notifications
# Used for: User notifications (game invites, game started, player joined/left, etc.)
kafka-topics --bootstrap-server $KAFKA_BROKER \
    --create \
    --if-not-exists \
    --topic notifications \
    --partitions 3 \
    --replication-factor 1 \
    --config retention.ms=259200000 \
    --config compression.type=lz4

echo "✓ Created topic: notifications"

echo ""
echo "Kafka topics created successfully!"
echo ""
echo "Listing all topics:"
kafka-topics --bootstrap-server $KAFKA_BROKER --list

echo ""
echo "Topic details:"
kafka-topics --bootstrap-server $KAFKA_BROKER --describe

