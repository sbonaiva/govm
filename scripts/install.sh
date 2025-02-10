#!/bin/bash

set -e

track_last_command() {
    last_command=$current_command
    current_command=$BASH_COMMAND
}
trap track_last_command DEBUG

echo_failed_command() {
    local exit_code="$?"
	if [[ "$exit_code" != "0" ]]; then
		echo "'$last_command': command failed with exit code $exit_code."
	fi
}
trap echo_failed_command EXIT

# Global variables
export GOVM_HOST="https://github.com/sbonaiva/govm/releases/download"
export GOVM_VERSION="0.0.1"
export GOVM_DIR="$HOME/.govm/bin"
export GOVM_TMP_DIR="${TMPDIR:-/tmp}/govm"

# Local variables
govm_path_snippet=$( cat << 'EOF'
# The next line updates PATH for the govm binary.
export PATH=$PATH:$HOME/.govm/bin
EOF
)

echo "
 ▗▄▄▖ ▗▄▖ ▗▖  ▗▖▗▖  ▗▖
▐▌   ▐▌ ▐▌▐▌  ▐▌▐▛▚▞▜▌
▐▌▝▜▌▐▌ ▐▌▐▌  ▐▌▐▌  ▐▌
▝▚▄▞▘▝▚▄▞▘ ▝▚▞▘ ▐▌  ▐▌                  
"

# Sanity checks
echo "Looking for tar..."
if ! command -v tar > /dev/null; then
	echo "Not found."
	echo "======================================================================================================"
	echo " Please install tar on your system using your favourite package manager."
	echo ""
	echo " Restart after installing tar."
	echo "======================================================================================================"
	echo ""
	exit 1
fi

echo "Looking for curl..."
if ! command -v curl > /dev/null; then
	echo "Not found."
	echo ""
	echo "======================================================================================================"
	echo " Please install curl on your system using your favourite package manager."
	echo ""
	echo " Restart after installing curl."
	echo "======================================================================================================"
	echo ""
	exit 1
fi

echo "Looking for sed..."
if [ -z $(command -v sed) ]; then
	echo "Not found."
	echo ""
	echo "======================================================================================================"
	echo " Please install sed on your system using your favourite package manager."
	echo ""
	echo " Restart after installing sed."
	echo "======================================================================================================"
	echo ""
	exit 1
fi

# Create directories
echo "Looking for a previous installation of govm..."

if [ -d "$GOVM_DIR" ]; then
	echo "Removing existing govm installation..."
	rm -rf "$GOVM_DIR"
fi
mkdir -p "$GOVM_DIR"

if [ -d "$GOVM_TMP_DIR" ]; then
	echo "Removing existing govm tmp..."
	rm -rf "$GOVM_TMP_DIR"
fi
mkdir -p "$GOVM_TMP_DIR"

echo "Checking platform..."

# infer platform
function infer_platform() {
	local kernel
	local machine

	kernel="$(uname -s)"
	machine="$(uname -m)"

	case $kernel in
	Linux)
	  case $machine in
	  i686)
		echo "linux_386"
		;;
	  x86_64)
		echo "linux_amd64"
		;;
	  aarch64)
		echo "linux_arm64"
		;;
	  *)
	  	echo "others"
	  	;;
	  esac
	  ;;
	Darwin)
	  case $machine in
	  x86_64)
		echo "darwin_amd64"
		;;
	  arm64)
		echo "darwin_arm64"
		;;
	  *)
	  	echo "darwin_amd64"
	  	;;
	  esac
	  ;;
	*)
	  echo "others"
	esac
}

export GOVM_PLATFORM="$(infer_platform)"

if [ "$GOVM_PLATFORM" == "others" ]; then
	echo "Unsupported platform: $(uname -s) $(uname -m)"
	exit 1
fi

echo "Detected platform: $(uname -s) $(uname -m)"
govm_tar_file="govm_${GOVM_VERSION}_${GOVM_PLATFORM}.tar.gz"
govm_tmp_file="$GOVM_TMP_DIR/$govm_tar_file"

echo "Downloading govm installation file..."
curl --fail --location --progress-bar "${GOVM_HOST}/v${GOVM_VERSION}/$govm_tar_file" > "$govm_tmp_file"

echo "Extracting govm installation file..."
tar -xzf "$govm_tmp_file" -C "$GOVM_DIR"

echo "Cleaning up..."
rm -rf "$GOVM_TMP_DIR"

echo "Attempting to add govm to your PATH..."

function infer_shell() {
	local shell
	shell="$SHELL"

	case $shell in
	*bash)
		echo "$HOME/.bashrc"
		;;
	*zsh)
		echo "$HOME/.zshrc"
		;;
	*fish)
		echo "$HOME/.config/fish/config.fish"
		;;
	*ksh)
		echo "$HOME/.kshrc"
		;;
	*)
		echo "others"
		;;
	esac
}

export GOVM_SHELL="$(infer_shell)"

if [ "$GOVM_SHELL" == "others" ]; then
	shells=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.config/fish/config.fish" "$HOME/.kshrc")

	sucedded=false
	for shell in "${shells[@]}"; do
		if test -w "$shell"; then
			if ! grep -qF "$govm_path_snippet" "$shell"; then
				echo "$govm_path_snippet" >> "$shell"
				sucedded=true
				export GOVM_SHELL="$shell"
			else
				echo "govm is already in your PATH."
				sucedded=true
				export GOVM_SHELL="$shell"
			fi
		fi
	done

	if [ "$sucedded" == false ]; then
		echo "Failed to add govm to your PATH. You will need to add the following line to your shell profile manually:"
		echo ""
		echo "$govm_path_snippet"
		exit 1
	fi
else
	if test -w "$GOVM_SHELL"; then
	    if ! grep -qF "$govm_path_snippet" "$GOVM_SHELL"; then
			echo "$govm_path_snippet" >> "$GOVM_SHELL"
		else
			echo "govm is already in your PATH."
		fi
	else
		echo "Failed to add govm to your PATH. You will need to add the following line to your shell profile manually:"
		echo ""
		echo "$govm_path_snippet"
		exit 1
	fi
fi

echo ""
echo ""
echo -e "\033[0;32mInstallation completed successfully!\033[0m"
echo ""
echo ""
echo "Please open a new terminal, or run the following in the existing one:"
echo ""
echo "    source $GOVM_SHELL"
echo ""
echo "Run the following command to start using govm:"
echo ""
echo "    govm --help"
echo ""
echo "Thanks for using govm!"
