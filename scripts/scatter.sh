#!/bin/bash

FILE=`realpath $1`
RESPONSE=`curl -F file=@$FILE https://scrots.michaelphillips.dev/file`
FILENAME=`echo $RESPONSE | jq -r ".Name"`

echo "https://scrots.michaelphillips.dev/files/$FILENAME" | xclip -selection c
