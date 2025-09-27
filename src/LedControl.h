#ifndef LED_CONTROL_H
#define LED_CONTROL_H

// Pino do relé (ajuste se necessário)
#ifndef RELAY_PIN
#define RELAY_PIN 26
#endif

// Inicializa o pino do relé.
void setupLeds();

// Atualiza o relé com base no estado da tranca.
// true = aberto (HIGH), false = trancado (LOW).
void atualizarLeds(bool trancaAberta);

#endif
