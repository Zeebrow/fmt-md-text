#!/usr/bin/env sh
# dont run these

cat README.md | ./fmt-md-text > /dev/null
[ "$?" -eq 0 ] && echo "ye"
cat README.md | ./fmt-md-text -l > /dev/null
[ "$?" -eq 0 ] && echo "ye"

./fmt-md-text -f README.md > /dev/null 
[ "$?" -eq 0 ] && echo "ye"
./fmt-md-text -f README.md -l > /dev/null
[ "$?" -eq 0 ] && echo "ye"

./fmt-md-text -f asdf > /dev/null
[ "$?" -eq 1 ] && echo "ye"

