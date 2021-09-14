package SenseHatIoT

/*
 * templates.go
 * HTML templates to render on the frontend!
 *
 * Principal author(s) : Abhishek Tiwary
 */

type HtmlTemplates struct{}

func (*HtmlTemplates) IndexPageTemplate() []byte {
	templ := []byte(`<!DOCTYPE html>
<html>
<head>
<title>How Hot Is Ab's Desk</title>

<link rel="stylesheet" type="text/css" href="/static/css/bootstrap.min.css">
<link rel="stylesheet" type="text/css" href="static/css/bootstrap-theme.min.css">

<script src="/static/js/moment.min.js"></script>
<script src="/static/js/Chart.js"></script>
<script src="/static/js/jquery-2.2.3.min.js"></script>
<script src="/static/js/bootstrap.min.js"></script>

<script language="javascript" type="text/javascript">
var wsuri = null;
var wsocket = null;

var data = [ {{ with .measurements }}{{ range . }} { "timestamp":{{ .Timestamp }}, "temperature":{{ .Temperature }} }, {{ end }}{{ end }} ];

</script>

</head>

<body>
<div style="width:75%;">
	<canvas id="canvas"></canvas>
</div>

<br/><br/>

<div id="value" style="position: absolute; bottom: 0; left: 0;"></div>

<script type="text/javascript">
window.chartColors = {
	red: 'rgb(255, 99, 132)',
	orange: 'rgb(255, 159, 64)',
	yellow: 'rgb(255, 205, 86)',
	green: 'rgb(75, 192, 192)',
	blue: 'rgb(54, 162, 235)',
	purple: 'rgb(153, 102, 255)',
	grey: 'rgb(201, 203, 207)'
};

var timeFormat = 'MM/DD/YYYY HH:mm:ss';

function PlotData() {
    var labels = [];
    for(var i=0; i < data.length; ++i)
        labels.push(new Date(data[i]["timestamp"]));

    var plota = [];
    for(var i=0; i < data.length; ++i)
        plota.push(data[i]["temperature"]);

    var color = Chart.helpers.color;
	var config = {
	    type: 'line',
		data: {
		    labels: labels,
			datasets: [{
				label: "Temperature",
				borderColor: window.chartColors.red,
				fill: false,
				data: plota,
				},
				]
		},
		options: {
            title:{
                text: "Chart.js Time Scale"
            },
		    scales: {
			    xAxes: [{
					type: "time",
					time: {
						format: timeFormat,
						// round: 'minute'
						tooltipFormat: 'll HH:mm:ss'
					},
					scaleLabel: {
						display: true,
						labelString: 'Date'
					}
				}, ],
				yAxes: [{
					scaleLabel: {
						display: true,
						labelString: 'value'
					}
				}]
			},
		}
	};

	var ctx = document.getElementById("canvas").getContext("2d");
	window.myLine = new Chart(ctx, config);

}

window.onload = function() {
    console.log("loaded!");

	wsuri = "{{ with .serverinfo }}ws://{{ .ServerIP }}:{{ .ServerPort }}/ws{{ end }}";
	wsocket = new WebSocket(wsuri);
	
	wsocket.onmessage = function(event) {
		console.log(event.data);
		var msg = JSON.parse(event.data);
		console.log(msg);
        document.getElementById("value").innerHTML = msg;
        msg = JSON.parse(msg);
        data.push({"timestamp" : msg.timestamp, "temperature" : msg.temperature});
		PlotData();
	};
	
    PlotData();
}

</script>

</body>
</html>
	`)

	return templ
}
