#!/bin/bash

for i in $(find . -name "*.go" | grep -v vendor); do
  echo $i
  sed -i '' -f license.tpl $i
done
