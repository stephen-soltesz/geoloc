local clients = import 'clients.json';

local values = [
  {
    metroDimension: [client.metro],
    geoDimension: [client.lat + ',' + client.lon],
    valuemetric: [1],
  }
  for client in clients
];

local data = {
  tables: {
    DEFAULT: values
  },
  fields: {
    metroDimension: [
      {
        id: 'qt_5dwcxhixvb',
        name: 'metro',
        type: 'TEXT',
        concept: 'DIMENSION'
      },
    ],
    geoDimension: [
      {
        id: 'qt_zgbqe2ixvb',
        name: 'latlon',
        type: 'NUMBER',
        concept: 'DIMENSION'
      }
    ],
    valueMetric: [
      {
        id: 'qt_zgbqe2abcd',
        name: 'value',
        type: 'NUMBER',
        concept: 'METRIC'
      }
    ]
  },
  style: {
    barColor: {
      value: {
        color: '#F44336',
        opacity: 1
      },
      defaultValue: '#000000'
    }
  },
  theme: {
    themeSeriesColor: [
      { color: "#FF1744" }, // red
      { color: "#d500f9" }, // pink
      { color: "#798dff" }, // purple
      { color: "#40c4ff" }, // blue
      { color: "#00e676" }, // green
      { color: "#eeff41" }, // yellow
      { color: "#ffab00" }, // orange
      { color: "#efebe9" }, // white
      { color: "#ffffff" } // white
    ],
  },
};

'data =' + std.toString(data) + ';'
