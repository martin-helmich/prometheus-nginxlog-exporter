# Contribution Guide

Thanks for considering to contribute to this project! :+1:

Do you want to...

- [Report a bug?](#reporting-bugs)
- [Fix a bug?](#fixing-bugs)
- [Add a new feature?](#adding-new-features)
- [Ask for support?](#asking-for-support)

## Reporting bugs

Before opening an issue for a bug, please follow the steps outlined in the [support document](./SUPPORT.md).

Should you not be able to resolve your issue on your own using the steps outlined there, please open a new issue with the ![label: bug][~bug] label. Also, make sure to supply all the information requested in the [support document](./SUPPORT.md).

## Fixing bugs

Pull requests that fix bugs are always welcome! To make sure that pull requests with bug fixes can be merged quickly and easily, please make sure of the following:

1. Make sure all your code is formatted by `go fmt`. Otherwise, the build pipeline will fail.
1. Make sure all existing behaviour tests pass. Otherwise, the build pipeline will fail.
1. If possible, try to add new test cases that verify the fixed bug (meaning that the test cases should fail without your bugfix, and pass with it).

If you want to provide a bugfix, but have questions on how to implement it, you could either...

1. Open up a new issue and ask for guidance on how to contribute
1. Open up a PR with whatever you've already got and ask for guidance on how to contribute whatever is missing

## Adding new features

Just as bugfixes, new features are also always welcome. However, for larger features, consider opening an issue for that feature first, so that it can be discussed. Use the ![label: enhancement][~enhancement] label for that.

When adding new features, please consider the following points:

1. All of [the points listed in "Fixing bugs"](#fixing-bugs)
1. Please make sure not to break backwards compatibility
1. Please also remember to document new behaviour in the README file

## Asking for support

See [SUPPORT.md](./SUPPORT.md).

[faq]: https://github.com/martin-helmich/prometheus-nginxlog-exporter#frequently-asked-questions
[pii]: https://en.wikipedia.org/w/index.php?title=Personally_identifiable_information&redirect=no
[~bug]: https://img.shields.io/badge/-bug-ee0701.svg
[~enhancement]: https://img.shields.io/badge/-enhancement-84b6eb.svg