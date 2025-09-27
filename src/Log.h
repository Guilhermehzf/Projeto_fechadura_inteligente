#pragma once
#include <Arduino.h>   // Print, size_t
#include <stdarg.h>    // va_list, va_start, va_end
#include <stdio.h>     // ::vsnprintf

class TeeLogger : public Print {
public:
  TeeLogger() : a(nullptr), b(nullptr) {}
  void begin(Print* _a, Print* _b = nullptr) { a = _a; b = _b; }

  size_t write(uint8_t c) override {
    size_t n = 0;
    if (a) n += a->write(c);
    if (b) n += b->write(c);
    return n;
  }
  using Print::print;
  using Print::println;

  void printf(const char* fmt, ...) {
    char buf[256];
    va_list ap;
    va_start(ap, fmt);
    ::vsnprintf(buf, sizeof(buf), fmt, ap);  // <-- note o "::"
    va_end(ap);
    print(buf);
  }

private:
  Print* a;
  Print* b;
};

extern TeeLogger LOG;

#define LOGF(...) do { LOG.printf(__VA_ARGS__); } while(0)
