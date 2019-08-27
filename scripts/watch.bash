#!/bin/bash

make test-internal ; fswatch . -e ".*" -i "\\.go$" | xargs -n1 -I{}  make clear test-internal
