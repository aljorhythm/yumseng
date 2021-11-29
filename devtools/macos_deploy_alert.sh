#!/bin/bash

while true; do
  data=$(curl yumseng.herokuapp.com/ping)

  echo "$data"

  if [ "$olddata" != "" ] && [ "$olddata" != $data ]; then
    osascript -e 'yumseng was deployed"'"$data"'"'
  fi
  olddata=$data
  sleep 5
done
