# locationengine

This Go library helps you to connect to Kontakt.io's MQTT broker and receive 
realtime location data of items near your Kontakt.io receiers.

After parsing the location data, the library triggers your callback functions,
notifying your application of new items appearing, items disappearing, as well
as changes in RSSI and proximity status.

More information on the API used can be found on the 
[Kontakt.io Location Engine Monitoring page](https://developer.kontakt.io/rest-api/api-guides/location-engine-monitoring/#mqtt)
