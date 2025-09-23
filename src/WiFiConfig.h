#ifndef WIFICONFIG_H
#define WIFICONFIG_H

void wifi_setup_nonblocking();  // não trava no setup
void wifi_tick();               // chama no loop()
bool wifi_is_connected();

#endif
