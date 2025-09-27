// src/MqttHandler.h
#ifndef MQTTHANDLER_H
#define MQTTHANDLER_H

// Inicializa cliente e configura callbacks/servidor
void mqtt_setup();

// Avança a FSM do MQTT (reconexão, loop, etc.). Chamar em loop()
void mqtt_tick();

// Publica o estado atual (usa retained). Se offline, marca como pendente
void publish_current_state();

#endif
