#!/bin/bash
head -n 1 "$1" > "$2" && tail -n +2 "$1" | sort >> "$2"
