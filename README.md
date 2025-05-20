
# Aviso

Este projeto é fornecido **sem garantias de segurança**. O uso é de sua total responsabilidade. Não nos responsabilizamos por eventuais vazamentos de dados, falhas de segurança ou quaisquer danos decorrentes do uso deste software.

**Recomendamos fortemente** que você revise, adapte e implemente medidas de segurança adequadas antes de utilizar este projeto em ambientes de produção.

# NFEs SP
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
![Downloads](https://img.shields.io/github/downloads/CaioBonalume/nfes/total)
![Go.mod](https://img.shields.io/github/go-mod/go-version/CaioBonalume/nfes)



**Apenas Notas Fiscais de Serviço**

**Linguagens e pacotes necessárias para executação da aplicação**
- GoLang
- PHP
- composer
- Soap
- openssl

O projeto se trata de um servidor local que cria um módulo que realiza requisições em formato de API e realiza conversão de responses do formato XML para JSON, com motor do projeto [NotaFiscalSP](https://github.com/kaio-souza/NotaFiscalSP) criado por [kaio-souza](https://github.com/kaio-souza). 

Que integra com o sistema de notas da Prefeitura de São Paulo (Nota do Milhão), possibilitando a automatização de serviços como emissão e consulta de Notas e outros serviços relacionados.

O Grande problema não apenas no Go mas em outras linguagens é a dificuldade em assinar o XML sem utilizar métodos em JAVA ou Python ou PHP.





## Features

- Dados do próprio CNPJ
#### Consultas
- Notas fiscais de serviço emitidas/recebidas
- Nota fiscal de serviço específica emitida/recebida

#### Emissão
-  Nota Fiscal de serviço

#### Cancelamento
- Cancelar nota fiscal de serviço (Não esta funcionando)


## Environment Variables

Para que este projeto funcione, você deve colocar na mesma pasta que o main.go o arquivo .env contendo duas variavéis.

`CERT_PATH="/caminho/para/certificado.pfx"`

`KEY_PATH="senha do certificado pfx"`


## Installation

- Após realizar a etapa acima será necessário instalar as linguagens e pacotes acima citados

```bash
  cd internal/src
  go run main.go
```
**\* recomendação: Caso ocorra algum erro ao adicionar configurações ao php.ini, realizar os procedimentos de segurança de um servidor php. Você pode encontrar um link sobre como fazer isso em referências.**

Se receber o erro
```bash
Erro ao adicionar configurações ao php.ini
```
Você pode alterar a rota onde o php.ini esta instalado em internal/Utils/installers.go, tanto do linux quanto do mac. Para saber onde o arquivo esta localizado utilize o comando:
```bash
php --ini
```
Então executar novamente e verificar se o erro permanece.

## API Reference

#### Get dados do CNPJ

```
  GET /cnpj
```

| Parameter | Type     | Exemplo | Description                       |
| :-------- | :------- |:------- | :-------------------------------- |
| `cnpj` | `string` | 13456789000100 | **Required**. Seu CNPJ |

A última linha tem sua Inscrição Municipal caso possua 

#### Dados CNPJ - Data Structure

```json
{
    "success": "true",
    "message": null,
    "xmlInput": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<p1:PedidoConsultaCNPJ xmlns:p1=\"http://www.prefeitura.sp.gov.br/nfe\" ...>",
    "status": {
        "InscricaoMunicipal": "12345678",
        "EmiteNFe": "true"
    },
    "im": "12345678"
}
```

#### Get todas notas de serviço emitidas ou recebidas

```
  GET /nfse/emitida/todas-notas
```
```
  GET /nfse/recebida/todas-notas
```

| Parameter | Type     | Exemplo | Description                       |
| :-------- | :------- |:------- | :-------------------------------- |
| `cnpj`      | `string` | 13456789000100 | **Required**. Seu CNPJ |
| `im`      | `string` | 12345678 | **Optional**. Sua Inscrição Municipal, altamente recomendado caso tenha |
| `ini`      | `string` | 2025-04-01 | **Required**. Data Inicial |
| `fim`      | `string` | 2025-04-30 | **Required**. Data final |

O espaço entre as datas não pode ter mais de 31 dias.

#### Todas notas serviço emitidas - Data Structure
```json
{
    "Cabecalho": {
        "@attributes": {
            "Versao": "1"
        },
        "Sucesso": "true"
    },
    "NFe": [
        {
            "ChaveNFe": {
                "InscricaoPrestador": "12345678",
                "NumeroNFe": "265",
                "CodigoVerificacao": "WNWWQQZQ"
            },
            "DataEmissaoNFe": "2025-04-06T12:10:37",
            "DataFatoGeradorNFe": "2025-04-06T12:10:37",
            "CPFCNPJPrestador": {
                "CNPJ": "13456789000100"
            },
            "RazaoSocialPrestador": "EMPRESA EXEMPLO LTDA",
            "EnderecoPrestador": {
                "TipoLogradouro": "R",
                "Logradouro": "DOS CURIOSOS",
                "NumeroEndereco": "15",
                "Bairro": "VILA VILAS",
                "Cidade": "3550308",
                "UF": "SP",
                "CEP": "1234010"
            },
            "StatusNFe": "N",
            "TributacaoNFe": "T",
            "OpcaoSimples": "4",
            "ValorServicos": "1000",
            "CodigoServico": "2919",
            "AliquotaServicos": "0",
            "ValorISS": "0",
            "ValorCredito": "10.00",
            "ISSRetido": "false",
            "CPFCNPJTomador": {
                "CNPJ": "42591651000143"
            },
            "InscricaoMunicipalTomador": "93784694",
            "RazaoSocialTomador": "EXEMPLO EMPREESA RECEBEDORA LTDA ME",
            "EnderecoTomador": {
                "TipoLogradouro": "R",
                "Logradouro": "DOS CURIOSOS",
                "NumeroEndereco": "15",
                "ComplementoEndereco": "CJ 1208",
                "Bairro": "VILA VILAS",
                "Cidade": "3550308",
                "UF": "SP",
                "CEP": "1234010"
            },
            "EmailTomador": "financeiro@empresaexemplo.com",
            "Discriminacao": "QUALQUER COISA QUE EU QUISER ESCREVER",
            "FonteCargaTributaria": []
        },
        ...
    ]
}
```

#### Get dados de uma nota de serviço emitida específica

```
  GET /nfse/emitida/nota
```

| Parameter | Type     | Exemplo | Description                       |
| :-------- | :------- |:------- | :-------------------------------- |
| `cnpj` | `string` | 13456789000100 | **Required**. Seu CNPJ |
| `num_nota` | `string` | 265 | **Required**. Número da nota não precisa conter os zeros |

```json
{
    "success": "true",
    "message": null,
    "xmlInput": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<p1:PedidoConsultaNFe xmlns:p1=\"http://www.prefeitura.sp.gov.br/nfe\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\">...>\n",
    "xmlOutput": "<?xml version=\"1.0\" encoding=\"UTF-8\"?><RetornoConsulta xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://www.prefeitura.sp.gov.br/nfe\">...>",
    "response": {
        "Cabecalho": {
            "@attributes": {
                "Versao": "1"
            },
            "Sucesso": "true"
        },
        "NFe": {
            "ChaveNFe": {
                "InscricaoPrestador": "12345678",
                "NumeroNFe": "265",
                "CodigoVerificacao": "WNWWQQZQ"
            },
            "DataEmissaoNFe": "2025-04-06T12:10:37",
            "DataFatoGeradorNFe": "2025-04-06T12:10:37",
            "CPFCNPJPrestador": {
                "CNPJ": "13456789000100"
            },
            "RazaoSocialPrestador": "EMPRESA EXEMPLO LTDA",
            "EnderecoPrestador": {
                "TipoLogradouro": "R",
                "Logradouro": "DOS CURIOSOS",
                "NumeroEndereco": "15",
                "Bairro": "VILA VILAS",
                "Cidade": "3550308",
                "UF": "SP",
                "CEP": "1234010"
            },
            "StatusNFe": "N",
            "TributacaoNFe": "T",
            "OpcaoSimples": "4",
            "ValorServicos": "1000",
            "CodigoServico": "2919",
            "AliquotaServicos": "0",
            "ValorISS": "0",
            "ValorCredito": "10.00",
            "ISSRetido": "false",
            "CPFCNPJTomador": {
                "CNPJ": "42591651000143"
            },
            "InscricaoMunicipalTomador": "93784694",
            "RazaoSocialTomador": "EXEMPLO EMPREESA RECEBEDORA LTDA ME",
            "EnderecoTomador": {
                "TipoLogradouro": "R",
                "Logradouro": "DOS CURIOSOS",
                "NumeroEndereco": "15",
                "ComplementoEndereco": "CJ 1208",
                "Bairro": "VILA VILAS",
                "Cidade": "3550308",
                "UF": "SP",
                "CEP": "1234010"
            },
            "EmailTomador": "financeiro@empresaexemplo.com",
            "Discriminacao": "QUALQUER COISA QUE EU QUISER ESCREVER",
            "FonteCargaTributaria": []
        }
    }
}
```

#### Post nota de serviço

```
POST /nfse/emitir
```
#### Modelo body request
```json
{
    "cnpj": "00027000000000",
    "im": "00000002", // Opcional, altamente recomendado se você tiver
    "numero_rps": "123", // Aqui é o próximo numero de nota fiscal, se a ultima foi 122 a proxima será 123.
//   "tipo_rps": "1", // Opcional, Este campo precisa ser ativado no código conforme apresentado a frente
    "valor_servicos": "100.00",
    "codigo_servico": "02919",
    "aliquota_servicos": "0.0",
    "cnpj_tomador": "20000004000100",
    "razao_social_tomador": "RAZAO SOCIAL TOMADOR LTDA",
    "tipo_logradouro": "R",
    "logradouro": "NOME DA RUA",
    "numero_endereco": "147",
    "bairro": "VILA TESTE",
    "cep": "00000000",
    "email_tomador": "teste@teste.com.br",
    "discriminacao": "Teste Emissão de Notas pela API"
}
```
No caso do POST a requisição deve ter um body ao invés de Query Parameter.

O campo tipo_rps pode receber três tipos de valores:
```php
    //Recibo Provisório de Serviços
    const RECIBO_PROVISORIO = 'RPS';

    //Recibo Provisório de Serviços proveniente de Nota Fiscal Conjugada (Mista)
    const RECIBO_PROVENIENTE_DE_NOTA_CONJUGADA = 'RPS-M';

    //Cupom
    const CUPOM = 'RPS-C';
```
Para descomentar o campo utilize a rota:
```bash
cd internal/php/NotaFiscalSP/meusphp/NF-emitir.php
```

Há dois campos omitidos que sempre vão junto a requisição que são estado e cidade que são constantes de São Paulo - SP. Se tiver interesse de testar outra localidade altere em:
```bash
cd internal/php/NotaFiscalSP/meusphp/NF-emitir.php
```

#### Remove nota de serviço

Não esta funcionando

## Referência
- [kaio-souza](https://github.com/kaio-souza)
- [NotaFiscalSP](https://github.com/kaio-souza/NotaFiscalSP)
- [Configuração servidor PHP em Linux](https://php.com.br/instalacao-php-linux)
## Support & Feedback

Para entrar em contato sobre erros ou sugestões de melhora, ou até mesmo fazer parte do projeto me manda uma mensagem no Discod.
**gafuz**


## Funding
Se tiver vontade de fazer uma contribuição financeira para este projeto utilize esta chave PIX
```
1863751c-2f89-4ce6-b012-45b041bb18d4
```


## License

MIT License

Copyright (c) 2025 Bonalumi Tecnologia

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.