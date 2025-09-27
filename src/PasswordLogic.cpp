#include "PasswordLogic.h"
#include "EEPROMHandler.h"
#include "LCDInterface.h"
#include "LockControl.h"
#include "secrets.h"

// buffers e estados só aqui
static String senhaAtual;
static const String senhaMestra = MASTER_PASSWORD;
static String bufferNormal;
static String bufferProg;
static bool   modoProg = false;

void password_init()
{
  senhaAtual = lerSenhaDaEeprom();
  bufferNormal = "";
  bufferProg   = "";
  modoProg     = false;
  exibirMensagemInicial();
  exibirDigitacaoNormal(bufferNormal);   // limpa a 2ª linha
}

static void entrarModoProg()
{
  modoProg = true;
  bufferProg = "";
  exibirModoProgramacao(bufferProg);
}

static void sairModoProg()
{
  modoProg = false;
  bufferProg = "";
  exibirMensagemInicial();
  exibirDigitacaoNormal(bufferNormal); // garante linha limpa ao sair
}

static void tratarNormal(char k)
{
  // backspace no modo normal: usa '*'
  if (k == '*') {
    if (bufferNormal.length() > 0) {
      bufferNormal.remove(bufferNormal.length() - 1);
    }
    exibirDigitacaoNormal(bufferNormal);
    return;
  }

  // ignora '#' no modo normal (sem ação)
  if (k == '#') {
    return;
  }

  // acumula a tecla e atualiza máscara
  bufferNormal += k;
  exibirDigitacaoNormal(bufferNormal);

  // senha mestra → entra programação
  if (bufferNormal.endsWith(senhaMestra)) {
    bufferNormal = "";
    entrarModoProg();
    return;
  }

  // senha normal → alterna tranca
  if (bufferNormal.endsWith(senhaAtual)) {
    bufferNormal = "";
    lock_toggle("keypad");        // atualiza LED/LCD
    exibirMensagemInicial();      // volta tela
    exibirDigitacaoNormal(bufferNormal); // limpa máscara
  }
}

static void tratarProg(char k)
{
  if (k == '#') {
    if (bufferProg.length() > 0) {
      salvarSenhaNaEeprom(bufferProg);
      senhaAtual = bufferProg;
    }
    sairModoProg();
    return;
  }

  if (k == '*') {
    if (bufferProg.length() > 0) bufferProg.remove(bufferProg.length() - 1);
    exibirModoProgramacao(bufferProg);
    return;
  }

  if (bufferProg.length() < 10) {
    bufferProg += k;
    exibirModoProgramacao(bufferProg);
  }
}

void password_onKey(char k)
{
  if (!k) return;
  if (modoProg) tratarProg(k);
  else          tratarNormal(k);
}
