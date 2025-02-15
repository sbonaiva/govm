package domain

var (
	ShellRunCommandsFiles = map[string]string{
		"/bin/bash":     ".bashrc",
		"/usr/bin/bash": ".bashrc",
		"/bin/zsh":      ".zshrc",
		"/usr/bin/zsh":  ".zshrc",
		"/bin/ksh":      ".kshrc",
		"/usr/bin/ksh":  ".kshrc",
		"/bin/fish":     ".config/fish/config.fish",
		"/usr/bin/fish": ".config/fish/config.fish",
	}
)
