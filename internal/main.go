package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	utils "github.com/CaioBonalume/nfes/internal/Utils"
	"github.com/CaioBonalume/nfes/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	if err := utils.InstallDependencies(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handlers.CertPath = os.Getenv("CERT_PATH")
	handlers.CertPass = os.Getenv("KEY_PATH")
	handlers.CertPem = os.Getenv("CERT_PEM")
	handlers.CertKey = os.Getenv("CERT_KEY")
	http.HandleFunc("/cnpj", handlers.ConsultaCNPJHandler)
	http.HandleFunc("/nfse/emitida/todas-notas", handlers.ConsultaNotasEnviadasHandler)
	http.HandleFunc("/nfse/emitida/nota", handlers.ConsultaNotaEnviadaEspecificaHandler)
	http.HandleFunc("/nfse/recebida/todas-notas", handlers.ConsultaNotasRecebidasHandler)
	http.HandleFunc("/nfse/emitir", handlers.EnviarNotaHandler)
	http.HandleFunc("/nfse/cancelar", handlers.CancelarNotaHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
