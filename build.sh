#!/bin/bash
docker build -t mimic-builder -f Dockerfile.build . 

docker run -v $(pwd)/output:/output mimic-builder  
