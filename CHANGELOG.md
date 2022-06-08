# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [8.0.0] - 2022-06-07
### Added
- Automatic file, line and package addition to error log when using `WithError`.

### Removed
- Handlers with dependencies; now encouraged to be separate packages.

### Changed
- Byte Pool now exposed through function only.
- Default error format function and output of wrapped error information.
- Tags & Types are now deduplicated in the default error format function.
- Updated to latest deps.
- CI to use GitHub Actions.
- Documentation.
- Default timestamp format to RFC3339Nano.
- Console logger uses builder pattern.


[Unreleased]: https://github.com/go-playground/log/compare/v8.0.0...HEAD
[8.0.0]: https://github.com/go-playground/log/compare/v7.0.2...v8.0.0