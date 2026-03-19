#!/bin/bash

current_dir=$(pwd)
migrations=$(ls "$current_dir/cmd/migrations")

if [ $? -ne 0 ]; then
  echo "Cannot list migrations. Aborting..."
  exit 1
fi

last_idx=0

for migration in $migrations; do
  # Parse filename to get the number part
  migration_idx=$(echo "$migration" | cut -d '_' -f 1)
  migration_idx=$((10#$migration_idx))

  if [[ $(( migration_idx - last_idx )) -eq 0 ]]; then
    echo "Invalid migration detected. Last migration: $last_idx. Current migration $migration_idx"
    exit 1
  fi

  last_idx=$migration_idx
done

echo "No duplication found"
exit 0
