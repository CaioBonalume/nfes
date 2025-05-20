package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const defaultTimeout = 15 * time.Second
const defaultUserAgent = "GoNFe/0.1"

func HttpClient(certPath, certKey string) (*http.Client, error) {
	tlsConfig := tls.Config{}

	cert, err := tls.LoadX509KeyPair(certPath, certKey)
	if err != nil {
		return nil, fmt.Errorf("erro no carregamento do certificado digital. Detalhes: %w", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("erro no carregamento da cadeia de certificados do sistema. Detalhes: %w", err)
	}

	tlsConfig.RootCAs = caCertPool
	tlsConfig.Renegotiation = tls.RenegotiateOnceAsClient

	client := http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
		},
	}

	return &client, nil
}

func newRequest(url string, soapAction string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if soapAction != "" {
		req.Header.Set("SOAPAction", soapAction)
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("User-Agent", defaultUserAgent)

	return req, nil
}

func SendRequest(obj interface{}, url, xmlns, soapAction string, client *http.Client, optReq ...func(req *http.Request)) ([]byte, error) {
	xmlfile, err := xml.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("erro na geração do XML de requisição. Detalhes: %w", err)
	}

	xmlfile, err = getSoapEnvelope(xmlfile, xmlns)
	if err != nil {
		return nil, fmt.Errorf("erro na geração do envelope SOAP. Detalhes: %w", err)
	}
	xmlfile = []byte(append([]byte(xml.Header), xmlfile...))

	req, err := newRequest(url, soapAction, xmlfile)
	if err != nil {
		return nil, fmt.Errorf("erro na criação da requisição (http.Request) para a URL %s. Detalhes: %w", url, err)
	}
	for _, opt := range optReq {
		opt(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição ao WebService %s. Detalhes: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, &WSError{url, resp.StatusCode, resp.Status, string(body)}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro na leitura do corpo da resposta: %w", err)
	}

	xmlfile, err = readSoapEnvelope(body)
	if err != nil {
		return nil, err
	}

	return xmlfile, err
}

func SendSOAPRequest(certPath, certKey, url, soapAction, signedXML string) ([]byte, error) {
	client, err := HttpClient(certPath, certKey)
	if err != nil {
		return nil, err
	}

	// Monta o envelope SOAP manualmente
	envelope := fmt.Sprintf(`
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    %s
  </soap:Body>
</soap:Envelope>`, signedXML)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(envelope)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	if soapAction != "" {
		req.Header.Set("SOAPAction", soapAction)
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
