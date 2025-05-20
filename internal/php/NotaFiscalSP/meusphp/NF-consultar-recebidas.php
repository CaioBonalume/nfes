<?php

require __DIR__ . '/../vendor/autoload.php';

use NotaFiscalSP\Constants\Params;
use NotaFiscalSP\Entities\Requests\NF\Period;
use NotaFiscalSP\NotaFiscalSP;

// Recebe os parâmetros do terminal
$options = getopt("", [
    "cnpj:",
    "cert_path:",
    "cert_pass:",
    "inicio:",
    "fim:",
    "im::" // <- Opcional
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

$nf = new NotaFiscalSP($config);

$period = new Period();
$period->setDtInicio($options['inicio'])
       ->setDtFim($options['fim'])
       ->setPagina(1); // Paginação padrão é 1 que traz os 50 primeiros registros

// Define IM se informado
if (!empty($options['im'])) {
    $period->setInscricaoMunicipal($options['im']);
}

$emitidas = $nf->notasRecebidas($period);

echo json_encode($emitidas->getResponse());
exit;
