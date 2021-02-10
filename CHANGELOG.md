# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Removed

## [1.1.0] - 2021-02-10

### Added

- `cern_egroup_members` adds a `mails` attribute to fetch the list of e-mail addresses associated with the users in the e-group. The flag `query_mails` controls whether this information should be fetched or not.

## [1.0.0] - 2020-12-11

### Added

- New resources to manage LanDB elements: vm, vm card and vm interface.
- Migrate data source to get Teigi secrets from [the previous dedicated provider](https://gitlab.cern.ch/batch-team/infra/terraform-provider-teigi).

## [0.1.0] - 2019-12-11

### Added

- New data source to get egroup members.
