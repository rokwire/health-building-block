# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Update README. 

## [1.9.0] - 2020-08-12
### Added
- CODEONWERS file. 
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

[Unreleased]: https://github.com/rokwire/health-building-block/compare/v1.9.0...HEAD