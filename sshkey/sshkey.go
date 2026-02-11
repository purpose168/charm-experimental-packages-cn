// Package sshkey 提供用于解析带密码短语支持的 SSH 私钥的工具。
package sshkey

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"golang.org/x/crypto/ssh"
)

// Open 读取路径并解析密钥。
func Open(keyPath string) (ssh.Signer, error) {
	pemBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("sshkey: %w", err)
	}
	return Parse(keyPath, pemBytes)
}

// Parse 尝试将给定的 PEM 解析为 ssh.Signer。
// 如果密钥已加密，它会要求输入密码短语。
// 'identifier' 用于在要求密码短语时向用户标识密钥。
func Parse(identifier string, pemBytes []byte) (ssh.Signer, error) {
	return doParse(identifier, pemBytes, ssh.ParsePrivateKey, ssh.ParsePrivateKeyWithPassphrase)
}

// ParseRaw 尝试将给定的 PEM 解析为私钥。
// 如果密钥已加密，它会要求输入密码短语。
// 'identifier' 用于在要求密码短语时向用户标识密钥。
func ParseRaw(identifier string, pemBytes []byte) (any, error) {
	return doParse(identifier, pemBytes, ssh.ParseRawPrivateKey, ssh.ParseRawPrivateKeyWithPassphrase)
}

func doParse[T any](
	identifier string,
	pemBytes []byte,
	parse func(pemBytes []byte) (T, error),
	parseWithPass func(pemBytes, passphrase []byte) (T, error),
) (T, error) {
	result, err := parse(pemBytes)
	if isPassphraseMissing(err) {
		passphrase, err := ask(identifier)
		if err != nil {
			return result, fmt.Errorf("sshkey: %w", err)
		}
		result, err := parseWithPass(pemBytes, passphrase)
		if err != nil {
			return result, fmt.Errorf("sshkey: %w", err)
		}
		return result, nil
	}
	if err != nil {
		return result, fmt.Errorf("sshkey: %w", err)
	}
	return result, nil
}

func isPassphraseMissing(err error) bool {
	var kerr *ssh.PassphraseMissingError
	return errors.As(err, &kerr)
}

func ask(path string) ([]byte, error) {
	var pass string
	if err := huh.Run(
		huh.NewInput().
			Inline(true).
			Value(&pass).
			Title(fmt.Sprintf("Enter the passphrase to unlock %q: ", path)).
			EchoMode(huh.EchoModePassword),
	); err != nil {
		return nil, fmt.Errorf("sshkey: %w", err)
	}
	return []byte(pass), nil
}
