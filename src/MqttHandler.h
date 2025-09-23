// src/MqttHandler.h
#ifndef MQTTHANDLER_H
#define MQTTHANDLER_H

void setup_mqtt();
void mqtt_loop();
void publish_current_state(); 
void mqtt_tick();
void mqtt_setup();

#endif