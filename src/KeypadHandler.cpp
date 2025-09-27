#include "KeypadHandler.h"
#include <Keypad.h>
#include "Log.h"

const byte ROWS = 4;
const byte COLS = 3;

char keys[ROWS][COLS] = {
  {'1','2','3'},
  {'4','5','6'},
  {'7','8','9'},
  {'*','0','#'}
};

// 🔧 TROQUE a ordem das colunas: era {16, 4, 2}
byte rowPins[ROWS] = {19, 18, 5, 17};
byte colPins[COLS] = { 33, 32, 25};  // ← invertido (da esquerda p/ direita)

Keypad kp = Keypad(makeKeymap(keys), rowPins, colPins, ROWS, COLS);

char lerTecla() {
  char k = kp.getKey();
  if (k) LOGF("[KEYPAD] key='%c'\n", k); // debug
  return k;
}
