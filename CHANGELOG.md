# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.2] - 2022-08-27
- Fix: change broadcast to no longer require recovery from edge case panics.

## [0.4.1] - 2022-08-25
- Fix: allow blocking writes to output channels and allow recovery from edge case panics.

## [0.4.0] - 2022-08-24
### Changed
- Minor change to the underlying `select` statement in `Coalesce` and `Broadcast` to make it more idiomatic.

## [0.3.0] - 2022-08-22
### Changed
- Renamed `GenerateOutputBroadcasters` to `GenerateOutputBroadcast`
- Renamed `GenerateInputOutputBroadcasters` to `GenerateBroadcast`

## [0.2.0] - 2022-08-21
### Added
- Examples in go documentation for [pkg.go.dev/github.com/JamesChung/fofi](https://pkg.go.dev/github.com/JamesChung/fofi)
- `examples/` directory with runnable examples of each function
- New `GenerateInputOutputBroadcasters` function

### Changed
- Renamed `GenerateBroadcasters` to `GenerateOutputBroadcasters`

## [0.1.0] - 2022-08-20
### Added
- CHANGELOG.md
- Initial release code
