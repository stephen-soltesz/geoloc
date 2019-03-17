var div = document.createElement('div');
div.id = "map-holder";
div.style.width = "100vw";
div.style.height = "99vh";
document.body.appendChild(div);
document.body.style.backgroundColor = "#2A2C39";

function randInt(max) {
    return parseInt((Math.random() * max).toFixed()).toString(16);
}

function randColor(alpha) {
    var c = "#000000";
    var tc = tinycolor(c);
    // Loop until we get a light enough color.
    while (tc.isDark()) {
        c = "#" + randInt(0x100) + randInt(0x100) + randInt(0x100) + alpha;
        tc = tinycolor(c);
    }
    return c
}

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
    locations.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
    if (points != '') {
        points.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
    }
    console.log([t.x, t.y, t.k]);
}

// Define map zoom behaviour
var zoom = d3
    .zoom()
    .on("zoom", zoomed);

function getTextBox(selection) {
    selection
        .each(function(d) {
            d.bbox = this
                .getBBox();
        });
}

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

// var metroMap = {};
// var metros = {};

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

    var m = [];
    var ms = {};
    sites.forEach(function(s) {
        if (!(s.metro in ms)) {
            ms[s.metro] = true;
            m.push([s.lon, s.lat, s.metro]);
        }
    });
    locations = svg.selectAll("circle")
        .data(m).enter()
        .append("circle")
        .attr("cx", function(d) {
            return projection(d)[0];
        })
        .attr("cy", function(d) {
            return projection(d)[1];
        })
        .attr('r', "10px")
        .style("fill", "#ffffff11")
        .on("click", function(d, i) {
            console.log(d, i);
            name = d[2];
            if (typeof metroMap === 'undefined') {
                console.log('IGNORING undefined metroMap');
                return
            }
            var remove = metroMap[name];
            metroMap[name] = !remove;
            if (remove) {
                var even = svg.selectAll("rect").filter(
                    function(d, i) {
                        return d[3] == name;
                    });
                even.remove();
            } else {
                if (metros[name] === undefined) {
                    return
                }
                // TODO: is there a more efficient way to do this?
                data = svg.selectAll("rect").data();
                points = svg.selectAll("rect")
                    .data(data.concat(metros[name]))
                    .join("rect")
                    .attr("x", function(e) {
                        return projection(e)[0];
                    })
                    .attr("y", function(e) {
                        return projection(e)[1];
                    })
                    .attr('width', ".4px")
                    .attr('height', ".4px")
                    .style("fill", function(e) {
                        return e[2];
                    });
            }
            points.attr("transform", "translate(" + [t.x, t.y] + ")scale(" + t.k + ")");
        });

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
    // Add a label group to each feature/country. This will contain the country name and a background rectangle
    // Use CSS to have class "countryLabel" initially hidden
    countryLabels = countriesGroup
        .selectAll("g")
        .data(json.features)
        .enter()
        .append("g")
        .attr("class", "countryLabel")
        .attr("id", function(d) {
            return "countryLabel" + d.properties.iso_a3;
        })
        .attr("transform", function(d) {
            return (
                "translate(" + path.centroid(d)[0] + "," + path.centroid(d)[1] + ")"
            );
        })
        // add mouseover functionality to the label
        .on("mouseover", function(d, i) {
            d3.select(this).style("display", "block");
        })
        .on("mouseout", function(d, i) {
            d3.select(this).style("display", "none");
        })
        // add an onlcick action to zoom into clicked country
        .on("click", function(d, i) {
            d3.selectAll(".country").classed("country-on", false);
            d3.select("#country" + d.properties.iso_a3).classed("country-on", true);
            boxZoom(path.bounds(d), path.centroid(d), 20);
        });
    initiateZoom();
};

var getPoints = function() {
    var p = [];
    var t = 0;
    var s = 0;
    $.each(metroMap, function(name, show) {
        t += 1;
        if (show) {
            s += 1;
            for (i=0; i < metros[name].length; i++ ) {
                p.push(metros[name][i]);
            }
        }
    });
    console.log('Showing:', s, ' of ', t);
    return p;
}

function loadData(data) {
    var i = 0;
    metroMap = {};
    metros = {};
    var colorMap = {};
    // var icolors = palette("rainbow", 16);
    // var icolors = palette("mpn65", 12);
    /*
    var icolors = [
        'rgb(255, 0, 41)',
        'rgb(0, 210, 213)',
        'rgb(255, 127, 0)',
        'rgb(179, 233, 0)',
        'rgb(247, 129, 191)',
        'rgb(251, 128, 114)',
        'rgb(252, 205, 229)',
        'rgb(255, 237, 111)',
        'rgb(0, 208, 103)',
    ];
    console.log(icolors);
    */

    var ic = 0;
    var p = [];
    data.tables.DEFAULT.forEach(function(c) {
        if (true) { // i % 3 == 0) {
            if (c.hasOwnProperty('geoDimension') && c.hasOwnProperty('metroDimension')) {
                var ll = c.geoDimension[0].split(',')
                var v = {
                    metro: c.metroDimension[0],
                    lat: parseFloat(ll[0]),
                    lon: parseFloat(ll[1]),
                };
                if (!(v.metro in metros)) {
                    metros[v.metro] = [];
                    color = data.theme.themeSeriesColor[metroColor[v.metro]].color;
                    colorMap[v.metro] = color;
                    console.log(v.metro + " " + color);
                    ic += 1;
                }
                if (true) {
                    metros[v.metro].push([v.lon, v.lat, colorMap[v.metro], v.metro]);
                    metroMap[v.metro] = true;
                }
            }
        }
        i += 1;
    });

    p = getPoints();
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
