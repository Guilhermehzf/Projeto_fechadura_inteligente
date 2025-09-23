#ifndef PASSWORDLOGIC_H
#define PASSWORDLOGIC_H

#include <Arduino.h>

void password_init();        // lê senha da EEPROM e prepara buffers
void password_onKey(char k); // tratar tecla (normal/programação)

#endif
