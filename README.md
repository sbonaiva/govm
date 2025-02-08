# Go Version Manager (govm)

![Go Version Manager](https://img.shields.io/badge/version-1.0.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**Go Version Manager (govm)** is a command-line tool written in Go that simplifies the management of Go installations on your machine. 
With a few simple commands, you can easily list, install, and uninstall different versions of the Go programming language, making it perfect for developers who need to switch between versions for various projects.

## Features

- **List Installed Versions**: Quickly check all the Go versions available and currently installed on your system.
- **Install Specific Version**: Install a specific version of Go with a single command.
- **Uninstall Version**: Remove any installed version of Go when it’s no longer needed.

## Installation

To get started with `govm`, you need to have [Go](https://golang.org/dl/) installed on your machine. 

After that, you can download the `govm` binary from the [releases page](https://github.com/sbonaiva/govm/releases).

Alternatively, you can build it from source:

```bash
git clone https://github.com/sbonaiva/govm.git
cd govm
make install
```

## Usage

### List Installed Versions

```bash
govm list
```

This command will display all Go versions currently installed on your system.

### Install a Specific Version

```bash
govm install [version]
```

Replace `[version]` with the desired Go version (e.g., `go1.17.6`). This command will download and install the specified version.

### Uninstall a Version

```bash
govm uninstall [version]
```

Replace `[version]` with the Go version you want to remove from your system.

## Contributing

Contributions are welcome! If you would like to contribute to `govm`, please follow these steps:

1. Fork the repository.
2. Create a new branch (\`git checkout -b feature/YourFeature\`).
3. Make your changes and commit them (\`git commit -m 'Add some feature'\`).
4. Push to the branch (\`git push origin feature/YourFeature\`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Go](https://golang.org) - For providing a powerful and efficient language.
- [GitHub](https://github.com) - For hosting this project and fostering open-source collaboration.

## Contact

For any questions or suggestions, feel free to open an issue.

---

Thank you for using Go Version Manager! Happy coding!
