package utils

import "fmt"

func XMLConsultaCNPJ(cnpj string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<p1:PedidoConsultaCNPJ xmlns:p1="http://www.prefeitura.sp.gov.br/nfe" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <Cabecalho Versao="1">
        <CPFCNPJRemetente>
            <CNPJ>%s</CNPJ>
        </CPFCNPJRemetente>
    </Cabecalho>
    <CNPJContribuinte>
        <CNPJ>%s</CNPJ>
    </CNPJContribuinte>
</p1:PedidoConsultaCNPJ>`, cnpj, cnpj)
}

func XMLConsultaNotaEspecifica(cnpj, im, numNota string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<p1:PedidoConsultaNFe xmlns:p1="http://www.prefeitura.sp.gov.br/nfe" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <Cabecalho Versao="1">
        <CPFCNPJRemetente>
            <CNPJ>%s</CNPJ>
        </CPFCNPJRemetente>
    </Cabecalho>
    <Detalhe>
        <ChaveNFe>
            <InscricaoPrestador>%s</InscricaoPrestador>
            <NumeroNFe>%s</NumeroNFe>
        </ChaveNFe>
    </Detalhe>
</p1:PedidoConsultaNFe>`, cnpj, im, numNota)
}
