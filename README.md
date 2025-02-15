# Go Version Manager (govm)

![Go Version Manager](https://img.shields.io/badge/version-0.0.2-blue.svg)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)

**Go Version Manager (govm)** is a command-line tool written in Go that simplifies the management of Go installations on your machine. 
With a few simple commands, you can easily list, install, and uninstall different versions of the Go programming language, making it perfect for developers who need to switch between versions for various projects.

## Requirements

- Linux or macOS
- Internet connection
- [curl](https://curl.se/)
- [tar](https://www.gnu.org/software/tar/)
- [sed](https://www.gnu.org/software/sed/)

## Installation

To get started with `govm`, you need to run the following in a terminal:

```bash
curl -s "https://raw.githubusercontent.com/sbonaiva/govm/refs/heads/main/scripts/install.sh" | bash
```

Alternatively, you can build it from source:

```bash
git clone https://github.com/sbonaiva/govm.git
cd govm
make install
```

## Usage

### List

```bash
govm list
```

This command will display all Go versions currently installed on your system, as well as the versions available for installation.

### Install

```bash
govm install [version]
```

Replace `[version]` with the desired Go version (e.g., `go1.23.6`). This command will download and install the specified version.

### Uninstall

```bash
govm uninstall
```

This command removes the currently installed version of Go from your system.

## Troubleshooting
If you encounter any issues while using the application, please follow these steps:

**Check the Log File**: 

The first step in troubleshooting is to review the log file for any errors or warnings. You can do this by running the following command in your terminal:

```bash
cat ~/.govm/govm.log
```

This will display the contents of the log file, which may provide insight into what went wrong.

**Identify the Issue**: 

Look for any error messages or unusual behavior in the log. Take note of any specific error codes or messages that may help in diagnosing the problem.

**Open an Issue**: 

If you are unable to resolve the issue after reviewing the log file, please open a new issue on our GitHub repository. Make sure to include the following information:
   - A description of the problem you are experiencing.
   - Steps to reproduce the issue.
   - Any relevant error messages or log snippets.

This will help us assist you more effectively.

## Contributing

Contributions are welcome! If you would like to contribute to `govm`, please follow these steps:

1. Fork the repository.
2. Create a new branch (\`git checkout -b feature/YourFeature\`).
3. Make your changes and commit them using [Conventional Commits](https://www.conventionalcommits.org/).
4. Push to the branch (\`git push origin feature/YourFeature\`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Go](https://golang.org) - For providing a powerful and efficient language.
- [Cobra](https://github.com/spf13/cobra) - For simplifying the creation of command-line applications in Go.
- [GitHub](https://github.com) - For hosting this project and fostering open-source collaboration.

## Contact

For any questions or suggestions, feel free to open an issue.

---

Thank you for using Go Version Manager! 

Proudly made in Brazil 🇧🇷.
