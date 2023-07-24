# Kinetic Commerce Infrastructure Tool (kci) 

This document provides detailed instructions on how to build, install and use the `kci` command line interface (CLI) program, an internal tool developed in Go. The `kci` is a critical tool for managing our kitchen collective infrastructure and this document will guide you on the most efficient ways to leverage its functionalities.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Building from Source](#building-from-source)
  - [Precompiled Binaries](#precompiled-binaries)
- [Usage](#usage)
  - [Common Commands](#common-commands)
- [Development](#development)
- [Testing](#testing)
- [Support](#support)

## Overview

`kci` is a CLI tool that simplifies the management of the Kinetic Commerce infrastructure. It allows for easy interaction with the system, providing a series of commands designed to facilitate tasks like querying installed services, system diagnostics, and routine maintenance. 

**note**: While this tools is in daily use, it is also under constant development. It should be considered _extremely_ rough around the edges and somewhat volatile.

## Prerequisites

Before you can install and use `kci`, you must have the following installed on your system:

- Go 1.16 or higher
- Git

## Installation

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/KineticCommerce/kci.git
```

2. Navigate into the `kci` directory:

```bash
cd kci
```

3. Build the program:

```bash
go build
```

4. Move the compiled binary to a directory within your `$PATH`:

```bash
sudo mv ./kci /usr/local/bin/
```

### Precompiled Binaries

Are not yet available.

## Usage

You can invoke `kci` from the command line by typing `kci` followed by the command you want to execute.

### Common Commands

Here are some common commands that you may find useful:

- `kci instance list`: list all instances 
- `kci instance list --filter bastion`: list all instances with a name includeing "bastion"
- `kci ssm list`: list all instances managed with ssm
- `kci ssm session --instanace i-123456`: launch an ssm session for instance i-123456 
- `kci help`: get help on all commands


## Development

Contributions to the `kci` project are welcome. To get started with development, please speak to the operations team.

## Testing

To run the tests for `kci`, navigate into the project's root directory and run:

```bash
go test ./...
```

## Support

For support or to report bugs, please file an issue on our [issues page](#) AND speak to someone on the operations team.

