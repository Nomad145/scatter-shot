#!/bin/bash

FILE=`realpath $1`
RESPONSE=`curl -F file=@$FILE https://files.michaelphillips.dev/upload`
FILENAME=`echo $RESPONSE | jq -r ".Name"`

echo "https://files.michaelphillips.dev/view/$FILENAME" | xclip -selection c
