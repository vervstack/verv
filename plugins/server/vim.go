package server

import (
	"context"
	"os"
	"path/filepath"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const (
	vimPackage        = "vim"
	vimrcLocalPath    = "/etc/vim/vimrc.local"
	vimrcLocalMode    = 0o644
	vimBasicPath      = "/usr/bin/vim.basic"
	editorAlternative = "editor"
	editorLinkPath    = "/etc/alternatives/editor"

	// Ubuntu's /etc/vim/vimrc ships with `syntax on` commented out and sources
	// /etc/vim/vimrc.local when it is readable.
	vimrcLocalContent = "\" managed by verv setup-server\nsyntax on\n\nautocmd FileType yaml setlocal tabstop=2 shiftwidth=2 softtabstop=2 expandtab\n"
)

func isVimSetUp(_ context.Context, _ Options) bool {
	req := cmd.Request{Tool: "dpkg-query", Args: []string{"-W", "-f", "${Status}", vimPackage}}

	status, err := cmd.Execute(req)
	if err != nil || !isDpkgInstalled(status) {
		return false
	}

	content, err := os.ReadFile(vimrcLocalPath)
	if err != nil || string(content) != vimrcLocalContent {
		return false
	}

	editorPath, err := filepath.EvalSymlinks(editorLinkPath)
	if err != nil {
		return false
	}

	return editorPath == vimBasicPath
}

func installVim(_ context.Context, _ Options) (string, error) {
	err := aptGet("update")
	if err != nil {
		return "", rerrors.Wrap(err, "error updating package index")
	}

	err = aptGet("install", "-y", vimPackage)
	if err != nil {
		return "", rerrors.Wrap(err, "error installing vim")
	}

	content := []byte(vimrcLocalContent)

	err = os.WriteFile(vimrcLocalPath, content, vimrcLocalMode)
	if err != nil {
		return "", rerrors.Wrap(err, "error writing system vimrc")
	}

	err = run("update-alternatives", "--set", editorAlternative, vimBasicPath)
	if err != nil {
		return "", rerrors.Wrap(err, "error setting vim as the default editor")
	}

	return "vim installed and set as the default editor", nil
}
