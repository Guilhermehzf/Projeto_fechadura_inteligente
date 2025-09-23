#ifndef LOCKCONTROL_H
#define LOCKCONTROL_H

#include <Arduino.h>

// exposto para outros módulos (MQTT/handlers)
extern bool trancaAberta;

// inicializa estado e hardware
void lock_init(bool abertaInicial);

// aplica estado (atualiza LED/LCD) e agenda publish (se possível)
void lock_apply(bool aberta, const char* method = "local");

// alterna estado atual
void lock_toggle(const char* method = "local");

#endif
