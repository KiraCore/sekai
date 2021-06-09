#!/bin/bash

sekaid query customgov proposals --log_level=debug
sekaid query customgov proposals --reverse=false --log_level=debug
sekaid query customgov proposals --limit=1 --log_level=debug