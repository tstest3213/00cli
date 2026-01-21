package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHDeployer implementa deploy via SSH
type SSHDeployer struct {
	Host     string
	Port     int
	User     string
	SSHKey   string
	Password string
}

// Execute executa comandos via SSH
func (d *SSHDeployer) Execute(commands []string) error {
	config := &ssh.ClientConfig{
		User:            d.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Em produção, use validação adequada
		Timeout:         10 * time.Second,
	}

	// Autenticação por chave SSH
	if d.SSHKey != "" {
		key, err := os.ReadFile(d.SSHKey)
		if err != nil {
			return fmt.Errorf("erro ao ler chave SSH: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("erro ao parsear chave SSH: %w", err)
		}

		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else if d.Password != "" {
		// Autenticação por senha
		config.Auth = []ssh.AuthMethod{
			ssh.Password(d.Password),
		}
	} else {
		return fmt.Errorf("nenhuma forma de autenticação configurada (ssh_key ou password)")
	}

	// Conectar ao servidor
	addr := fmt.Sprintf("%s:%d", d.Host, d.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("erro ao conectar via SSH: %w", err)
	}
	defer client.Close()

	// Executar comandos
	for i, cmd := range commands {
		fmt.Printf("  [%d/%d] Executando: %s\n", i+1, len(commands), cmd)

		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("erro ao criar sessão SSH: %w", err)
		}

		// Capturar stdout e stderr
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if err := session.Run(cmd); err != nil {
			session.Close()
			return fmt.Errorf("erro ao executar comando '%s': %w", cmd, err)
		}

		session.Close()
	}

	return nil
}

// UploadFile faz upload de um arquivo via SCP
func (d *SSHDeployer) UploadFile(localPath, remotePath string) error {
	config := &ssh.ClientConfig{
		User:            d.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	if d.SSHKey != "" {
		key, err := os.ReadFile(d.SSHKey)
		if err != nil {
			return fmt.Errorf("erro ao ler chave SSH: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("erro ao parsear chave SSH: %w", err)
		}

		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else if d.Password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(d.Password),
		}
	}

	addr := fmt.Sprintf("%s:%d", d.Host, d.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("erro ao conectar via SSH: %w", err)
	}
	defer client.Close()

	// Abrir arquivo local
	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo local: %w", err)
	}
	defer srcFile.Close()

	// Criar sessão SCP
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("erro ao criar sessão: %w", err)
	}
	defer session.Close()

	// Executar SCP
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%04o %d %s\n", 0644, getFileSize(srcFile), filepath.Base(remotePath))
		io.Copy(w, srcFile)
		fmt.Fprint(w, "\x00")
	}()

	cmd := fmt.Sprintf("scp -t %s", remotePath)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("erro ao fazer upload: %w", err)
	}

	return nil
}

func getFileSize(file *os.File) int64 {
	info, err := file.Stat()
	if err != nil {
		return 0
	}
	return info.Size()
}
