#!/bin/bash

# util script to remove all scratch files

if [ $(basename "$PWD") = "testutil" ]; then
    find ./.. -type f -name "*_scratch*" -exec rm {} +
else
    find ./ -type f -name "*_scratch*" -exec rm {} +
fi