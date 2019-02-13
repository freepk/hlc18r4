#!/bin/bash
docker run --net host -v $(pwd)/tank:/var/loadtest -v $(pwd)/tmp:/tmp -it direvius/yandex-tank
