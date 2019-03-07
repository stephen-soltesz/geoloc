// Create and add the canvas.

function drawWM(data) {
    console.log(JSON.stringify(data, null, 2));
}

// subscribe to data and style changes.
dscc.subscribeToData(drawWM, {transform: dscc.objectTransform});
