#ifndef LCD_INTERFACE_H
#define LCD_INTERFACE_H

#include <Arduino.h>

void setupLcd();
void exibirMensagemInicial();
void exibirModoProgramacao(const String& bufferNovaSenha);
void exibirAcessoLiberado();
void exibirTrancado();
void verificarTimeoutMensagem();
bool lcdEstaMostrandoMensagem();

// mostra na 2ª linha a senha digitada no modo normal, mascarada com '*'
void exibirDigitacaoNormal(const String& buffer);

// opcional: telinha de boot com infos rápidas
void exibirInfoRede(const char* broker, int port, int relayPin);

#endif
