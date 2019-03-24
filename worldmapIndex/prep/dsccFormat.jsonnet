local clients = import 'clients.json';

local values = [
  {
    geoDimension: [client.latlon],
    valueMetric: [client.index],
  }
  for client in clients
];

local data = {
  tables: {
    DEFAULT: values
  },
  fields: {
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
      { color: "#d32f2f" }, // red
      { color: "#ffc107" }, // pink
      { color: "#43a047" }, // purple
      { color: "#4caf50" }, // blue
      { color: "#81c784" }, // green
      { color: "#c8e6c9" }, // yellow
      { color: "#64b5f6" }, // orange
      { color: "#2196f3" }, // white
      { color: "#ffffff" } // white
    ],
  },
};

'data =' + std.toString(data) + ';'
