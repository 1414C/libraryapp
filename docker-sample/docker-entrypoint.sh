#!/bin/sh

# get the ipv4 address assigned to eth0
replace=$(ifconfig eth0 | grep "inet addr" | cut -d ':' -f 2 | cut -d ' ' -f 1)

# set a variable with the value we are planning to replace
search="xxx.xxx.xxx.xxx"

# check that variable replace has something in it
if [ -z "$replace" ]; then
  echo "Did not get an address value for eth0"
elif [ -n "$replace" ]; then
  echo "${replace} found"
# replace all instances of 'xxx.xxx.xxx.xxx' in .dev.config.json
# with the ipv4 address in the ${replace} variable
  sed -i "s/${search}/${replace}/g" .dev.config.json
  exec /main "$@"
fi
