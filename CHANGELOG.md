# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]
## [0.14.0] - 2022-05-26
### Added
- InvalidState

## [0.13.0] - 2022-05-26
### Added
- NoChange and Warn failures

## [0.12.0] - 2022-05-25
### Added
- OutOfRange failure
 
## [0.11.0] - 2022-05-17
### Added
- MissingFromContext, AlreadyExists failures
- tests for Panic

## [0.10.0] - 2022-05-11
### Added 
- Panic, RestAPI failures

## [0.9.1] - 2022-05-09
### Fixed
- IsMultiple returned false when wrapped.

## [0.9.0] - 2022-05-06
### Added
- support for multiple errors 
- startup failure

## [0.8.0] - 2022-05-02
### Added
- Timeout, IgnoreRetry

### Changed
- moved from deprecated github.com/pkg/errors to golang errors pkg

### Removed 
- InvalidInput until I can rework a better api

## [0.7.0] - 2022-04-21
### Added
- NotAuthorized
- NotAuthenticated
- Forbidden

## [0.6.2] - 2022-04-21
### Added
- input error now has InvalidInputMsg function
### Removed
- used InputOptions struct

## [0.6.1] - 2022-04-21
### Fixed
- input error now supports fields

## [0.6.0] - 2022-04-21
### Added
- bad request failure to be used in api and middleware code

## [0.5.0] - 2022-04-21
### Added
- shutdown error used when a server is signaled to shut down

## [0.4.0] - 2022-04-11 Added
### Added
- invalid param error was added 


## [0.3.0] - 2021-09-19
### Added
- config was added to support configuration errors
### Removed
- platform error type was removed, not needed

## [0.2.0] - 2021-08-21
### Added
- defer error for use inside defer functions

## [0.1.0] - 2021-08-21
- platform, system and server failures
- not found, validation, input and ignore failures