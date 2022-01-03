#!/bin/bash

sekaid query customgov proposals
sekaid query customgov proposals --reverse
sekaid query customgov proposals --limit=1
sekaid query customgov proposals --limit=1 --output=json --reverse
