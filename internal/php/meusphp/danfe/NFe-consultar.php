<?php
require __DIR__ . '/../../vendor/autoload.php';

use NFePHP\NFe\Tools;
use NFePHP\Common\Certificate;

// Parâmetros do terminal
$options = getopt("", [
    "cnpj:",
    "cert_path:",
    "cert_pass:",
    "ult_nsu::" // Opcional: último NSU consultado
]);

if (empty($options['cnpj']) || empty($options['cert_path']) || empty($options['cert_pass'])) {
    echo json_encode(['error' => 'Parâmetros obrigatórios: --cnpj, --cert_path, --cert_pass']);
    exit(1);
}

$cert = file_get_contents($options['cert_path']);
$certificate = Certificate::readPfx($cert, $options['cert_pass']);

$config = [
    "atualizacao" => date('Y-m-d H:i:s'),
    "tpAmb" => 2, // 1=produção, 2=homologação
    "razaosocial" => "EMPRESA EXEMPLO LTDA",
    "siglaUF" => "SP",
    "cnpj" => "59909600000110", //$options['cnpj'],
    "schemes" => "PL_009_V4",
    "versao" => "4.00"
];

$tools = new Tools(json_encode($config), $certificate);

$ultNSU = isset($options['ult_nsu']) ? $options['ult_nsu'] : '000000000000000';

try {
    // Consulta documentos recebidos a partir do NSU informado
    $response = $tools->sefazDistDFe($options['cnpj'], 'NSU', $ultNSU);
    echo $response;
} catch (Exception $e) {
    echo json_encode(['error' => $e->getMessage()]);
    exit(1);
}