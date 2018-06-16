# technical

[![CircleCI](https://circleci.com/gh/blacklabcapital/technical.svg?style=svg)](https://circleci.com/gh/blacklabcapital/technical)
[![GoDoc](https://godoc.org/github.com/blacklabcapital/technical?status.svg)](https://godoc.org/github.com/blacklabcapital/technical)

## Description

`technical` is a Go package implementing various technical trading indicators and statistical functions.

It aims to provide simple, performant, open-source implementations of commonly used mathematical and statistical functions and indicators used in technical and quantitative trading.

## Usage

```
import "github.com/blacklabcapital/technical"
```

You can use `technical` as a utility library for use in techinical/quantitative analysis or trading programs.

`technical` can be and has been used in real time trading, even in low-latency microsecond environments. Benchmarks have been provided with comparisons to other statistical libraries for common functions.

The library has extensive unit tests which also double as examples for common usage. Please see the godoc for package documentation.

The main indicators (currently) explicitly implemented are:

- **Exponentially Weighted Moving Average**

- **Bollinger Bands**
- **Average True Range**

Various other indicators can be trivially composed with the included stats functions, such as a Simple Moving Average.

## Contributing

`master` holds the latest current stable version of technical. Commits with a minor version are guaranteed to have no breaking API changes, only feature additions and bug fixes.

`dev` holds the latest commits and is where active development takes place. If you submit a pull request it should be against the `dev` branch.

`<major.minor>` are version branches. Tested changes from `dev` are staged for a release by merging into the appropriate version branch.
