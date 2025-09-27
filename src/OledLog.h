// OledLog.h
#pragma once
#include <Adafruit_SSD1306.h>
#include "OledTerminal.h"

// pinos e endereço do OLED
#ifndef OLED_SDA
#define OLED_SDA   21
#endif
#ifndef OLED_SCL
#define OLED_SCL   22
#endif
#ifndef OLED_ADDR
#define OLED_ADDR  0x3C   // troque para 0x3D se o scanner mostrar
#endif

// instâncias globais de display e “terminal”
extern Adafruit_SSD1306 OLED_DISPLAY;
extern OledTerminal      OLED_TERM;

// inicializa Wire + OLED + terminal
void oledlog_setup();
void oledlog_tick();   // << novo
