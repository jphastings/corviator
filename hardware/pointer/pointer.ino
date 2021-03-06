#include <WiFiManager.h>
#include <PubSubClient.h>
#include <AccelStepper.h>
#include <ArduinoJson.h>
#include <math.h>
#include "vars.h";

#define APP_NAME "jan-poka:pointer"
#define ONE_ROTATION_STEPS 8192

int INNER_DIR = 14; // D5
int INNER_STP = 12; // D6
int OUTER_DIR = 5;  // D1
int OUTER_STP = 4;  // D2
int DISABLE   = 16; // D0
int BOOTING   = 2;  // D4

int MAX_SPEED = 9000;

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

AccelStepper inner = AccelStepper(AccelStepper::DRIVER, INNER_STP, INNER_DIR);
AccelStepper outer = AccelStepper(AccelStepper::DRIVER, OUTER_STP, OUTER_DIR);

void setup() {
  Serial.begin(115200);
  pinMode(BOOTING, OUTPUT);
  digitalWrite(BOOTING, HIGH);
  pinMode(DISABLE, OUTPUT);
  digitalWrite(DISABLE, HIGH);

  mqttClient.setServer(MQTT_HOST, MQTT_PORT);

  inner.setMaxSpeed(MAX_SPEED);
  outer.setMaxSpeed(MAX_SPEED);

  mqttConnectionLoop();
  Serial.println("\nBooted");
  digitalWrite(BOOTING, LOW);
}

void mqttConnectionLoop() {
  if (!mqttClient.connected()) {
    if (!mqttClient.connect(APP_NAME, MQTT_USER, MQTT_PASS)) {
      Serial.println("Failed to connect to MQTT broker");
      // TODO: Backoff
      return;
    }
    Serial.println("MQTT connected");
    mqttClient.setCallback(handleGeoTarget);
    mqttClient.subscribe(MQTT_TOPIC);
  }
  mqttClient.loop();
}

void handleGeoTarget(char* topic, byte* payload, unsigned int length) {
  /* Copied from EspMQTTClient: https://github.com/plapointe6/EspMQTTClient/blob/master/src/EspMQTTClient.cpp#L649 */
  // Convert the payload into a String
  // First, We ensure that we dont bypass the maximum size of the PubSubClient library buffer that originated the payload
  // This buffer has a maximum length of _mqttClient.getBufferSize() and the payload begin at "headerSize + topicLength + 1"
  unsigned int strTerminationPos;
  if (strlen(topic) + length + 9 >= mqttClient.getBufferSize()) {
    strTerminationPos = length - 1;
  } else {
    strTerminationPos = length;
  }
  
  // Second, we add the string termination code at the end of the payload and we convert it to a String object
  payload[strTerminationPos] = '\0';
  String payloadStr((char*)payload);
  /* end */

  StaticJsonDocument<512> jsonDoc;
  DeserializationError err = deserializeJson(jsonDoc, payload);
  if (err != DeserializationError::Ok) {
    Serial.print("Unsuccessful at parsing MQTT JSON: ");
    Serial.println(err.f_str());
    return;
  }
  goTo(jsonDoc["azi"], jsonDoc["ele"]);
}

void goTo(double azimuth, double elevation) {
  double zRotate = azimuth / 180.0;
  double xyRotate = elevation / 270.0;
  long innerRotation = floor(zRotate * ONE_ROTATION_STEPS);
  long outerRotation = floor((xyRotate - 2 * zRotate) * ONE_ROTATION_STEPS);

  long stepsForInner = innerRotation - inner.currentPosition();
  long stepsForOuter = outerRotation - outer.currentPosition();

  inner.moveTo(innerRotation);
  outer.moveTo(outerRotation);

  double innerSpeed = MAX_SPEED;
  double outerSpeed = MAX_SPEED;
  if (stepsForInner > stepsForOuter) {
    outerSpeed = innerSpeed * stepsForInner / stepsForOuter;
  } else {
    innerSpeed = outerSpeed * stepsForOuter / stepsForInner;
  }

  inner.setSpeed(innerSpeed);
  outer.setSpeed(outerSpeed);
}

void loop() {
  bool noMove = inner.distanceToGo() == 0 && outer.distanceToGo() == 0;
  digitalWrite(DISABLE, noMove ? HIGH : LOW);
  
  inner.runSpeedToPosition();
  outer.runSpeedToPosition();
  mqttConnectionLoop();
}
