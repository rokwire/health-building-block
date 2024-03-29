# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.13.0] - 2021-10-05
### Changed
- Add "exempt" changes to UINOverride entity [#152](https://github.com/rokwire/health-building-block/issues/152)

## [2.12.1] - 2021-09-01
### Added
- Add activation time to UINOverride record [#147](https://github.com/rokwire/health-building-block/issues/147)

## [2.12.0] - 2021-08-04
### Security
- Upgrade mongo driver. [#144](https://github.com/rokwire/health-building-block/issues/144)

## [2.11.0] - 2021-07-29
### Security
- Update "jwt-go" library. [#141](https://github.com/rokwire/health-building-block/issues/141)

## [2.10.0] - 2021-07-20
### Fixed
- Inconsistent get user API return. [#136](https://github.com/rokwire/health-building-block/issues/136)

## [2.8.0] - 2021-06-16
- Add "consent_vaccine" flag [#133](https://github.com/rokwire/health-building-block/issues/133)

## [2.7.0] - 2021-04-26
### Added
= Create Create a Security.md [#150] (https://github.com/rokwire/health-building-block/issues/150)
- Expose external API to check if the user exists. [#128](https://github.com/rokwire/health-building-block/issues/128)
- Expose get time client API. [#130](https://github.com/rokwire/health-building-block/issues/130)

### Changed
- Get email and phone from the Rokmetro system. [#126](https://github.com/rokwire/health-building-block/issues/126)
- Make the log files on AWS readable. [#122](https://github.com/rokwire/health-building-block/issues/122)

## [2.6.0] - 2021-02-26
### Changed
- Add date created field to the app version entity. [#117](https://github.com/rokwire/health-building-block/issues/117)

### Added
- Expose get configs admin API. [#119](https://github.com/rokwire/health-building-block/issues/119)
- Expose APIs for external join group approvement. [#121](https://github.com/rokwire/health-building-block/issues/121)

## [2.5.0] - 2021-01-20
### Changed
- Update the process manual test admin API. [#114](https://github.com/rokwire/health-building-block/issues/114)

## [2.4.0] - 2021-01-11
### Security
- Prepare the APIs to be consumed by web app. [#111](https://github.com/rokwire/health-building-block/issues/111)

## [2.3.0] - 2020-12-23
### Security
- New authentication mechanism. [#88](https://github.com/rokwire/health-building-block/issues/88)

## [2.2.0] - 2020-12-16
### Changed
- Add encrypted PK field to the user entity. [#106](https://github.com/rokwire/health-building-block/issues/106)

## [2.1.0] - 2020-12-11
### Changed
- Sub accounts feature internal improvements. [#103](https://github.com/rokwire/health-building-block/issues/103)
- Update sub accounts required fields. [#101](https://github.com/rokwire/health-building-block/issues/101)

## [2.0.0] - 2020-12-10
### Added
- Sub accounts feature. [#92](https://github.com/rokwire/health-building-block/issues/92)

## [1.34.0] - 2020-12-02
### Changed
- Allow to search rosters by substring. [#96](https://github.com/rokwire/health-building-block/issues/96)

## [1.33.0] - 2020-12-01
### Changed
- Rosters update. [#93](https://github.com/rokwire/health-building-block/issues/93)

## [1.32.0] - 2020-11-19
### Added
- Implement admin API for updating single roster. [#89](https://github.com/rokwire/health-building-block/issues/89)

## [1.31.0] - 2020-11-16
### Changed
- Capitol Staff update. [#85](https://github.com/rokwire/health-building-block/issues/85)

## [1.30.0] - 2020-11-04
### Changed
- Enhance Roster APIs. [#82](https://github.com/rokwire/health-building-block/issues/82)

## [1.29.0] - 2020-10-27
### Fixed
- Not able to delete county status. [#77](https://github.com/rokwire/health-building-block/issues/77)

## [1.28.0] - 2020-10-23
### Changed
- Get audit API update. [#74](https://github.com/rokwire/health-building-block/issues/74)

## [1.27.0] - 2020-10-14
### Security
- Prepare the admin APIs for the web application. [#71](https://github.com/rokwire/health-building-block/issues/71)

## [1.26.0] - 2020-10-08
### Changed
- UIN Override changes. [#68](https://github.com/rokwire/health-building-block/issues/68)

## [1.25.0] - 2020-10-05
### Added
- UIN override extentions. [#47](https://github.com/rokwire/health-building-block/issues/47)

## [1.24.0] - 2020-10-01
### Added
- App versions handling improvements. [#63](https://github.com/rokwire/health-building-block/issues/63)

## [1.23.0] - 2020-09-30
### Added
- Add 2.7 and 2.8 as supported app versions. [#60](https://github.com/rokwire/health-building-block/issues/60)

## [1.22.0] - 2020-09-29
### Fixed
- Fix timezone database Docker issue. [#57](https://github.com/rokwire/health-building-block/issues/57)

## [1.21.0] - 2020-09-29
### Added
- Reset wait times to gray when locations close (are closed). [#52](https://github.com/rokwire/health-building-block/issues/52)

### Changed
- Admin APIs update - audit data. [#43](https://github.com/rokwire/health-building-block/issues/43)

### Security
- Allow location admins to retrieve providers, counties and test types. [#53](https://github.com/rokwire/health-building-block/issues/53)

## [1.20.0] - 2020-09-25
### Added
- Building access APIs. [#48](https://github.com/rokwire/health-building-block/issues/48)

## [1.19.0] - 2020-09-24
### Added
- UIN overrides. [#44](https://github.com/rokwire/health-building-block/issues/44)

### Changed
- Audit access contol update. [#41](https://github.com/rokwire/health-building-block/issues/41)

## [1.18.0] - 2020-09-18
### Added
- Audit trail client data. [#38](https://github.com/rokwire/health-building-block/issues/38)

## [1.17.0] - 2020-09-17
### Added
- Test location wait time. [#35](https://github.com/rokwire/health-building-block/issues/35)

## [1.16.0] - 2020-09-16
### Added
- Expose Rules APIs. [#28](https://github.com/rokwire/health-building-block/issues/28)

### Fixed
- Audit for "action" entity. [#30](https://github.com/rokwire/health-building-block/issues/30)

### Changed
- Disable test type result validation fields. [#32](https://github.com/rokwire/health-building-block/issues/32)

## [1.15.0] - 2020-09-10
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

[Unreleased]: https://github.com/rokwire/health-building-block/compare/v2.13.0...HEAD
[2.13.0]: https://github.com/rokwire/health-building-block/compare/v2.12.1...v2.13.0
[2.12.1]: https://github.com/rokwire/health-building-block/compare/v2.12.0...v2.12.1
[2.12.0]: https://github.com/rokwire/health-building-block/compare/v2.11.0...v2.12.0
[2.11.0]: https://github.com/rokwire/health-building-block/compare/v2.10.0...v2.11.0
[2.10.0]: https://github.com/rokwire/health-building-block/compare/v2.9.0...v2.10.0
[2.9.0]: https://github.com/rokwire/health-building-block/compare/v2.8.0...v2.9.0
[2.8.0]: https://github.com/rokwire/health-building-block/compare/v2.7.0...v2.8.0
[2.7.0]: https://github.com/rokwire/health-building-block/compare/v2.6.0...v2.7.0
[2.6.0]: https://github.com/rokwire/health-building-block/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/rokwire/health-building-block/compare/v2.4.0...v2.5.0
[2.4.0]: https://github.com/rokwire/health-building-block/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/rokwire/health-building-block/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/rokwire/health-building-block/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/rokwire/health-building-block/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/rokwire/health-building-block/compare/v1.34.0...v2.0.0
[1.34.0]: https://github.com/rokwire/health-building-block/compare/v1.33.0...v1.34.0
[1.33.0]: https://github.com/rokwire/health-building-block/compare/v1.32.0...v1.33.0
[1.32.0]: https://github.com/rokwire/health-building-block/compare/v1.31.0...v1.32.0
[1.31.0]: https://github.com/rokwire/health-building-block/compare/v1.30.0...v1.31.0
[1.30.0]: https://github.com/rokwire/health-building-block/compare/v1.29.0...v1.30.0
[1.29.0]: https://github.com/rokwire/health-building-block/compare/v1.28.0...v1.29.0
[1.28.0]: https://github.com/rokwire/health-building-block/compare/v1.27.0...v1.28.0
[1.27.0]: https://github.com/rokwire/health-building-block/compare/v1.26.0...v1.27.0
[1.26.0]: https://github.com/rokwire/health-building-block/compare/v1.25.0...v1.26.0
[1.25.0]: https://github.com/rokwire/health-building-block/compare/v1.24.0...v1.25.0
[1.24.0]: https://github.com/rokwire/health-building-block/compare/v1.23.0...v1.24.0
[1.23.0]: https://github.com/rokwire/health-building-block/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/rokwire/health-building-block/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/rokwire/health-building-block/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/rokwire/health-building-block/compare/v1.19.0...v1.20.0
[1.19.0]: https://github.com/rokwire/health-building-block/compare/v1.18.0...v1.19.0
[1.18.0]: https://github.com/rokwire/health-building-block/compare/v1.17.0...v1.18.0
[1.17.0]: https://github.com/rokwire/health-building-block/compare/v1.16.0...v1.17.0
[1.16.0]: https://github.com/rokwire/health-building-block/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/rokwire/health-building-block/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/rokwire/health-building-block/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/rokwire/health-building-block/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/rokwire/health-building-block/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/rokwire/health-building-block/compare/v1.10.1...v1.11.0
[1.10.1]: https://github.com/rokwire/health-building-block/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/rokwire/health-building-block/compare/v1.9.0...v1.10.0
