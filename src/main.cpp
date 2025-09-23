#include "main.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "KeypadHandler.h"
#include "EEPROMHandler.h"
#include "WiFiConfig.h"
#include "MqttHandler.h"
#include "LockControl.h"
#include "PasswordLogic.h"

void setup() {
  Serial.begin(115200);

  // hardware primeiro (sempre offline-friendly)
  setupLeds();
  setupLcd();
  setupEeprom();

  lock_init(true);       // mostra LED/LCD conforme estado inicial
  password_init();       // carrega senha + prepara buffers

  // rede NÃO-bloqueante
  wifi_setup_nonblocking();
  mqtt_setup();
}

void loop() {
  // rede em background
  wifi_tick();
  mqtt_tick();

  // teclado SEMPRE funciona
  char k = lerTecla();
  if (k) password_onKey(k);

  verificarTimeoutMensagem(); // se você usa timeouts na LCD
}
