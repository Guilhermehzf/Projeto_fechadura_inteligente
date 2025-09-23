#include "EEPROMHandler.h"
#include <EEPROM.h>
#include "secrets.h"

const int EEPROM_START_ADDR = 0;
const int EEPROM_SIZE = 64;
const String SENHA_PADRAO = PASSWORD;
const int TAMANHO_MAX_SENHA = 10;

void setupEeprom() {
    EEPROM.begin(EEPROM_SIZE);
}

String lerSenhaDaEeprom() {
    String senha = "";
    for (int i = 0; i < TAMANHO_MAX_SENHA; i++) {
        char c = EEPROM.read(EEPROM_START_ADDR + i);
        if (c == '\0' || c == 0xFF) { // Terminador nulo ou byte vazio
            break;
        }
        senha += c;
    }
    return senha.length() > 0 ? senha : SENHA_PADRAO;
}

void salvarSenhaNaEeprom(const String& senha) {
    unsigned int i = 0;
    for (i = 0; i < senha.length() && i < TAMANHO_MAX_SENHA; i++) {
        EEPROM.write(EEPROM_START_ADDR + i, senha[i]);
    }
    // Adiciona o terminador nulo para marcar o fim da string
    EEPROM.write(EEPROM_START_ADDR + i, '\0'); 
    EEPROM.commit();
}