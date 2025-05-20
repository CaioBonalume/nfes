<?php


require __DIR__ . '/../vendor/autoload.php';

use NotaFiscalSP\Constants\Params;
use NotaFiscalSP\NotaFiscalSP;

// Recebe os parâmetros do terminal
$options = getopt("", [
    "cnpj:",
    "cert_path:",
    "cert_pass:",
]);

if (empty($options['cnpj']) || empty($options['cert_path']) || empty($options['cert_pass'])) {
    echo json_encode(['error' => 'Parâmetros obrigatórios: --cnpj, --cert_path, --cert_pass']);
    exit(1);
}

$config = [
    Params::CNPJ => $options['cnpj'],
    Params::CERTIFICATE_PATH => $options['cert_path'],
    Params::CERTIFICATE_PASS => $options['cert_pass']
];

$nfSP = new NotaFiscalSP($config);

// Consulta o CNPJ
$response = $nfSP->cnpjInfo($options['cnpj']);

echo json_encode($response);
exit;