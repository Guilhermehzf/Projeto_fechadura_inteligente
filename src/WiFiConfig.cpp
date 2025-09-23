#include "WiFiConfig.h"
#include <WiFi.h>
#include "secrets.h"

// DESCOMENTE para Wokwi
#define WOKWI_SIMULATION

#ifdef WOKWI_SIMULATION
  static const char* ssid = "Wokwi-GUEST";
  static const char* pass = "";
#else
  static const char* ssid = SECRET_WIFI_SSID;
  static const char* pass = SECRET_WIFI_PASS;
#endif

enum class WifiState { Idle, Connecting, Connected };
static WifiState st = WifiState::Idle;
static unsigned long lastTry = 0;
static const unsigned long RETRY_MS = 3000;
static bool everPrintIP = false;

void wifi_setup_nonblocking()
{
  WiFi.mode(WIFI_STA);
  st = WifiState::Idle;
  everPrintIP = false;
}

void wifi_tick()
{
  unsigned long now = millis();

  if (st == WifiState::Idle) {
    if (now - lastTry >= 10) {
      lastTry = now;
      WiFi.begin(ssid, pass);
      st = WifiState::Connecting;
      Serial.printf("[WiFi] Conectando em '%s'...\n", ssid);
    }
    return;
  }

  if (st == WifiState::Connecting) {
    wl_status_t s = WiFi.status();
    if (s == WL_CONNECTED) {
      st = WifiState::Connected;
      if (!everPrintIP) {
        everPrintIP = true;
        Serial.printf("[WiFi] Conectado, IP=%s\n", WiFi.localIP().toString().c_str());
      }
    } else if (s == WL_CONNECT_FAILED || s == WL_NO_SSID_AVAIL || s == WL_DISCONNECTED) {
      if (now - lastTry >= RETRY_MS) {
        lastTry = now;
        WiFi.disconnect(true);
        WiFi.begin(ssid, pass);
        Serial.println("[WiFi] Re-tentando...");
      }
    }
    return;
  }

  if (st == WifiState::Connected) {
    if (WiFi.status() != WL_CONNECTED) {
      st = WifiState::Idle; // caiu → recomeça
      Serial.println("[WiFi] Desconectado.");
    }
  }
}

bool wifi_is_connected()
{
  return WiFi.status() == WL_CONNECTED;
}
