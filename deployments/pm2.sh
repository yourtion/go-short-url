#!/usr/bin/env bash

set -e

pm2 stop short-url
mv short-url short-url.old
mv short-url.new short-url
pm2 start short-url

pm2 logs short-url