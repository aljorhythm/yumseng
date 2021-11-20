#!/bin/sh

formatted_files=$(go fmt ./...)
echo $formatted_files

[ -n "$formatted_files" ] && echo "formatting occured/detected please commit" && exit 1

exit 0