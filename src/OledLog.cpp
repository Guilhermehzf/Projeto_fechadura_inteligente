#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include "OledTerminal.h"
#include "OledLog.h"
#include <WiFi.h>  // Necessário para acessar o WiFi

Adafruit_SSD1306 OLED_DISPLAY(128, 64, &Wire, -1);
OledTerminal OLED_TERM(&OLED_DISPLAY);

void oledlog_setup() {
  if (!OLED_DISPLAY.begin(SSD1306_SWITCHCAPVCC, OLED_ADDR)) {
    return; // segue sem OLED
  }
  OLED_TERM.begin();
  OLED_TERM.println("OLED log pronto");
}

// Desenha o sinal de Wi-Fi no canto superior direito com barras crescentes
void desenharSinalWifi() {
  long rssi = WiFi.RSSI();
  int signalLevel = 0;

  if (rssi >= -50) signalLevel = 4;
  else if (rssi >= -60) signalLevel = 3;
  else if (rssi >= -70) signalLevel = 2;
  else if (rssi >= -80) signalLevel = 1;
  else signalLevel = 0;

  const int iconBaseX = 108;
  const int iconBaseY = 10;
  const int barWidth = 3;
  const int barGap = 1;
  int barHeights[] = { 2, 4, 6, 8, 10 };

  // CORREÇÃO AQUI: Era SSD136_BLACK, agora é SSD1306_BLACK
  OLED_DISPLAY.fillRect(iconBaseX, 0, 128 - iconBaseX, iconBaseY + 1, SSD1306_BLACK);

  for (int i = 0; i < 5; i++) {
    if (i <= signalLevel) {
      int currentBarX = iconBaseX + (i * (barWidth + barGap));
      OLED_DISPLAY.fillRect(currentBarX, iconBaseY - barHeights[i], barWidth, barHeights[i], SSD1306_WHITE);
    }
  }
}

// Função que é chamada a cada loop
void oledlog_tick() {
  // 1. Sua biblioteca de terminal provavelmente desenha o texto e atualiza a tela.
  OLED_TERM.tick();
  // 2. Agora desenhamos o sinal de Wi-Fi no buffer (por cima do que já estava lá).
  desenharSinalWifi();
  // 3. ESSENCIAL: Enviamos o buffer atualizado (com o texto E o ícone) para a tela.
  OLED_DISPLAY.display(); 
}
