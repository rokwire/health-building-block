# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Add order number for the providers tests. [#25](https://github.com/rokwire/health-building-block/issues/25)

## [1.14.0] - 2020-09-08
### Security
- Admin APIs authentication update. [#20](https://github.com/rokwire/health-building-block/issues/20)

### Fixed
- Public health group quarantine error. [#22](https://github.com/rokwire/health-building-block/issues/22)

## [1.13.0] - 2020-08-31
### Fixed
- Provider group fails hitting the admin APIs. [#17](https://github.com/rokwire/health-building-block/issues/17)

## [1.12.0] - 2020-08-31
### Security
- Authorization improvements. [#14](https://github.com/rokwire/health-building-block/issues/14)

## [1.11.0] - 2020-08-27
### Added
- Audit trail. [#11](https://github.com/rokwire/health-building-block/issues/11)

## [1.10.1] - 2020-08-24
### Changed
- Update issue templates. [#5](https://github.com/rokwire/health-building-block/issues/5)
- Improve logs around the admin users authentication. [#7](https://github.com/rokwire/health-building-block/issues/7)

## [1.10.0] - 2020-08-14
### Changed
- Update Trace APIs. [#2](https://github.com/rokwire/health-building-block/issues/2)
- Update README. 

## [1.9.0] - 2020-08-12
### Added
- CODEOWNERS file. 
- Implement admin APIs for putting into and removing a user from quarantine. 

## [1.8.0] - 2020-08-10
### Fixed
- Concurrency issue. 
- Unable to set 'manual_test' to 'false' for Provider. 

## [1.7.0] - 2020-08-07
### Added
- Send firebase notification when ctest is added. 
- Mark the new created Shibboleth user as re-post = true. 
- Allow to turn off manual test reporting for a provider. 

## [1.6.0] - 2020-08-06
### Added
- Re-post test. 

## [1.5.0] - 2020-08-05
### Added
- Support different statuses for different apps/app versions. 

## [1.4.0] - 2020-07-31
### Changed
- Do not keep in the database already verified manual tests. 

## [1.3.0] - 2020-07-28
### Added
- Add code of conduct file. 

### Security
- Fix ECR image vulnerability. Reduced the image size. 

## [1.2.0] - 2020-07-24
### Changed
- Use Swagger for documentation instead of README file. 

### Security
- List of API keys support for the external APIs. 

## [1.1.0] - 2020-07-22
### Added
- Added changelog file. 
- Add license header to all source files. 

### Changed
- Cleanup unused/old entities. 

### Security
- Increase the level of protection for the user related APIs. 

## [1.0.632] - 2020-07-03
### Added
- Update access rule admin api.
- Delete access rule admin api.

[Unreleased]: https://github.com/rokwire/health-building-block/compare/v1.14.0...HEAD
[1.14.0]: https://github.com/rokwire/health-building-block/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/rokwire/health-building-block/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/rokwire/health-building-block/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/rokwire/health-building-block/compare/v1.10.1...v1.11.0
[1.10.1]: https://github.com/rokwire/health-building-block/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/rokwire/health-building-block/compare/v1.9.0...v1.10.0
