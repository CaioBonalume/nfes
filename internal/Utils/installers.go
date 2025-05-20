package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func checkAndInstall(cmd, installCmd string, sudo bool) error {
	_, err := exec.LookPath(cmd)
	if err == nil {
		fmt.Printf("%s já está instalado.\n", cmd)
		return nil
	}
	fmt.Printf("%s não encontrado. Instalando...\n", cmd)
	if sudo {
		installCmd = "sudo " + installCmd
	}
	c := exec.Command("bash", "-c", installCmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func appendToPHPIni(lines []string, phpIniPath string) error {
	f, err := os.OpenFile(phpIniPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func InstallDependencies() error {
	osType := runtime.GOOS

	var phpInstall, composerInstall, opensslInstall, phpIniPath string
	var sudo bool

	switch osType {
	case "linux":
		phpInstall = "apt-get update && apt-get install -y php"
		composerInstall = "apt-get install -y composer"
		opensslInstall = "apt-get install -y openssl"
		// O caminho pode variar, mas geralmente é:
		phpIniPath = "/etc/php/8.3/apache2/php.ini"
		sudo = true
	case "darwin":
		phpInstall = "brew install php"
		composerInstall = "brew install composer"
		opensslInstall = "brew install openssl"
		// O caminho pode variar, mas geralmente é:
		phpIniPath = "/opt/homebrew/etc/php/8.4/php.ini"
		sudo = false
	default:
		fmt.Println("Sistema operacional não suportado para instalação automática.")
		os.Exit(1)
	}

	if err := checkAndInstall("php", phpInstall, sudo); err != nil {
		fmt.Printf("Erro ao instalar PHP: %v\n", err)
		os.Exit(1)
	}
	if err := checkAndInstall("composer", composerInstall, sudo); err != nil {
		fmt.Printf("Erro ao instalar Composer: %v\n", err)
		os.Exit(1)
	}
	if err := checkAndInstall("openssl", opensslInstall, sudo); err != nil {
		fmt.Printf("Erro ao instalar OpenSSL: %v\n", err)
		os.Exit(1)
	}

	// Adiciona configurações ao php.ini
	fmt.Printf("Adicionando configurações ao php.ini em %s...\n", phpIniPath)
	phpIniLines := []string{
		"session.name = CUSTOMSESSID",
		"session.use_only_cookies = 1",
		"session.cookie_httponly = true",
		"session.use_trans_sid = 0",
	}
	if err := appendToPHPIni(phpIniLines, phpIniPath); err != nil {
		fmt.Printf("Erro ao adicionar configurações ao php.ini: %v\n", err)
		// Não faz exit, apenas alerta
	}

	// Executa composer install no diretório do projeto PHP
	fmt.Println("Executando 'composer install' no diretório do projeto PHP...")
	cmd := exec.Command("composer", "install")
	cmd.Dir = "php/NotaFiscalSP"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Erro ao executar composer install: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ambiente PHP pronto!")

	return nil
}
