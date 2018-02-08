#!/bin/bash

>0.log && find . -name "*.go" | xargs cat >> 0.log
