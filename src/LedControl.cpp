#include "LedControl.h"
#include <Arduino.h>

const int ledVerde = 25;
const int ledVermelho = 27;

void setupLeds() {
    pinMode(ledVerde, OUTPUT);
    pinMode(ledVermelho, OUTPUT);
}

void atualizarLeds(bool trancaAberta) {
    digitalWrite(ledVerde, trancaAberta ? HIGH : LOW);
    digitalWrite(ledVermelho, trancaAberta ? LOW : HIGH);
}