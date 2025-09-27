#include "LCDInterface.h"
#include <LiquidCrystal_I2C.h>

LiquidCrystal_I2C lcd(0x27, 16, 2);

static bool mostrandoMensagemTemporaria = false;
static unsigned long tempoInicioMensagem = 3000; // 3S
static const long TIMEOUT_MENSAGEM = 5000; // 5s

static void clearLine(uint8_t row) {
  lcd.setCursor(0, row);
  for (int i = 0; i < 16; i++) lcd.print(' ');
  lcd.setCursor(0, row);
}

void setupLcd() {
  lcd.init();
  lcd.backlight();
}

void exibirMensagemInicial() {
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Digite a senha:");
  clearLine(1); // linha de digitação começa limpa
}

void exibirDigitacaoNormal(const String& buffer) {
  // máscara **** na segunda linha
  lcd.setCursor(0, 1);
  clearLine(1);
  uint8_t n = buffer.length();
  if (n > 16) n = 16;
  for (uint8_t i = 0; i < n; i++) lcd.print('*');
}

void exibirModoProgramacao(const String& bufferNovaSenha) {
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Modo programacao");
  lcd.setCursor(0, 1);
  lcd.print("Nova: ");
  uint8_t n = bufferNovaSenha.length();
  if (n > 10) n = 10; // limite visual
  for (uint8_t i = 0; i < n; i++) lcd.print('*');
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
  if (mostrandoMensagemTemporaria && (millis() - tempoInicioMensagem >= TIMEOUT_MENSAGEM)) {
    exibirMensagemInicial();
    mostrandoMensagemTemporaria = false;
  }
}

bool lcdEstaMostrandoMensagem() {
  return mostrandoMensagemTemporaria;
}
