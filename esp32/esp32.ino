#include <LiquidCrystal_I2C.h>

#include <Wire.h>

/*
 *  This sketch demonstrates how to scan WiFi networks.
 *  The API is based on the Arduino WiFi Shield library, but has significant changes as newer WiFi functions are supported.
 *  E.g. the return value of `encryptionType()` different because more modern encryption is supported.
 */
#include "WiFi.h"
#include <ESPAsyncWebServer.h>

#define EMPTY "-"

const char* ssid = "wifi-ssid";
const char* password = "wifi-password";

struct SavingsData {
  String total;
  String source;
};

AsyncWebServer server(80);
LiquidCrystal_I2C lcd(0x27, 16,2);
String ip = EMPTY;
bool ipPrinted = false;
SavingsData sd = SavingsData {EMPTY, EMPTY};
bool dataChanged = false;

void setup() {
  Serial.begin(115200);

  connectToWiFi();
  setupServer();

  initLCD();
}

void connectToWiFi() {
  delay(1000);

  WiFi.mode(WIFI_STA); //Optional
  WiFi.begin(ssid, password);
  Serial.println("\nConnecting");

  while(WiFi.status() != WL_CONNECTED){
      Serial.print(".");
      delay(100);
  }

  Serial.println("\nConnected to the WiFi network");
  Serial.print("Local ESP32 IP: ");
  Serial.println(WiFi.localIP().toString().c_str());

  ip = WiFi.localIP().toString();
}

void setupServer() {
  // Define a route to serve the HTML page
  server.on("/data", HTTP_POST, [](AsyncWebServerRequest* request) {
    Serial.println("ESP32 Web Server: New request received:");  // for debugging
    Serial.println("POST /data");        // for debugging
    AsyncWebParameter *total = request->getParam(0);
    AsyncWebParameter *source = request->getParam(1);
    Serial.print("DATA -> ");
    Serial.print(total->value());
    Serial.print(" : ");
    Serial.print(source->value());
    sd = SavingsData {total->value(), source->value()};
    dataChanged = true;
    request->send(200, "text/html", "<html><body><h1>Hello, ESP32!</h1></body></html>");
  });

  // Start the server
  server.begin();
}

void initLCD() {
  lcd.init();   
  lcd.backlight(); 
}

void loop() {
  if (isSavingsDataEmpty()) {
    printIpToLCD();
  } else {
    printSavingsDataToLCD();
    ipPrinted = false;
  }
  delay(1000);
}

void printIpToLCD() {
  if (ipPrinted) {
    return;
  }
  
  lcd.clear();
  lcd.setCursor(0,0);
  lcd.print(ip);

  ipPrinted = true;
}

void printSavingsDataToLCD() {
  if (!dataChanged) {
    return;
  }

  lcd.clear();
  lcd.setCursor(0,0);
  lcd.print(sd.source);
  lcd.setCursor(0,1);
  lcd.print(sd.total);

  dataChanged = false;
}

bool isSavingsDataEmpty() {
  return sd.total == EMPTY && sd.source == EMPTY;
}
