package utils

import "fmt"

func SOAPTemplate(signedXML string) string {
	return fmt.Sprintf(`
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Header/>
  <soapenv:Body>
    %s
  </soapenv:Body>
</soapenv:Envelope>`, signedXML)
}
