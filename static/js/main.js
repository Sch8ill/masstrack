var map = L.map("map").setView([51.1642292, 10.4541194], 8);
L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
const device = urlParams.get("device");

const heatmap = urlParams.get("heatmap");
if (heatmap === "true") {
    displayHeatmap();
    document.getElementById("heatmap").checked = true;
} else {
    displayAllDevices();
    if (device !== null) {
        displayDevicePath(device);
    }
}

async function displayAllDevices() {
    const locations = await fetchLocations();
    for (let i = 0; i < locations.length; i++) {
        var l = locations[i];
        var marker = L.marker([l.latitude, l.longitude]).addTo(map);
        marker.bindPopup(
            l.timestamp +
            "</br>" +
            deviceType(l.device) +
            "</br>" +
            "<a href='/?device=" + l.device + "'>track</a>"
        )
    }
    displayStats(locations.length);
}


async function displayHeatmap() {
    const locations = await fetchLocations();
    var coordinates = new Array()
    for (let i = 0; i < locations.length; i++) {
        var l = locations[i];
        coordinates.push([l.latitude, l.longitude, 200]) // + intensity
    }
    console.log(coordinates);
    L.heatLayer(coordinates, { minOpacity: 0.00001, radius: 25 }).addTo(map);
    displayStats(locations.length);
}

async function displayDevicePath(device) {
    const locations = await fetchDeviceLocations(device);
    for (let i = 0; i < locations.length; i++) {
        var l = locations[i];

        L.marker([l.latitude, l.longitude]).addTo(map).bindPopup(l.timestamp);

        // line to previous location
        if (i != 0) {
            L.polyline([[l.latitude, l.longitude], [locations[i - 1].latitude, locations[i - 1].longitude]], { color: "rgb(255,0,0)" }).addTo(map);
        }

        // set view point to last location
        if (i == locations.length - 1) {
            map.setView([l.latitude, l.longitude], 13);
        }
    }
    document.getElementById("deviceInfo").innerHTML = "<h4>ID: </h4>" + device + "<h4>Type: </h4>" + deviceType(device) + "<h4>First seen: </h4>" + locations[0].timestamp + "<h4>Last seen: </h4>" + locations[locations.length - 1].timestamp;

    displayStats(locations.length);
}

async function displayStats(locations) {
    document.getElementById("pointsCount").textContent = locations;
    console.log("displaying: " + locations);
}

async function fetchDeviceLocations(device) {
    const res = await fetch("/api/v1/locations?device=" + device);
    return await res.json();
}

async function fetchLocations() {
    const res = await fetch("/api/v1/locations" + window.location.search);
    return await res.json();
}

function deviceType(device) {
    var type = "iphone";
    if (device.length === 40) {
        type = "android"
    }
    return type;
}

function go(e) {
    e.preventDefault();
    const start = document.getElementById("start").value;
    const end = document.getElementById("end").value;
    const heatmap = document.getElementById("heatmap").checked;

    var url = "/?heatmap=" + encodeURIComponent(heatmap) + "&start=" + encodeURIComponent(start) + "&end=" + encodeURIComponent(end);
    if (device !== null) {
        url += "&device=" + device;
    }

    window.location.href = url;
}

function reset() {
    window.location.href = "/";
}