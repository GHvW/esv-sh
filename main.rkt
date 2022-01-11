#lang racket/base

(require racket/cmdline)

(define verses-to-pull
  (command-line
   #:program "esv.sh"))