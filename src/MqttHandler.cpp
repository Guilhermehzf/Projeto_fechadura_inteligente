// src/MqttHandler.cpp
#include "MqttHandler.h"
#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include "main.h"         
#include "LedControl.h"   
#include "LCDInterface.h"
#include <ArduinoJson.h>
#include "secrets.h"

// --- SUAS CREDENCIAIS DO HIVEMQ ---
const char* MQTT_BROKER = SECRET_MQTT_BROKER;
const int   MQTT_PORT = SECRET_MQTT_PORT;
const char* MQTT_USER = SECRET_MQTT_USER;
const char* MQTT_PASS = SECRET_MQTT_PASS;

// --- TÓPICOS ---
const char* MQTT_TOPIC_COMMANDS = "fechadura/comandos";
const char* MQTT_TOPIC_STATE = "fechadura/estado";

WiFiClientSecure espClient;
PubSubClient client(espClient);

void publish_current_state() {
  if (client.connected()) {
    // 1. Cria um documento JSON para montar a mensagem
    JsonDocument doc;
    doc["tranca_aberta"] = trancaAberta;

    // 2. Converte o JSON para uma string
    String output;
    serializeJson(doc, output);

    // 3. Publica a string JSON no tópico de estado COM A FLAG DE RETENÇÃO (true)
    client.publish(MQTT_TOPIC_STATE, output.c_str(), true);
    Serial.printf("Estado JSON (%s) publicado (com retenção) no tópico %s\n", output.c_str(), MQTT_TOPIC_STATE);
  }
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.printf("Mensagem recebida no tópico %s\n", topic);
  JsonDocument doc;
  DeserializationError error = deserializeJson(doc, payload, length);

  if (error) {
    Serial.printf("Falha ao analisar JSON: %s\n", error.c_str());
    return;
  }

  const char* command = doc["command"];
  if (command && strcmp(command, "toggle") == 0) {
    Serial.println("Comando 'toggle' remoto recebido! Alterando estado.");
    trancaAberta = !trancaAberta;
    
    atualizarLeds(trancaAberta);
    if (trancaAberta) {
      exibirAcessoLiberado();
    } else {
      exibirTrancado();
    }
    
    // Após um comando remoto, publicamos o novo estado como confirmação
    publish_current_state();
  }
}

void reconnect() {
  while (!client.connected()) {
    Serial.print("Tentando conectar ao Broker MQTT...");
    String clientId = "esp32-fechadura-" + WiFi.macAddress();
    if (client.connect(clientId.c_str(), MQTT_USER, MQTT_PASS)) {
      Serial.println("conectado!");
      client.subscribe(MQTT_TOPIC_COMMANDS);
      publish_current_state();
    } else {
      Serial.printf("falhou, rc=%d tentando novamente em 5 segundos\n", client.state());
      delay(5000);
    }
  }
}

void setup_mqtt() {
  espClient.setInsecure();
  client.setServer(MQTT_BROKER, MQTT_PORT);
  client.setCallback(callback);
}

void mqtt_loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();
}