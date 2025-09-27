#include "main.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "KeypadHandler.h"
#include "EEPROMHandler.h"
#include "WiFiConfig.h"
#include "MqttHandler.h"
#include "LockControl.h"
#include "PasswordLogic.h"

#include "Log.h"
#include "OledLog.h"

void setup() {
  Serial.begin(115200);

  // Configura o OLED
  oledlog_setup();
  
  // Tee dos logs: Serial + OLED
  LOG.begin(&Serial, &OLED_TERM);
  LOG.println("Logger: Serial + OLED");

  // Resto do setup
  setupLeds();
  setupLcd();
  setupEeprom();
  lock_init(true);
  password_init();

  wifi_setup_nonblocking();
  mqtt_setup();
}

void loop() {
  oledlog_tick();  // Atualiza o OLED com log e Wi-Fi
  
  wifi_tick();
  mqtt_tick();

  char k = lerTecla();
  if (k) password_onKey(k);

  verificarTimeoutMensagem();
}
