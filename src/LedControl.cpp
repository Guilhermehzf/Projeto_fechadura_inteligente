#include "LedControl.h"
#include <Arduino.h>

void setupLeds() {
  pinMode(RELAY_PIN, OUTPUT);
  // estado inicial: ABERTO = LOW
  digitalWrite(RELAY_PIN, LOW);
}

void atualizarLeds(bool trancaAberta) {
  // LOW = aberto; HIGH = trancado
  digitalWrite(RELAY_PIN, trancaAberta ? LOW : HIGH);
}