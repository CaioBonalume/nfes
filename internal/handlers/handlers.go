package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type NotaRequest struct {
	CNPJ               string `json:"cnpj"`
	IM                 string `json:"im"`
	NumeroRPS          string `json:"numero_rps"`
	TipoRPS            string `json:"tipo_rps"`
	ValorServicos      string `json:"valor_servicos"`
	CodigoServico      string `json:"codigo_servico"`
	AliquotaServicos   string `json:"aliquota_servicos"`
	CNPJTomador        string `json:"cnpj_tomador"`
	RazaoSocialTomador string `json:"razao_social_tomador"`
	TipoLogradouro     string `json:"tipo_logradouro"`
	Logradouro         string `json:"logradouro"`
	NumeroEndereco     string `json:"numero_endereco"`
	Bairro             string `json:"bairro"`
	CEP                string `json:"cep"`
	EmailTomador       string `json:"email_tomador"`
	Discriminacao      string `json:"discriminacao"`
}

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

func ConsultaNotasEnviadasHandler(w http.ResponseWriter, r *http.Request) {
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
		"php/NotaFiscalSP/meusphp/NF-consultar-enviadas.php",
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

func ConsultaNotaEnviadaEspecificaHandler(w http.ResponseWriter, r *http.Request) {
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
		"php/NotaFiscalSP/meusphp/NF-consultar-emitida-especifica.php",
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

func ConsultaNotasRecebidasHandler(w http.ResponseWriter, r *http.Request) {
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
		"php/NotaFiscalSP/meusphp/NF-consultar-recebidas.php",
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

func EnviarNotaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req NotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.CNPJ == "" || CertPath == "" || CertPass == "" {
		http.Error(w, "Parâmetros de cabeçalho obrigatórios ausentes", http.StatusBadRequest)
		return
	}

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args := []string{
		"php/NotaFiscalSP/meusphp/NF-emitir.php",
		fmt.Sprintf("--cnpj=%s", req.CNPJ),
		fmt.Sprintf("--im=%s", req.IM),
		fmt.Sprintf("--cert_path=%s", CertPath),
		fmt.Sprintf("--cert_pass=%s", CertPass),
		fmt.Sprintf("--valor-servicos=%s", req.ValorServicos),
		fmt.Sprintf("--codigo-servico=%s", req.CodigoServico),
		fmt.Sprintf("--aliquota-servicos=%s", req.AliquotaServicos),
		fmt.Sprintf("--cnpj-tomador=%s", req.CNPJTomador),
		fmt.Sprintf("--razao-social-tomador=%s", req.RazaoSocialTomador),
		fmt.Sprintf("--tipo-logradouro=%s", req.TipoLogradouro),
		fmt.Sprintf("--logradouro=%s", req.Logradouro),
		fmt.Sprintf("--numero-endereco=%s", req.NumeroEndereco),
		fmt.Sprintf("--bairro=%s", req.Bairro),
		fmt.Sprintf("--cep=%s", req.CEP),
		fmt.Sprintf("--email-tomador=%s", req.EmailTomador),
		fmt.Sprintf("--discriminacao=%s", req.Discriminacao),
	}

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func CancelarNotaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
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
		"php/NotaFiscalSP/meusphp/NF-cancelar.php",
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
