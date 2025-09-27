# 🚪 Fechadura Inteligente com ESP32, MQTT e Duplo Display

![ESP32](https://img.shields.io/badge/ESP32-WROOM-blueviolet)
![Arduino](https://img.shields.io/badge/Framework-Arduino-00979D)
![MQTT](https://img.shields.io/badge/Protocolo-MQTT-red)

## 📖 Sobre o Projeto

Este projeto é um sistema de controle de fechadura inteligente construído com o microcontrolador **ESP32**. Ele oferece controle de acesso tanto localmente, através de um teclado matricial, quanto remotamente, via protocolo MQTT.

A principal característica do sistema é sua **interface de usuário dupla**: um display **LCD 16x2** é usado para interação direta com o usuário (digitação de senha, status), enquanto um display **OLED 128x64** serve como um console de log em tempo real, exibindo status de conexão, mensagens de depuração e um indicador de sinal Wi-Fi.

---

## ✨ Funcionalidades Principais

* ✅ **Controle Duplo:** Abertura por senha local no teclado ou por comando remoto via MQTT.
* ✅ **Interface Dupla:** Display LCD para o usuário e Display OLED para logs e status técnico.
* ✅ **Conectividade IoT:** Conecta-se à rede Wi-Fi e a um broker MQTT para monitoramento e controle à distância.
* ✅ **Persistência de Senha:** A senha de acesso é armazenada na memória EEPROM, mantendo-se mesmo após uma reinicialização.
* ✅ **Modo de Configuração:** Permite a alteração da senha de usuário de forma segura através de uma senha mestra.
* ✅ **Feedback Visual:** LEDs (no módulo relé) indicam claramente o estado da fechadura (trancada/destrancada).

---

## 🛠️ Hardware Necessário

* 1x Placa ESP32 DevKitC V4 (ou similar)
* 1x Display LCD 16x2 com Módulo I2C
* 1x Display OLED 128x64 I2C (SSD1306)
* 1x Teclado de Membrana Matricial 4x3
* 1x Módulo Relé de 1 Canal
* 1x Fechadura Solenoide (12V)
* 1x Fonte de Alimentação Externa 12V com Conector Jack
* 1x Protoboard e Jumpers

---

## 🔌 Esquema do Circuito

A imagem abaixo ilustra a montagem completa do circuito.

**![Insira a imagem do seu esquema aqui](https://i.imgur.com/your-schematic-image.png)**
*(Substitua o link acima pelo caminho da imagem do seu esquema)*

### Tabela de Pinagem

| Componente              | Pino do Componente | Pino no ESP32 |
| :---------------------- | :----------------- | :------------ |
| **Display OLED (I2C)** | `SDA`              | **GPIO 21** |
|                         | `SCL`              | **GPIO 22** |
|                         | `VCC`              | **3.3V** |
|                         | `GND`              | **GND** |
| **Display LCD (I2C)** | `SDA`              | **GPIO 21** |
|                         | `SCL`              | **GPIO 22** |
|                         | `VCC`              | **VIN (5V)** |
|                         | `GND`              | **GND** |
| **Teclado 4x3** | Linhas (R1-R4)     | **19, 18, 5, 17** |
|                         | Colunas (C1-C3)    | **33, 32, 25**|
| **Módulo Relé** | `IN` (Sinal)       | **GPIO 26** |
|                         | `VCC`              | **VIN (5V)** |
|                         | `GND`              | **GND** |

---

## 🚀 Configuração do Ambiente

Para compilar e carregar o código no ESP32, você precisará da IDE do Arduino ou do PlatformIO com as seguintes bibliotecas instaladas:

* `Keypad` by Mark Stanley, Alexander Brevig
* `LiquidCrystal_I2C` by Frank de Brabander
* `Adafruit_GFX` by Adafruit
* `Adafruit_SSD1306` by Adafruit
* `PubSubClient` by Nick O'Leary
* `ArduinoJson` by Benoit Blanchon

### 🔐 Configuração de Credenciais (`secrets.h`)

Para que o projeto funcione, é **essencial** criar um arquivo chamado `secrets.h` dentro da pasta do projeto (ou na pasta `src/`). Este arquivo armazenará todas as suas informações sensíveis, mantendo-as separadas do código principal.

Copie o conteúdo abaixo para o seu arquivo `secrets.h` e **substitua os valores** pelos seus.

```cpp
#ifndef SECRETS_H
#define SECRETS_H

// --- Wi-Fi ---
// Substitua com o nome (SSID) e a senha da sua rede Wi-Fi.
#define SECRET_WIFI_SSID "SUA_REDE_WIFI"
#define SECRET_WIFI_PASS "SUA_SENHA_WIFI"

// --- Broker MQTT ---
// Insira os dados do seu broker MQTT (ex: HiveMQ, Mosquitto).
#define SECRET_MQTT_BROKER "ENDERECO_DO_BROKER"
#define SECRET_MQTT_PORT 8883 // Porta padrão para TLS
#define SECRET_MQTT_USER "SEU_USUARIO_MQTT"
#define SECRET_MQTT_PASS "SUA_SENHA_MQTT"

// --- Senhas do Sistema ---
// Senha padrão que será gravada na EEPROM na primeira inicialização.
#define PASSWORD "123456" 

// Senha mestra para entrar no modo de alteração de senha.
#define MASTER_PASSWORD "9999" 

#endif // SECRETS_H
```

---

## 👨‍💻 Como Usar

1.  **Primeira Utilização:** Ao ligar, o sistema estará no estado "Acesso Liberado" e a senha padrão definida em `PASSWORD` será carregada.
2.  **Operação Normal:**
    * O LCD exibirá "Digite a senha:".
    * Insira a senha de 6 dígitos. A tela mostrará `******`.
    * Se a senha estiver correta, a fechadura alternará seu estado (de trancada para aberta, ou vice-versa).
    * A tecla `*` funciona como "backspace" para apagar o último dígito.
3.  **Alterar a Senha:**
    * Na tela inicial, digite a `MASTER_PASSWORD`.
    * O sistema entrará em "Modo programacao".
    * Digite a nova senha (de até 10 dígitos).
    * Pressione `#` para salvar a nova senha e sair do modo de programação.
    * Pressione `*` para apagar o último dígito da nova senha.
4.  **Controle Remoto (MQTT):**
    * Envie uma mensagem JSON para o tópico `fechadura/comandos`.
    * Payload: `{"command":"toggle"}`
    * Isso fará a fechadura alternar seu estado. O novo estado será publicado no tópico `fechadura/estado`.