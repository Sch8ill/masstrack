var map = L.map("map").setView([51.1642292, 10.4541194], 8);
L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
displayAllDevices();

var device = urlParams.get("device");
if (device !== null) {
    displayDevicePath(device);
}


async function displayAllDevices() {
    const locations = await fetchLocations();
    for (let i = 0; i < locations.length; i++) {
        var l = locations[i];
        var marker = L.marker([l.latitude, l.longitude]).addTo(map);
        if (l.device.length === 40) {
            var device = "android"
        } else {
            var device = "iphone"
        }
        marker.bindPopup(
            l.timestamp +
            "</br>" +
            device +
            "</br>" +
            "<a href='/?device=" + l.device + "'>track</a>"
        )
    }
    document.getElementById("pointsCount").textContent = locations.length;
}

async function displayDevicePath(device) {
    const locations = await fetchDeviceLocations(device);
    for (let i = 0; i < locations.length; i++) {
        var l = locations[i];
        var marker = L.marker([l.latitude, l.longitude]).addTo(map);
        marker.bindPopup(l.timestamp);

        if (i == 0) {
            map.setView([l.latitude, l.longitude], 12);
        } else {
            L.polyline([[l.latitude, l.longitude], [locations[i - 1].latitude, locations[i - 1].longitude]], { color: "rgb(255,0,0)" }).addTo(map);
        }
    }
}

async function fetchDeviceLocations(device) {
    const res = await fetch("/api/v1/locations?device=" + device);
    return await res.json();
}

async function fetchLocations() {
    const res = await fetch("/api/v1/locations");
    return await res.json();
}