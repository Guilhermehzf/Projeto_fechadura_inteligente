#include "LockControl.h"
#include "LedControl.h"
#include "LCDInterface.h"
#include "MqttHandler.h"   // publish_current_state()

bool trancaAberta = true;

static void render_lock()
{
  atualizarLeds(trancaAberta);
  if (trancaAberta) exibirAcessoLiberado();
  else              exibirTrancado();
}

void lock_init(bool abertaInicial)
{
  trancaAberta = abertaInicial;
  render_lock();
}

void lock_apply(bool aberta, const char* /*method*/)
{
  trancaAberta = aberta;
  render_lock();

  // publica se estiver conectado (ou marca pendente)
  publish_current_state();
}

void lock_toggle(const char* /*method*/)
{
  lock_apply(!trancaAberta);
}
