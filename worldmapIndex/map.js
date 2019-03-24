var div = document.createElement('div');
div.id = "map-holder";
div.style.width = "100vw";
div.style.height = "99vh";
document.body.appendChild(div);
document.body.style.backgroundColor = "#2A2C39";

// DEFINE VARIABLES
// Define size of map group
// Full world map is 2:1 ratio
// Using 12:5 because we will crop top and bottom of map
w = dscc.getWidth()-5;
h = dscc.getHeight()-15;

// variables for catching min and max zoom factors
var minZoom;
var maxZoom;

// DEFINE FUNCTIONS/OBJECTS
// Define map projection
var projection = d3
    .geoEquirectangular()
    .center([0, 15]) // set centre to further North as we are cropping more off bottom of map
    .scale([w / (2 * Math.PI)]) // scale to fit group width
    .translate([w / 2, h / 2]) // ensure centred in group
;

// Define map path
var path = d3
    .geoPath()
    .projection(projection);

var points = '';

// Create function to apply zoom to countriesGroup
function zoomed() {
    t = d3.event.transform;
    countriesGroup.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
    // locations.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
    if (points != '') {
        points.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
    }
    console.log([t.x, t.y, t.k]);
}

// Define map zoom behaviour
var zoom = d3
    .zoom()
    .on("zoom", zoomed);

// Function that calculates zoom/pan limits and sets zoom to default value
function initiateZoom() {
    // Define a "minzoom" whereby the "Countries" is as small possible without leaving white space at top/bottom or sides
    minZoom = Math.max($("#map-holder").width() / w, $("#map-holder").height() / h);
    // set max zoom to a suitable factor of this value
    maxZoom = 20 * minZoom;
    // set extent of zoom to chosen values
    // set translate extent so that panning can't cause map to move out of viewport
    zoom
        .scaleExtent([minZoom, maxZoom])
        .translateExtent([
            [0, 0],
            [w, h]
        ]);
    // define X and Y offset for centre of map to be shown in centre of holder
    midX = ($("#map-holder").width() - 1.5 * w) / 4;
    midY = ($("#map-holder").height() - 1.5 * h) / 4;
    // change zoom transform to min zoom and centre offsets
    svg.call(zoom.transform, d3.zoomIdentity.translate(midX, midY).scale(1.5));
}

// zoom to show a bounding box, with optional additional padding as percentage of box size
function boxZoom(box, centroid, paddingPerc) {
    minXY = box[0];
    maxXY = box[1];
    // find size of map area defined
    zoomWidth = Math.abs(minXY[0] - maxXY[0]);
    zoomHeight = Math.abs(minXY[1] - maxXY[1]);
    // find midpoint of map area defined
    zoomMidX = centroid[0];
    zoomMidY = centroid[1];
    // increase map area to include padding
    zoomWidth = zoomWidth * (1 + paddingPerc / 100);
    zoomHeight = zoomHeight * (1 + paddingPerc / 100);
    // find scale required for area to fill svg
    maxXscale = $("svg").width() / zoomWidth;
    maxYscale = $("svg").height() / zoomHeight;
    zoomScale = Math.min(maxXscale, maxYscale);
    // handle some edge cases
    // limit to max zoom (handles tiny countries)
    zoomScale = Math.min(zoomScale, maxZoom);
    // limit to min zoom (handles large countries and countries that span the date line)
    zoomScale = Math.max(zoomScale, minZoom);
    // Find screen pixel equivalent once scaled
    offsetX = zoomScale * zoomMidX;
    offsetY = zoomScale * zoomMidY;
    // Find offset to centre, making sure no gap at left or top of holder
    dleft = Math.min(0, $("svg").width() / 2 - offsetX);
    dtop = Math.min(0, $("svg").height() / 2 - offsetY);
    // Make sure no gap at bottom or right of holder
    dleft = Math.max($("svg").width() - w * zoomScale, dleft);
    dtop = Math.max($("svg").height() - h * zoomScale, dtop);
    // set zoom
    svg
        .transition()
        .duration(500)
        .call(
            zoom.transform,
            d3.zoomIdentity.translate(dleft, dtop).scale(zoomScale)
        );
}

// on window resize
$(window).resize(function() {
    // Resize SVG
    svg
        .attr("width", $("#map-holder").width())
        .attr("height", $("#map-holder").height());
    initiateZoom();
});

// create an SVG
var svg = d3
    .select("#map-holder")
    .append("svg")
    // set to the same size as the "map-holder" div
    .attr("width", $("#map-holder").width())
    .attr("height", $("#map-holder").height())
    // add zoom functionality
    .call(zoom);

// get map data
var loadMap = function(json) {

    //Bind data and create one path per GeoJSON feature
    countriesGroup = svg.append("g").attr("id", "map");
    // add a background rectangle
    countriesGroup
        .append("rect")
        .attr("x", 0)
        .attr("y", 0)
        .attr("width", w)
        .attr("height", h)
        .style("fill", "#2A2C39");

    // draw a path for each feature/country
    countries = countriesGroup
        .selectAll("path")
        .data(json.features)
        .enter()
        .append("path")
        .attr("d", path)
        .attr("id", function(d, i) {
            return "country" + d.properties.iso_a3;
        })
        .style("fill", "#404040") /* country colour */
        .style("stroke", "#2A2C39") /* country border colour */
        .style("stroke-width", 1) /* country border width */
        .attr("class", "country")
        // add a mouseover action to show name label for feature/country
        .on("mouseover", function(d, i) {
            d3.select("#countryLabel" + d.properties.iso_a3).style("display", "block");
        })
        .on("mouseout", function(d, i) {
            d3.select("#countryLabel" + d.properties.iso_a3).style("display", "none");
        })
        // add an onclick action to zoom into clicked country
        .on("click", function(d, i) {
            d3.selectAll(".country").classed("country-on", false);
            d3.select(this).classed("country-on", true);
            boxZoom(path.bounds(d), path.centroid(d), 20);
        });

    initiateZoom();
};

function loadData(data) {
    var i = 0;

    var p = [];
    data.tables.DEFAULT.forEach(function(c) {
        if (c.hasOwnProperty('geoDimension') && c.hasOwnProperty('valueMetric')) {
            var ll = c.geoDimension[0].split(',');
            var v = {
                lat: parseFloat(ll[0]),
                lon: parseFloat(ll[1]),
            };
            var index = c.valueMetric[0];
            color = data.theme.themeSeriesColor[index % 20].color;
            p.push([v.lon, v.lat, color]);
        }
        i += 1;
    });

    points = svg.selectAll("rect")
        .data(p).join("rect")
        .attr("x", function(d) {
            return projection(d)[0];
        })
        .attr("y", function(d) {
            return projection(d)[1];
        })
        .attr('width', ".4px")
        .attr('height', ".4px")
        .style("fill", function(d) {
            return d[2];
        });
    points.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
}

loadMap(wm);
console.log("load data");
dscc.subscribeToData(loadData, {
    transform: dscc.objectTransform
});
