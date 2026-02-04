package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/CaioBonalume/nfes/config"
)

type NotaRequest struct {
	CNPJ               string `json:"cnpj"`
	IM                 string `json:"im"`
	NumeroRPS          string `json:"numero-rps"`
	TipoRPS            string `json:"tipo-rps"`
	ValorServicos      string `json:"valor-servicos"`
	CodigoServico      string `json:"codigo-servico"`
	AliquotaServicos   string `json:"aliquota-servicos"`
	CNPJTomador        string `json:"cnpj-tomador"`
	RazaoSocialTomador string `json:"razao-social-tomador"`
	TipoLogradouro     string `json:"tipo-logradouro"`
	Logradouro         string `json:"logradouro"`
	NumeroEndereco     string `json:"numero-endereco"`
	Bairro             string `json:"bairro"`
	CEP                string `json:"cep"`
	EmailTomador       string `json:"email-tomador"`
	Discriminacao      string `json:"discriminacao"`
}

func validaCertificado() error {
	if _, err := os.Stat(config.Env.CERT_PATH); err != nil {
		return fmt.Errorf("certificado não encontrado: %v", err)
	}
	if config.Env.CERT_PASS == "" {
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
		"../../internal/php/meusphp/nfse/Consulta-CNPJ.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
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
