# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [8.1.2] - 2023-08-16
### Fixed
- errors dependency which contains fixes for wrapped/wrapping errors.

## [8.1.1] - 2023-08-16
### Fixed
- errors.Link output in error function after updating dependency.

## [8.1.0] - 2023-08-15
### Added
- log.G as shorthand for adding a set of Grouped fields. This ability has always been present but is now fully supported in the default logger and with helper function for ease of use.
- slog support added in Go 1.21+ both to use as an slog.Handler or redirect.

### Fixed
- errors.Chain handling from default withErrorFn handler after dep upgrade.

## [8.0.2] - 2023-06-22
### Fixed
- Corrected removal of default logger upon registering a custom one.

## [8.0.1] - 2022-06-23
### Fixed
- Handling un-hashable tag values during dedupe process.

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
- Removed colors from built in console logger.
- Removed ability to remove individual log levels externally; RemoveHandler+AddHandler can do the same.


[Unreleased]: https://github.com/go-playground/log/compare/v8.1.2...HEAD
[8.1.2]: https://github.com/go-playground/log/compare/v8.1.1...v8.1.2
[8.1.1]: https://github.com/go-playground/log/compare/v8.1.0...v8.1.1
[8.1.0]: https://github.com/go-playground/log/compare/v8.0.2...v8.1.0
[8.0.2]: https://github.com/go-playground/log/compare/v8.0.1...v8.0.2
[8.0.1]: https://github.com/go-playground/log/compare/v8.0.0...v8.0.1
[8.0.0]: https://github.com/go-playground/log/compare/v7.0.2...v8.0.0