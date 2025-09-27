#include "MqttHandler.h"
#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include "main.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "secrets.h"

#include "Log.h"

// Credenciais
static const char* MQTT_BROKER = SECRET_MQTT_BROKER;
static const int   MQTT_PORT   = SECRET_MQTT_PORT;
static const char* MQTT_USER   = SECRET_MQTT_USER;
static const char* MQTT_PASS   = SECRET_MQTT_PASS;

// Tópicos
static const char* MQTT_TOPIC_COMMANDS = "fechadura/comandos";
static const char* MQTT_TOPIC_STATE    = "fechadura/estado";

static WiFiClientSecure espClient;
static PubSubClient client(espClient);

// Publicação pendente quando offline
static bool publishPending = false;

// Backoff de reconexão
static unsigned long lastMqttAttempt = 0;
static const unsigned long MQTT_RETRY_MS = 3000;

static void mqtt_on_message(char* topic, byte* payload, unsigned int length) {
  LOGF("[MQTT] RX topic=%s len=%u\n", topic, length);

  StaticJsonDocument<128> doc;
  DeserializationError err = deserializeJson(doc, payload, length);
  if (err) {
    LOGF("[MQTT] JSON inválido: %s\n", err.c_str());
    return;
  }

  const char* command = doc["command"];
  if (command && strcmp(command, "toggle") == 0) {
    trancaAberta = !trancaAberta;
    atualizarLeds(trancaAberta);
    if (trancaAberta) exibirAcessoLiberado();
    else              exibirTrancado();

    // confirma novo estado
    publish_current_state();
  }
}

void publish_current_state() {
  // Monta JSON { "tranca_aberta": <bool> }
  StaticJsonDocument<64> doc;
  doc["tranca_aberta"] = trancaAberta;

  char buf[64];
  size_t n = serializeJson(doc, buf, sizeof(buf));

  if (client.connected()) {
    bool ok = client.publish(
      MQTT_TOPIC_STATE,
      reinterpret_cast<const uint8_t*>(buf),
      static_cast<unsigned int>(n),
      true
    );// retained
    LOGF("[MQTT] publish estado '%s': %s\n", buf, ok ? "OK" : "FAIL");
    if (!ok) publishPending = true;
  } else {
    publishPending = true; // envia quando reconectar
  }
}

static void mqtt_connect_once() {
  if (client.connected()) return;
  if (WiFi.status() != WL_CONNECTED) return;

  String cid = "esp32-fechadura-" + WiFi.macAddress();
  if (client.connect(cid.c_str(), MQTT_USER, MQTT_PASS)) {
    LOG.println("[MQTT] conectado!");
    client.subscribe(MQTT_TOPIC_COMMANDS);

    // publica retained na conexão (ou o que ficou pendente)
    if (publishPending) {
      publishPending = false;
      publish_current_state();
    } else {
      publish_current_state();
    }
  } else {
    LOGF("[MQTT] falhou rc=%d\n", client.state());
  }
}

void mqtt_setup() {
  espClient.setInsecure();               // TLS sem verificação de CA
  client.setServer(MQTT_BROKER, MQTT_PORT);
  client.setCallback(mqtt_on_message);
}

void mqtt_tick() {
  if (!client.connected()) {
    unsigned long now = millis();
    if (now - lastMqttAttempt >= MQTT_RETRY_MS) {
      lastMqttAttempt = now;
      mqtt_connect_once();
    }
    return;
  }
  client.loop();
}
