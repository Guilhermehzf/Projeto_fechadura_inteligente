#ifndef API_HANDLER_H
#define API_HANDLER_H

// Função de polling, para ser chamada no loop principal.
void api_loop();

// Função para ser chamada quando a senha correta for digitada.
void api_toggle_state();

#endif