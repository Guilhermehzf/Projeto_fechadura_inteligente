#include "WiFiConfig.h"
#include <WiFi.h>

// DESCOMENTE A LINHA ABAIXO PARA RODAR NO SIMULADOR WOKWI
#define WOKWI_SIMULATION

#ifdef WOKWI_SIMULATION
  // --- CONFIGURAÇÃO PARA O SIMULADOR WOKWI ---
  const char* ssid = "Wokwi-GUEST";
  const char* password = "";
#else
  // --- CONFIGURAÇÃO PARA SUA PLACA REAL ---
  // --- COLOQUE AS CREDENCIAIS DA SUA REDE WI-FI AQUI ---
  const char* ssid = "NOME_DA_SUA_REDE_WIFI";
  const char* password = "SENHA_DA_SUA_REDE_WIFI";
#endif

void setup_wifi() {
  delay(10);
  Serial.println();
  Serial.print("Conectando-se a ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.println("WiFi conectado!");
  Serial.print("Endereço IP: ");
  Serial.println(WiFi.localIP());
}