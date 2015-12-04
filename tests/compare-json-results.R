#!/usr/bin/env Rscript
library("httr")
library("jsonlite")

#
# Compare JSON results from production and dev copy.
#

#
# Main
#
main <- function (args) {
  cat("DEBUG args", args, sep = "\n")
  # Local dev setup
  dev_url <- paste(
    Sys.getenv("ASPACE_PROTOCOL"),
    "://",
    Sys.getenv("ASPACE_HOST"),
    ":",
    Sys.getenv("ASPACE_PORT"), sep = "")
  
  dev_username <- Sys.getenv("ASPACE_USERNAME")
  dev_password <- Sys.getenv("ASPACE_PASSWORD")
  cat(dev_url, dev_username, dev_password, sep = "\n")
}
main(commandArgs(TRUE))
