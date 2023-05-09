package main

import (
    "log"
    "os"
)

var (
    INFO = log.New(os.Stderr, "[RANDO] [INFO] ", 0)
    WARN = log.New(os.Stderr, "[RANDO] [WARNING] ", 0)
    ERR  = log.New(os.Stderr, "[RANDO] [ERROR] ", 0)
)
