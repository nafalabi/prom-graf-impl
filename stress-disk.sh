#!/bin/sh

ALLOCATION_SIZE=$1

if [ -z "$ALLOCATION_SIZE" ]; then
  echo
  echo "Invalid Usage..."
  echo
  echo "Usage:"
  echo
  echo "  sh stress-disk.sh <allocation_size>"
  echo "  example: sh stress-disk.sh 20G"
  echo
  echo
  exit
fi

echo "  Allocating 20G..."
echo

fio --name=test --rw=write --bs=1M --iodepth=32 --filename=/tmp/test --size=$ALLOCATION_SIZE

