#include "main.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "KeypadHandler.h"
#include "EEPROMHandler.h"
#include "WiFiConfig.h"
#include "MqttHandler.h"

// --- Variáveis Globais ---
String senhaAtual;
const String senhaMestra = "0000";
String bufferNormal = "";
String bufferProgramacao = "";
bool modoProgramacao = false;
bool trancaAberta = true; 

// --- Protótipos ---
void tratarModoProgramacao(char tecla);
void tratarModoNormal(char tecla);

void setup() {
    Serial.begin(115200);
    setupLeds();
    setupLcd();
    setupEeprom();
    setup_wifi();
    setup_mqtt();
    senhaAtual = lerSenhaDaEeprom();
    atualizarLeds(trancaAberta);
    exibirMensagemInicial();
}

void loop() {
    verificarTimeoutMensagem();
    mqtt_loop();
    char tecla = lerTecla();
    if (tecla) {
        if (modoProgramacao) {
            tratarModoProgramacao(tecla);
        } else {
            tratarModoNormal(tecla);
        }
    }
}

// === FUNÇÃO CORRIGIDA ===
void tratarModoNormal(char tecla) {
    // A LINHA QUE FALTAVA FOI ADICIONADA AQUI
    bufferNormal += tecla;

    if (bufferNormal.endsWith(senhaMestra)) {
        bufferNormal = "";
        modoProgramacao = true;
        exibirModoProgramacao(bufferProgramacao);
        return;
    }

    if (bufferNormal.endsWith(senhaAtual)) {
        bufferNormal = "";
        
        // 1. AÇÃO IMEDIATA: Mude o estado e o hardware localmente
        Serial.println("Senha correta! Trocando estado localmente...");
        trancaAberta = !trancaAberta; 
        atualizarLeds(trancaAberta);
        if (trancaAberta) {
            exibirAcessoLiberado();
        } else {
            exibirTrancado();
        }

        // 2. NOTIFICAÇÃO: Avise a nuvem sobre o novo estado
        publish_current_state();
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