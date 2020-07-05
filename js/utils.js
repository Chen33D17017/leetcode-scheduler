function calculatePoint(i, intervalSize, colorRangeInfo) {
  var { colorStart, colorEnd, useEndAsStart } = colorRangeInfo;
  return useEndAsStart
    ? colorEnd - i * intervalSize
    : colorStart + i * intervalSize;
}

function interpolateColors(dataLength, colorScale, colorRangeInfo) {
  var { colorStart, colorEnd } = colorRangeInfo;
  var colorRange = colorEnd - colorStart;
  var intervalSize = colorRange / dataLength;
  var i, colorPoint;
  var colorArray = [];

  for (i = 0; i < dataLength; i++) {
    colorPoint = calculatePoint(i, intervalSize, colorRangeInfo);
    colorArray.push(colorScale(colorPoint));
  }

  return colorArray;
}

function createChart(chartType, chartData, chartId, COLORS) {
  const chartElement = document.getElementById(chartId);
  const myChart = new Chart(chartElement, {
    type: chartType.type,
    data: {
      labels: chartData.labels,
      datasets: [
        {
          backgroundColor: COLORS,
          // hoverBackgroundColor: COLORS,
          data: chartData.data,
        },
      ],
    },
    options: {
      responsive: true,
      legend: {
        display: chartType.chartLabelDisplay,
      },
      onClick: function (e, el) {
        if (!el || el.length === 0) return;
        console.log("onClick : label " + el[0]._model.label);
      },
      onHover: function (e) {
        var point = this.getElementAtEvent(e);
        e.target.style.cursor = point.length ? "pointer" : "default";
      },
      elements: {
        center: {
          text: chartData.centerText,
          color: chartData.centerDisplayColor, // Default is #000000
          fontStyle: "Arial", // Default is Arial
          sidePadding: 20, // Defualt is 20 (as a percentage)
          display: chartData.centerDisplay,
        },
      },
      scales: {
        yAxes: [
          {
            ticks: {
              beginAtZero: true,
            },
          },
        ],
      },
    },
  });

  return myChart;
}

function dynamicColorChart(
  chartType,
  chartId,
  chartData,
  colorScale,
  colorRangeInfo = []
) {
  if (colorRangeInfo) {
    colorRangeInfo = {
      colorStart: 0,
      colorEnd: 1,
      useEndAsStart: true,
    };
  }

  const dataLength = chartData.data.length;
  /* Create color array */
  var COLORS = interpolateColors(dataLength, colorScale, colorRangeInfo);

  /* Create chart */
  return createChart(chartType, chartData, chartId, COLORS);
}

function staticColorChart(chartType, chartId, chartData, colorSchema) {
  const dataLength = chartData.data.length;
  /* Create color array */
  var COLORS = [];
  for (i = 0; i < dataLength; i++) {
    COLORS.push(colorSchema[i]);
  }
  // var COLORS = interpolateColors(dataLength, colorScale, colorRangeInfo);

  /* Create chart */
  return createChart(chartType, chartData, chartId, COLORS);
}
