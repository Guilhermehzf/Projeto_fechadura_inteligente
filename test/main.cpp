#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <LiquidCrystal_I2C.h>

// Definições dos displays
#define I2C_SDA 21
#define I2C_SCL 22
#define OLED_ADDRESS 0x3C
#define LCD_ADDRESS  0x27 // Ou 0x26, dependendo do seu hardware

// Criação dos objetos
Adafruit_SSD1306 display(128, 64, &Wire, -1);
LiquidCrystal_I2C lcd(LCD_ADDRESS, 16, 2);

void setup() {
  Serial.begin(115200);

  // Inicia o barramento I2C
  Wire.begin(I2C_SDA, I2C_SCL);

  // Inicializa o OLED
  if(!display.begin(SSD1306_SWITCHCAPVCC, OLED_ADDRESS)) {
    Serial.println(F("Falha ao alocar SSD1306"));
    for(;;); // Loop infinito se falhar
  }

  // Inicializa o LCD
  lcd.init();
  lcd.backlight();
  
  // Mensagens iniciais
  lcd.setCursor(0, 0);
  lcd.print("LCD Pronto!");
  
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(WHITE);
  display.setCursor(0,0);
  display.println("OLED Pronto!");
  display.display();
}

void loop() {
  // Exemplo de atualização das telas
  
  // Atualiza o LCD com um contador
  lcd.setCursor(0, 1);
  lcd.print("Contador: ");
  lcd.print(millis() / 1000);

  // Atualiza o OLED com outra informação
  display.setCursor(0, 20);
  display.print("RSSI Wi-Fi: ");
  display.display();

  delay(500);
}