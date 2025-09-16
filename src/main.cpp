#include "main.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "KeypadHandler.h"
#include "EEPROMHandler.h"
#include "WiFiConfig.h"
#include "ApiHandler.h"

// A definição da variável continua aqui
bool trancaAberta = true; 

// Demais variáveis globais...
String senhaAtual;
const String senhaMestra = "0000";
String bufferNormal = "";
String bufferProgramacao = "";
bool modoProgramacao = false;

// Protótipos
void tratarModoProgramacao(char tecla);
void tratarModoNormal(char tecla);

void setup() {
    Serial.begin(115200);
    setupLeds();
    setupLcd();
    setupEeprom();
    setupWiFi();
    senhaAtual = lerSenhaDaEeprom();
    atualizarLeds(trancaAberta);
    exibirMensagemInicial();
}

void loop() {
    verificarTimeoutMensagem();
    api_loop(); // Chama a função de polling

    char tecla = lerTecla();
    if (tecla) {
        if (modoProgramacao) {
            tratarModoProgramacao(tecla);
        } else {
            tratarModoNormal(tecla);
        }
    }
}

void tratarModoNormal(char tecla) {
    bufferNormal += tecla;

    if (bufferNormal.endsWith(senhaMestra)) {
        bufferNormal = "";
        modoProgramacao = true;
        exibirModoProgramacao(bufferProgramacao);
        return;
    }

    if (bufferNormal.endsWith(senhaAtual)) {
        bufferNormal = "";
        // EM VEZ DE MUDAR O ESTADO LOCALMENTE, PEDIMOS PARA A API MUDAR
        api_toggle_state();
    }
}

void tratarModoProgramacao(char tecla) {
    if (tecla == '#') {
        if (bufferProgramacao.length() > 0) {
            salvarSenhaNaEeprom(bufferProgramacao);
            senhaAtual = bufferProgramacao;
        }
        bufferProgramacao = "";
        modoProgramacao = false;
        exibirMensagemInicial();
    } else if (tecla == '*') {
        if (bufferProgramacao.length() > 0) {
            bufferProgramacao.remove(bufferProgramacao.length() - 1);
        }
        exibirModoProgramacao(bufferProgramacao);
    } else {
        if (bufferProgramacao.length() < 10) {
            bufferProgramacao += tecla;
        }
        exibirModoProgramacao(bufferProgramacao);
    }
}