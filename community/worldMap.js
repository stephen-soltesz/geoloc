// Create and add the canvas.

function drawWM(data) {
    console.log(JSON.stringify(data));
}

// subscribe to data and style changes.
dscc.subscribeToData(drawWM, {transform: dscc.objectTransform});
