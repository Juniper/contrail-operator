#!/bin/bash

mount -o "lowerdir=/lib/modules,upperdir=/opt/modules,workdir=/opt/modules.wd" -t overlay overlay /lib/modules