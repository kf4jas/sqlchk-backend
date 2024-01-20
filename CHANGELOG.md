# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

[Also based on](https://github.com/conventional-changelog/standard-version/blob/master/CHANGELOG.md) so decending.

## [0.1.4] - 2024-01-20
### Added
- adds compose info into the repo

## [0.1.3] - 2024-01-20
### Added
- adds sqlite3 code thru polymorphism
- adds the rest of the sqlite code for scanning tables and columns
- adds config file setup info to readme

### Removed
- cleans up web plugin

### Changed
- fixes missing openconn

## [0.1.2] - 2024-01-16
### Changed
- updates fmt of makefile

## [0.1.1] - 2024-01-13
### Added
- adds default target var to makefile
- adds playbooks and ansible code
- adds package building code for rpm and deb

## [0.1.0] - 2024-01-13
### Added
- adds changelog and versionupdater
- adds version file to makefile
- adds the data post uri
- adds comments for processrows
- adds deployment for local docker env
- adds db container and organizes container info
- adds mq system for better load balancing includes machinary

### Changed
- updates the order of ops for the dockerfile
- fixes the postgresql ssl setup