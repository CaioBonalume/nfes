<?php

require __DIR__ . '/../../vendor/autoload.php';

use NotaFiscalSP\Constants\FieldData\RPSType;
use NotaFiscalSP\Entities\Requests\NF\Rps;
use NotaFiscalSP\NotaFiscalSP;
use NotaFiscalSP\Constants\Params;

// Recebe os parâmetros do terminal
$options = getopt("", [
    "cnpj:",
    "cert_path:",
    "cert_pass:",
    "im:", // <- Opcional
    "numero-rps:",
    "tipo-rps::", // <- Opcional
    "valor-servicos:",
    "codigo-servico:",
    "aliquota-servicos:",
    "cnpj-tomador:",
    "razao-social-tomador:",
    "tipo-logradouro:",
    "logradouro:",
    "numero-endereco:",
    "bairro:",
    "cep:",
    "email-tomador:",
    "discriminacao:"
]);

$config = [
    Params::CNPJ => $options['cnpj'],
    Params::CERTIFICATE_PATH => $options['cert_path'],
    Params::CERTIFICATE_PASS => $options['cert_pass']
];
// Adiciona IM apenas se for informado
if (!empty($options['im'])) {
    $config[Params::IM] = $options['im'];
}

// Instancie a Classe
$nf = new NotaFiscalSP($config);

// Monte a RPS
$rps = new Rps();
date_default_timezone_set('America/Sao_Paulo');
$rps->setDataEmissao(date('Y-m-d'));
$rps->setNumeroRps($options['numero-rps']);
// $rps->setTipoRps(RPSType::RECIBO_PROVENIENTE_DE_NOTA_CONJUGADA); // Verificar depois
$rps->setValorServicos(isset($options['valor-servicos']) ? (float)$options['valor-servicos'] : 0.0);
$rps->setCodigoServico(isset($options['codigo-servico']) ? (int)$options['codigo-servico'] : 0);
$rps->setAliquotaServicos(isset($options['aliquota-servicos']) ? (float)$options['aliquota-servicos'] : 0.0);
$rps->setCnpj($options['cnpj-tomador']);
$rps->setRazaoSocialTomador($options['razao-social-tomador']);
$rps->setTipoLogradouro($options['tipo-logradouro']);
$rps->setLogradouro($options['logradouro']);
$rps->setNumeroEndereco(isset($options['numero-endereco']) ? (int)$options['numero-endereco'] : 0);
$rps->setBairro($options['bairro']);
$rps->setCidade('3550308'); // São Paulo
$rps->setUf('SP'); // São Paulo
$rps->setCep($options['cep']);
$rps->setEmailTomador($options['email-tomador']);
$rps->setDiscriminacao($options['discriminacao']);

// Envie a Requisição
$request = $nf->enviarNota($rps);

// Utilize algum dos métodos do response para verificar o resultado
echo $request->getXmlOutput();
exit;