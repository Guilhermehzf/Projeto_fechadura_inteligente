#ifndef LCD_INTERFACE_H
#define LCD_INTERFACE_H

#include <Arduino.h>

void setupLcd();
void exibirMensagemInicial();
void exibirModoProgramacao(const String& bufferNovaSenha);
void exibirAcessoLiberado();
void exibirTrancado();
void verificarTimeoutMensagem();

#endif