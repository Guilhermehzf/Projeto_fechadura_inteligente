#ifndef LOCKCONTROL_H
#define LOCKCONTROL_H

#include <Arduino.h>

// Estado global exposto (usado pelo MQTT/Password)
extern bool trancaAberta;

// Pino do relé/fechadura (ajuste se necessário para seu hardware)
#ifndef RELAY_PIN
#define RELAY_PIN 26
#endif

// Inicializa estado e hardware (LED, LCD e relé)
void lock_init(bool abertaInicial);

// Aplica estado (atualiza LED, LCD, relé) e publica no MQTT
void lock_apply(bool aberta, const char* method = "local");

// Alterna o estado atual
void lock_toggle(const char* method = "local");

#endif
