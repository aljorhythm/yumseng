#!/bin/bash

while true; do
  data=$(curl localhost/ping)

  echo "$data"

  tag=$(echo $data | jq ".Tag")
  echo tag "$tag" oldtag "$oldtag"

  if [ "$oldtag" != "" ] && [ "$oldtag" != "$tag" ]; then
    echo changed $tag
    osascript \
      -e "on run(argv)" \
      -e "return display alert \"Yumseng was deployed!\" message item 1 of argv" \
      -e "end" \
      -- "$tag"
  fi

  oldtag=$tag

  sleep 5
done
