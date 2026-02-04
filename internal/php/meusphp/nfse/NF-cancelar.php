<?php


require __DIR__ . '/../../vendor/autoload.php';

use NotaFiscalSP\Constants\Params;
use NotaFiscalSP\NotaFiscalSP;

// Recebe os parâmetros do terminal
$options = getopt("", [
    "cnpj:",
    "num_nota:",
    "cert_path:",
    "cert_pass:",
]);

if (empty($options['num_nota']) || empty($options['cert_path']) || empty($options['cert_pass'])) {
    echo json_encode(['error' => 'Parâmetros obrigatórios: --num_nota, --cert_path, --cert_pass']);
    exit(1);
}

$config = [
    Params::CNPJ => $options['cnpj'],
    Params::CERTIFICATE_PATH => $options['cert_path'],
    Params::CERTIFICATE_PASS => $options['cert_pass']
];

$nfSP = new NotaFiscalSP($config);

// Cancela nota fiscal
$response = $nfSP->cancelarNota($options['num_nota']);

echo json_encode($response);
exit;