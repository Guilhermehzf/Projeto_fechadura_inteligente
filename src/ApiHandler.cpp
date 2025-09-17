#include "ApiHandler.h"
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "main.h"         // Para acessar a variável global trancaAberta
#include "LedControl.h"   
#include "LCDInterface.h" 

// --- Configurações da API ---
// COLOQUE AQUI A SUA URL ATUAL DO NGROK!
const char* baseUrl = "http://arduino.ghzds.com.br";

const unsigned long POLLING_INTERVAL = 1000; // 1 segundos
unsigned long ultimaVerificacao = 0;

void processarRespostaApi(String payload) {
    JsonDocument doc;
    DeserializationError error = deserializeJson(doc, payload);

    if (error) {
        Serial.print("Falha ao analisar JSON: ");
        Serial.println(error.c_str());
        return;
    }

    bool estadoDoServidor = doc["tranca_aberta"];

    if (estadoDoServidor != trancaAberta) {
        Serial.println("Estado dessincronizado, atualizando...");
        trancaAberta = estadoDoServidor;

        atualizarLeds(trancaAberta);
        if (trancaAberta) {
            exibirAcessoLiberado();
        } else {
            exibirTrancado();
        }
    } else {
        Serial.println("Estado já sincronizado.");
    }
}

// Função de polling para sincronizar com o servidor
void api_loop() {
    if (millis() - ultimaVerificacao >= POLLING_INTERVAL) {
        ultimaVerificacao = millis();
        Serial.println("Verificando estado na API via GET /status...");
        
        String statusUrl = String(baseUrl) + "/status";
        Serial.println("URL: " + statusUrl);
        HTTPClient http;
        http.begin(statusUrl);
        int httpCode = http.GET();
        Serial.println("HTTP Code: " + String(httpCode));

        if (httpCode == HTTP_CODE_OK) {
            processarRespostaApi(http.getString());
        } else {
            Serial.printf("[HTTP] GET... falhou, erro: %d\n", httpCode);
        }
        http.end();
    }
}

// Função para comandar a mudança de estado
void api_toggle_state() {
    Serial.println("Enviando comando para API via POST /toggle...");
    
    String toggleUrl = String(baseUrl) + "/toggle";
    HTTPClient http;
    http.begin(toggleUrl);
    int httpCode = http.POST("");

    if (httpCode == HTTP_CODE_OK) {
        processarRespostaApi(http.getString());
    } else {
        Serial.printf("[HTTP] POST... falhou, erro: %d\n", httpCode);
    }
    http.end();
}