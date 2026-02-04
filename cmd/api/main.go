package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CaioBonalume/nfes/config"
	utils "github.com/CaioBonalume/nfes/internal/Utils"
	"github.com/CaioBonalume/nfes/internal/infrastructure/http/handlers"
)

func main() {
	// Configura timezone
	time.LoadLocation("America/Sao_Paulo")
	// Carrega as variáveis de ambiente do arquivo .env
	config.LoadEnv()

	if err := utils.InstallDependencies(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("/cnpj", handlers.ConsultaCNPJHandler)
	http.HandleFunc("/nfse/emitida/todas-notas", handlers.ConsultaNotasEnviadasHandler)
	http.HandleFunc("/nfse/emitida/nota", handlers.ConsultaNotaEnviadaEspecificaHandler)
	http.HandleFunc("/nfse/recebida/todas-notas", handlers.ConsultaNotasRecebidasHandler)
	http.HandleFunc("/nfse/emitir", handlers.EnviarNotaHandler)
	http.HandleFunc("/nfse/cancelar", handlers.CancelarNotaHandler)
	http.HandleFunc("/danfe/consultar", handlers.ConsultaDanfeHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
