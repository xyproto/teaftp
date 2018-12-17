#!/bin/sh
if [ -z $1 ]; then
  echo 'Please provide a directory as the first argument.'
  exit 1
fi
echo 'Starting TeaFTP server...'
echo 'You can try retrieving files with: curl tftp://localhost/srv/example.txt'
echo "Serving $1 as the TFTP path /srv"
docker run -i -t --rm -v "$(realpath $1):/srv" --net=host teaftp
