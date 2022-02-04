# ReferenceDigitalTwin
A networking framework for connecting nodes running digital twins. Each node has an MQTT client for IoT messaging and a MongoDB database.

Requirements: 
- Eclipse Mosquitto v5, installed and accessible on path.
- MongoDB v5, installed and accessible on path.

Run with:

$ go build -o /bin/dt.exe

$ ./bin/dt.exe

Todo:
* Unsure how sendMessage is (intermittently?) connecting to MQTTClient
* Is receiveMessage connecting to SubMessages?
* Remove the channels, split timed and event driven behaviour
	* Timed behaviour can request data from DB if needed
	* Event driven can be triggered by MQTT messages, just need to be smarter about where receiver funcs go.

	* MQTT Transport layer security https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang
	* Weather api for glasgow?
	*	 Pollution calculation? --Geographic layout and wind direction
