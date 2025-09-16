#ifndef EEPROM_HANDLER_H
#define EEPROM_HANDLER_H

#include <Arduino.h>

void setupEeprom();
String lerSenhaDaEeprom();
void salvarSenhaNaEeprom(const String& senha);

#endif
