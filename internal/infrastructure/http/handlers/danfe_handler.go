package handlers

import (
	"fmt"
	"net/http"

	"github.com/CaioBonalume/nfes/config"
)

func ConsultaDanfeHandler(w http.ResponseWriter, r *http.Request) {
	chave := r.URL.Query().Get("chave")
	if chave == "" {
		http.Error(w, "Chave da NF-e não informada", http.StatusBadRequest)
		return
	}

	if err := validaCertificado(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args := []string{
		"../../internal/php/meusphp/danfe/NFe-consultar.php",
		fmt.Sprintf("--chave=%s", chave),
		fmt.Sprintf("--cert_path=%s", config.Env.CERT_PATH),
		fmt.Sprintf("--cert_pass=%s", config.Env.CERT_PASS),
	}

	output, ok := executaPHP(w, args)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(output)
}
