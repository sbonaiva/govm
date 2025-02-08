package domain

var (
	ShellRunCommandsFiles = map[string]string{
		"/usr/bin/bash": ".bashrc",
		"/usr/bin/zsh":  ".zshrc",
		"/usr/bin/ksh":  ".kshrc",
		"/usr/bin/fish": ".config/fish/config.fish",
	}
)
