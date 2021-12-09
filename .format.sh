#!/bin/sh

formatted_files=$(go fmt ./...)
echo $formatted_files

[ -n "$formatted_files" ] && echo "formatting occured or detected please try again. run 'git add .' if changed files are not committed" && exit 1

exit 0