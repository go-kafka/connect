Contributing Guide
------------------

Thank you for your interest in contributing! Improvements and reports of bugs
or documentation flaws are welcomed.

This is a small and simple project, but as my first public Go project I've
tried to establish good habits to follow. If you're unsure about anything,
don't hesitate to ask.

Please submit fixes or enhancements via GitHub pull requests, ensuring that
changes have passing test coverage and a clean bill of health from `gofmt`,
`golint`, and preferably `go vet`. The latter checks are not yet automated so
your diligence is appreciated.

If you wish to contribute a change that involves updating dependencies, please
use Glide in order to `glide update` or `glide get` the source in `vendor/`.
Use a version constraint, unless the package is imported via gopkg.in with v1
or above.
