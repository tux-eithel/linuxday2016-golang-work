var server = {
    socket: null,
};
server.init = function () {
    console.log("here!");
    server.socket = new WebSocket("ws://127.0.0.1:8080/stats");

    server.socket.onmessage = function (message) {
        var messageObj = JSON.parse(message.data);
        server.receive(messageObj);
    };

    server.socket.onerror = function (error) {
        console.log("some error")
    }

    server.socket.onclose = function (close) {
        console.log("by server!")
    }

};
server.receive = function (jsonObj) {

    if (jsonObj.Name == "get-hit") {

        labels = [];
        data = [];
        for (var i = 0; i < jsonObj.Data.length; i++) {
            labels.push(jsonObj.Data[i].Page)
            data.push(jsonObj.Data[i].Hit)
        }
        DataTopHit.labels = labels;
        DataTopHit.datasets[0].data = data;
        window.TopHitChart.update();

    }

    if (jsonObj.Name == "date-time") {

        DataDateRequestChart.datasets[0].data = jsonObj.Data;
        window.DateRequestChart.update();

    }

};

function MonthFromInt(value) {
    switch(value) {
        case 1:
            return "Gennaio";
        case 2:
            return "Febbraio";
        case 3:
            return "Marzo";
        case 4:
            return "Aprile";
        case 5:
            return "Maggio";
        case 6:
            return "Giugno";
        case 7:
            return "Luglio";
        case 8:
            return "Agosto";
        case 9:
            return "Settembre";
        case 10:
            return "Ottobre";
        case 11:
            return "Novembre";
        case 12:
            return "Dicembre";
    }
}

var DataTopHit = {
    labels: [],
    datasets: [{
        backgroundColor: "#FF6384",
        hoverBackgroundColor: "#FF6384",
        data: []
    }]
};
var DataDateRequestChart = {
    datasets: [{
        label: "hours-time",
        backgroundColor: "rgba(75,192,192,1)",
        hoverBackgroundColor: "rgba(75,192,192,1)",
        data: []
    }]
}
window.onload = function () {
    var ctx = document.getElementById("topHit").getContext("2d");
    window.TopHitChart = new Chart(ctx, {
        type: 'horizontalBar',
        data: DataTopHit,
        options: {
            elements: {
                rectangle: {
                    borderWidth: 2,
                    borderSkipped: 'left'
                }
            },
            title: {
                display: true,
                text: 'Top hit pages',
                fontSize: 18
            },
            legend: {
                display: false,
            }
        }
    })

    var ctx1 = document.getElementById("dateTime").getContext("2d");
    window.DateRequestChart = new Chart(ctx1, {
        type: 'bubble',
        data: DataDateRequestChart,
        options: {
            tooltips:{
                callbacks: {
                    label: function(tooltipItem, data) {
					var datasetLabel = data.datasets[tooltipItem.datasetIndex].label || '';
					var dataPoint = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index];
					return datasetLabel + ': (' + dataPoint.x + ', ' + dataPoint.y + ', ' + Math.round(dataPoint.r*1000/15) + ')';
				}
                }
            },
            title: {
                display: true,
                text: 'Visite Ore - Mesi',
                fontSize: 18
            },
            legend: {
                display: false,
            },
            scales: {                
                xAxes: [
                    {
                        scaleLabel: {
                            display: true,
                            labelString: 'Ore'
                        },
                        ticks: {
                            stepSize: 1,
                            autoSkip: false
                        }
                    }
                ],
                yAxes: [
                    {
                        scaleLabel: {
                            display: true,
                            labelString: 'Mesi'
                        },
                        ticks: {
                            stepSize: 1,
                            autoSkip: false,
                            callback: function(value) {
                                return MonthFromInt(value)
                            }
                        },
                    }
                ]
            }

        }
    })
};

server.init();