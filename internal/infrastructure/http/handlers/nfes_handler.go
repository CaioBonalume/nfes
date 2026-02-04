package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CaioBonalume/nfes/config"
)

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
		"../../internal/php/meusphp/nfse/NF-consultar-enviadas.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
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
		"../../internal/php/NotaFiscalSP/meusphp/NF-consultar-emitida-especifica.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--num_nota=%s", numNota),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
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
		"../../internal/php/NotaFiscalSP/meusphp/NF-consultar-recebidas.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
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
	if req.CNPJ == "" || config.Env.CERT_PATH == "" || config.Env.CERT_PASS == "" {
		http.Error(w, "Parâmetros de cabeçalho obrigatórios ausentes", http.StatusBadRequest)
		return
	}

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar ultimo numero de RPS gerado
	// var ultimoRPS string

	args := []string{
		"../../internal/php/meusphp/nfse/NF-emitir.php",
		fmt.Sprintf("--cnpj=%s", req.CNPJ),
		fmt.Sprintf("--im=%s", req.IM),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
		fmt.Sprintf("--numero-rps=%s", req.NumeroRPS), // <-- Aqui
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
		"../../internal/php/NotaFiscalSP/meusphp/NF-cancelar.php",
		fmt.Sprintf("--cnpj=%s", cnpj),
		fmt.Sprintf("--num_nota=%s", numNota),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
	}

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
