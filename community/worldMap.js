// create and add the canvas
// do this one time
/*
var canvas = document.createElement('canvas');
var ctx = canvas.getContext('2d');
canvas.id = 'myViz';
document.body.appendChild(canvas);

function drawViz(data){
  // clear the canvas
  var ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // viz code goes here

}
dscc.subscribeToData(drawViz, {transform: dscc.objectTransform})
*/

// create and add the canvas
var canvas = document.createElement('canvas');
var ctx = canvas.getContext('2d');
canvas.id = 'myViz';
document.body.appendChild(canvas);
var lastData = '';

function drawViz(data) {
  var ctx = canvas.getContext('2d');
  lastData = data;

  // clear the canvas.
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // set the canvas width and height
  ctx.canvas.width = dscc.getWidth();
  ctx.canvas.height = dscc.getHeight();

  ctx.fillStyle = data.style.barColor.value.color || data.style.barColor.defaultValue;
  ctx.fillRect(10, 10, 100, 100);

  loadData(canvas, data);
}

// subscribe to data and style changes.
dscc.subscribeToData(drawViz, {transform: dscc.objectTransform});