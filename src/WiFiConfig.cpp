#include "WiFiConfig.h"
#include <WiFi.h>

// DESCOMENTE A LINHA ABAIXO QUANDO ESTIVER RODANDO NO SIMULADOR WOKWI
#define WOKWI_SIMULATION

#ifndef WOKWI_SIMULATION
  // Este código só será compilado para o hardware real
  #include <WiFiManager.h>
#endif

void setupWiFi() {

  #ifdef WOKWI_SIMULATION
    // --- CÓDIGO PARA O SIMULADOR WOKWI ---
    Serial.println("Modo de simulação Wokwi: Conectando à rede virtual...");
    WiFi.begin("Wokwi-GUEST", "", 6); // O 6 é o canal, pode ser necessário

    while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
    }
    Serial.println("");
    Serial.println("Conectado à rede Wokwi com sucesso!");
    Serial.print("Endereço IP: ");
    Serial.println(WiFi.localIP());

  #else
    // --- CÓDIGO PARA O HARDWARE REAL (com WiFiManager) ---
    WiFiManager wm;
    wm.setConfigPortalTimeout(180);

    if (!wm.autoConnect("Fechadura-Config", "senha1234")) {
      Serial.println("Falha ao conectar e o tempo limite expirou. Reiniciando...");
      delay(3000);
      ESP.restart();
    }
    
    Serial.println("");
    Serial.println("Conectado à rede Wi-Fi com sucesso!");
    Serial.print("Endereço IP: ");
    Serial.println(WiFi.localIP());

  #endif
}