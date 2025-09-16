#ifndef LED_CONTROL_H
#define LED_CONTROL_H

// Inicializa os pinos dos LEDs.
void setupLeds();

// Atualiza os LEDs com base no estado da tranca.
// true para aberta (verde), false para fechada (vermelho).
void atualizarLeds(bool trancaAberta);

#endif