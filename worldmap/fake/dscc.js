dscc = {
    objectTransform: 1,
    subscribeToData: function(f, transform) {
        // TODO: handle different transform types.
        window.setTimeout(function() {
            f(data);
        }, 1000);
    },
    getWidth: function() {
        return 2400;
    },
    getHeight: function() {
        return 1000;
    },
}
