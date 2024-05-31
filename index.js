const axios = require("axios").default;
const fs = require("fs").promises;

const areaIdList = ["3608225259", "3608225592"];
const minDistance = 20;

//
// Warn if overriding existing method
if (Array.prototype.equals)
  console.warn(
    "Overriding existing Array.prototype.equals. Possible causes: New API defines the method, there's a framework conflict or you've got double inclusions in your code.",
  );
Array.prototype.equals = function (array) {
  if (!array) return false;
  if (array === this) return true;
  if (this.length != array.length) return false;
  for (var i = 0, l = this.length; i < l; i++) {
    if (this[i] instanceof Array && array[i] instanceof Array) {
      if (!this[i].equals(array[i])) return false;
    } else if (this[i] != array[i]) {
      return false;
    }
  }
  return true;
};
Object.defineProperty(Array.prototype, "equals", { enumerable: false });

function getQuery(id) {
  const query = `/*
This has been generated by the overpass-turbo wizard.
The original search was:
“((building=* or amenity=* or shop=* or office=* or place=* or public_transport=*) and name is not null and (type:node or type:way)) in Surabaya”
*/
[out:json][timeout:25]; // fetch area “Surabaya” to search in
area(id:${id})->.searchArea;
// gather results
(
  node["building"]["name"](area.searchArea);
  way["building"]["name"](area.searchArea);
  node["amenity"]["name"](area.searchArea);
  node["shop"]["name"](area.searchArea);
  node["office"]["name"](area.searchArea);
  node["place"]["name"](area.searchArea);
  way["place"]["name"](area.searchArea);
  node["public_transport"]["name"](area.searchArea);
  node["church"]["name"](area.searchArea);
  node["mosque"]["name"](area.searchArea);
  way["public_transport"]["name"](area.searchArea);
);
// print results
out center;`;

  return query;
}

const haversineDistance = ([lat1, lon1], [lat2, lon2]) => {
  const RADIUS_OF_EARTH_IN_KM = 6371;
  const toRadian = (angle) => (Math.PI / 180) * angle;
  const distance = (a, b) => (Math.PI / 180) * (a - b);

  const dLat = distance(lat2, lat1);
  const dLon = distance(lon2, lon1);

  lat1 = toRadian(lat1);
  lat2 = toRadian(lat2);

  const a =
    Math.pow(Math.sin(dLat / 2), 2) +
    Math.pow(Math.sin(dLon / 2), 2) * Math.cos(lat1) * Math.cos(lat2);
  const c = 2 * Math.asin(Math.sqrt(a));

  let finalDistance = RADIUS_OF_EARTH_IN_KM * c;

  return finalDistance * 1000;
};

const checkDistance = ([lat1, lon1], [lat2, lon2]) => {
  const distance = haversineDistance([lat1, lon1], [lat2, lon2]);

  if (distance > minDistance) {
    console.log(
      "[greater] coor1: ",
      [lat1, lon1],
      "coor2: ",
      [lat2, lon2],
      "distance: ",
      distance,
    );
    return true;
  }
  console.log(
    "[lesser] coor1: ",
    [lat1, lon1],
    "coor2: ",
    [lat2, lon2],
    "distance: ",
    distance,
  );
  return false;
};

const getPlace = async () => {
  try {
    const json = [];
    console.log("Fetching overpass api data");

    for (const id of areaIdList) {
      const { data } = await axios.post(
        "https://overpass-api.de/api/interpreter",
        {
          data: getQuery(id),
        },
        {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
        },
      );
      json.push(...data["elements"]);
      console.log(json);
    }

    await fs.writeFile(
      "result_orig.json",
      JSON.stringify(json, null, 2),
      (error) => {
        if (error) {
          console.error(error);
        }
      },
    );
    return json;
  } catch (error) {
    console.error("Error fetching data from Overpass API:", error);
  }
};

const getLatLon = (item) => {
  if (item.lat !== undefined && item.lon !== undefined) {
    return [item.lat, item.lon];
  } else if (
    item.center !== undefined &&
    item.center.lat !== undefined &&
    item.center.lon !== undefined
  ) {
    return [item.center.lat, item.center.lon];
  }
  return null;
};

const main = async () => {
  await getPlace();

  const places = JSON.parse(await fs.readFile("result_orig.json", "utf8"));

  for (let i = 0; i < places.length; i++) {
    const coord1 = getLatLon(places[i]);
    if (!coord1) return;

    for (let j = 0; j < places.length; j++) {
      const coord2 = getLatLon(places[j]);

      if (!coord2) continue;
      if (coord1.equals(coord2)) continue;

      console.log(
        "place1: ",
        places[i].tags.name,
        "place2: ",
        places[j].tags.name,
      );
      if (!checkDistance(coord1, coord2)) {
        places.splice(j, 1);
      }
    }
  }

  console.log("total places: ", places.length);

  await fs.writeFile(
    "result.json",
    JSON.stringify(places, null, 2),
    (error) => {
      if (error) {
        console.error(error);
      }
    },
  );
};

main();
