## A web server written in Go that takes data from a SenseHat, on a Raspberry Pi 3, and sends it to a websocket enabled client

I wrote this code circa 2017, while playing around with the Raspberry Pi and the SenseHat. CGo was used to interface with the SenseHat, extract temperature data, and store it in memory. Clients could access that data using a front-end powered by Chart.js. The client-side chart auto-refreshed using a websocket connection to the server.

