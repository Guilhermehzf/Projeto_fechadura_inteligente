#include "LCDInterface.h"
#include <LiquidCrystal_I2C.h>

// O objeto do LCD é "privado" deste módulo.
LiquidCrystal_I2C lcd(0x27, 16, 2);

bool mostrandoMensagemTemporaria = false;
unsigned long tempoInicioMensagem = 0;
const long TIMEOUT_MENSAGEM = 5000; // 5 segundos

void setupLcd() {
    lcd.init();
    lcd.backlight();
}

void exibirMensagemInicial() {
    lcd.clear();
    lcd.print("Digite a senha:");
}

void exibirModoProgramacao(const String& bufferNovaSenha) {
    lcd.clear();
    lcd.print("Modo programacao");
    lcd.setCursor(0, 1);
    lcd.print("Nova: ");
    for (unsigned int i = 0; i < bufferNovaSenha.length(); i++) {
        lcd.print("*");
    }
}

void exibirAcessoLiberado() {
    lcd.clear();
    lcd.print("Acesso liberado");
    mostrandoMensagemTemporaria = true;
    tempoInicioMensagem = millis();
}

void exibirTrancado() {
    lcd.clear();
    lcd.print("Trancado");
    mostrandoMensagemTemporaria = true;
    tempoInicioMensagem = millis();
}

void verificarTimeoutMensagem() {
    // Se uma mensagem temporária está ativa e o tempo expirou...
    if (mostrandoMensagemTemporaria && (millis() - tempoInicioMensagem >= TIMEOUT_MENSAGEM)) {
        exibirMensagemInicial();
        mostrandoMensagemTemporaria = false;
    }
}