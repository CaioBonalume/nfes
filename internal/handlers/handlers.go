package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	CertPath string
	CertPass string
	CertPem  string
	CertKey  string
)

func validaCertificado() error {
	if _, err := os.Stat(CertPath); err != nil {
		return fmt.Errorf("certificado não encontrado: %v", err)
	}
	if CertPass == "" {
		return fmt.Errorf("senha do certificado não informada")
	}
	return nil
}

func executaPHP(w http.ResponseWriter, args []string) ([]byte, bool) {
	cmd := exec.Command("php", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Erro ao executar PHP: %s", err)
		log.Printf("Stderr: %s", stderr.String())
		http.Error(w, "Erro ao executar script PHP: "+err.Error()+"\n"+stderr.String(), http.StatusInternalServerError)
		return nil, false
	}
	return output, true
}

func ConsultaCNPJHandler(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	if cnpj == "" {
		http.Error(w, "CNPJ não informado", http.StatusBadRequest)
		return
	}

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args := []string{
		"php/NotaFiscalSP/meusphp/Consulta-CNPJ.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--cert_path=%s", CertPath),
		fmt.Sprintf("--cert_pass=%s", CertPass),
	}

	cmd := exec.Command("php", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func ConsultaNotasHandler(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	im := r.URL.Query().Get("im") // opcional
	inicio := r.URL.Query().Get("inicio")
	fim := r.URL.Query().Get("fim")

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Monta os argumentos básicos
	args := []string{
		"php/NotaFiscalSP/meusphp/NF-consultar.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--cert_path=%s", CertPath),
		fmt.Sprintf("--cert_pass=%s", CertPass),
		fmt.Sprintf("--inicio=%s", inicio),
		fmt.Sprintf("--fim=%s", fim),
	}

	// Se IM for passado, adiciona
	if im != "" {
		args = append(args, fmt.Sprintf("--im=%s", im))
	}

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func ConsultaNotaEspecificaHandler(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	numNota := r.URL.Query().Get("num_nota")

	if cnpj == "" {
		http.Error(w, "CNPJ não informado", http.StatusBadRequest)
		return
	}
	if numNota == "" {
		http.Error(w, "Número da nota não informado", http.StatusBadRequest)
		return
	}

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args := []string{
		"php/NotaFiscalSP/meusphp/NF-consultar-especifica.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--num_nota=%s", numNota),
		fmt.Sprintf("--cert_path=%s", CertPath),
		fmt.Sprintf("--cert_pass=%s", CertPass),
	}

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
