#!/bin/bash

./userapi \
-redis_address=$REDIS_HOST \
-redis_max_retries=$REDIS_MAX_RETRIES \
-persist_path=$PERSISTING_DIR \
-data_ttl=$DATA_TTL
