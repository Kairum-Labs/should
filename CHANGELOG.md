## [0.2.0] - 2025-09-13

### Features

- Add BeSorted assertion
- Add BeSameTime assertion for time comparisons
- Add BeError, BeErrorIs and BeErrorAs assertions for error
- Add formatted message support for assertions
- Add NotBeError assertion to verify no error occurs
- Add BeWithin assertion for float comparisons

### Refactor

- Consolidate assertion options into options.go
- Rename ContainFunc to AnyMatch
- Streamline BeEqual assertion for primitive types
- Remove actual value from BeOfType error output
- Update time msgs and type check in equality

### Documentation

- Update README with Go Reference badge and remove unused SliceOrArray constraint from types
- Enhance CONTRIBUTING.md with golangci-lint guidelines and usage instructions
- Update README and refactor assertion types
- Fix formatting of NotBeError example in web docs

### Miscellaneous Tasks

- Update golangci-lint configuration
- Add dupword linter and update README

## [0.1.0] - 2025-07-16

### Features

- Prepare for v0.1.0

### Refactor

- Rename IgnoreCase to WithIgnoreCase for consistency

### Miscellaneous Tasks

- Add CHANGELOG.md for project versioning and document initial release features

## [0.1.0-rc.5] - 2025-07-09

### Features

- Add StartsWith, EndsWith, ContainKey/Value and NotContainDuplicates
- Add NotBeEqual, NotContainKey, and NotContainValue assertions
- Introduce BeLessOrEqualThan assertion
- Add ContainSubstring assertion

### Bug Fixes

- CI

### Refactor

- Update assertion functions to use specific types and improve error handling
- Enhance and restructure CI compatibility workflow
- Simplify CI workflow
- Remove unnecessary generics and fix naming consistency in assertions
- Update CI workflow to include vulnerability checks
- Remove string support from ordered comparison assertions
- Organize and expand support for custom assertion messages
- Enhance error handling in assertions and improved fail function
- Streamline error handling in assertions using early return logic
- Rename BeNotEmpty and BeNotNil to NotBeEmpty and NotBeNil
- NotContainDuplicates
- Rename BeGreaterOrEqualThan/BeLessOrEqualThan to BeGreaterOrEqualTo/BeLessOrEqualTo

### Miscellaneous Tasks

- Update .gitignore and add compatibility workflow
- Update Go version requirements in documentation and workflows
- Optimize Go module and vulnerability checks

## [0.1.0-rc.4] - 2025-06-24

### Bug Fixes

- Numeric context for contain assertions

### Refactor

- Refactor assertions to use functional options API

### Documentation

- Update docs for assertion function signatures and features

## [0.1.0-rc.3] - 2025-06-21

### Features

- Add HaveLength, BeOfType, and BeOneOf assertions

### Bug Fixes

- CI
- Add t.Helper() to assertion functions

## [0.1.0-rc.2] - 2025-06-21

### Refactor

- Improve assertion logic, tests and update README
- Overhaul the API for improved clarity and flexibility

### Testing

- Improvement test coverage

## [0.1.0-rc.1] - 2025-06-15

### Features

- Adds tests, updates CI configs and general improvements
- Expose the public assertions API through should.go

### Bug Fixes

- CI and readme
