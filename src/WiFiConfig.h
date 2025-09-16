#ifndef WIFICONFIG_H
#define WIFICONFIG_H

// Configura o WiFi. Se não houver credenciais salvas, 
// entra em modo de configuração (AP).
// A função bloqueia a execução até que a conexão seja estabelecida.
void setupWiFi();

#endif