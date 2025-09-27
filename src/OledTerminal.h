#pragma once
#include <Arduino.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <Print.h>
#include <stdarg.h>

class OledTerminal : public Print {
public:
  OledTerminal(Adafruit_SSD1306* d, int cols=21, int rows=8)
    : disp(d), COLS(cols), ROWS(rows), col(0), row(0),
      autoscroll(true), lastRedrawMs(0), redrawIntervalMs(80), pending(false)
  {
    buffer = new char[COLS * ROWS];
    clearBuffer();
  }

  ~OledTerminal() { delete [] buffer; }

  void begin() {
    disp->clearDisplay();
    disp->setTextColor(SSD1306_WHITE);
    disp->setTextSize(1);                // 6x8 px/char → ~21x8 em 128x64
    redrawNow();                         // desenha 1x
  }

  void clear() {
    clearBuffer();
    col = row = 0;
    redrawSoon();                        // agenda redraw
  }

  void setAutoscroll(bool on) { autoscroll = on; }

  // Chame periodicamente (ex.: a cada loop) para aplicar redraw se necessário
  void tick() {
    if (pending && (millis() - lastRedrawMs >= redrawIntervalMs)) {
      redrawNow();
    }
  }

  // Print API
  size_t write(uint8_t c) override {
    if (c == '\r') return 1;
    if (c == '\n') { newline(); redrawSoon(); return 1; }
    if (col >= COLS) newline();
    cell(row, col) = (char)c;
    col++;
    // não redesenha aqui; deixa para tick()
    redrawSoon();
    return 1;
  }
  using Print::print;
  using Print::println;

  void printf(const char* fmt, ...) {
    char tmp[256];
    va_list ap; va_start(ap, fmt);
    vsnprintf(tmp, sizeof(tmp), fmt, ap);
    va_end(ap);
    print(tmp);          // escreve no buffer
    redrawSoon();        // agenda redraw (tick() executa)
  }

private:
  Adafruit_SSD1306* disp;
  const int COLS, ROWS;
  int col, row;
  bool autoscroll;
  char* buffer; // ROWS x COLS

  unsigned long lastRedrawMs;
  unsigned long redrawIntervalMs; // ~80ms
  bool pending;

  inline char& cell(int r, int c) { return buffer[r*COLS + c]; }

  void clearBuffer() { memset(buffer, ' ', COLS*ROWS); }

  void scrollUp() {
    memmove(buffer, buffer + COLS, COLS*(ROWS-1));
    memset(buffer + COLS*(ROWS-1), ' ', COLS);
  }

  void newline() {
    col = 0;
    if (row < ROWS - 1) row++;
    else if (autoscroll) scrollUp();
  }

  void redrawSoon() { pending = true; }

  void redrawNow() {
    pending = false;
    disp->clearDisplay();
    for (int r = 0; r < ROWS; r++) {
      disp->setCursor(0, r * 8);
      disp->write((const uint8_t*)&cell(r,0), COLS);
    }
    disp->display();
    lastRedrawMs = millis();
  }
};
