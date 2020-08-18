#!/bin/bash
docker build -t laiye/crtest .
docker run -p 8000:8000 laiye/crtest