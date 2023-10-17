package main

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

func gitClone(repoSSHAddress, keyPath, branch, workDir string) error {
	sshKey, key_err := os.ReadFile(keyPath)
	if key_err != nil {
		logger.Fatal().
			Err(key_err).
			Str("module", "git").
			Msg("")
		//os.Exit(0)
	}

	signer, sign_err := ssh.ParsePrivateKey([]byte(sshKey))
	if sign_err != nil {
		logger.Fatal().
			Err(sign_err).
			Str("module", "git").
			Msg("")
		//os.Exit(0)
	}

	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}

	cloneOpts := &git.CloneOptions{
		Auth:          auth,
		URL:           repoSSHAddress,
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
	}

	_, cl_err := git.PlainClone(workDir, false, cloneOpts)
	if cl_err != nil {
		logger.Fatal().
			Err(cl_err).
			Str("module", "git").
			Msg("")
		//os.Exit(0)
	}

	return nil
}
