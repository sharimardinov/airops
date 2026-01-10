import http from "k6/http";
import { check, sleep, fail } from "k6";

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const V1 = `${BASE}/api/v1`;

export const options = {
    vus: 1,
    duration: __ENV.DURATION || "10s",
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<800"],
    },
};

function must200(res, name) {
    check(res, { [`${name} 200`]: (r) => r.status === 200 }) || fail(`${name} status=${res.status} body=${res.body}`);
}

export function setup() {
    // Берём flightId + airplaneCode из реальных данных
    const fList = http.get(`${V1}/flights?limit=1&offset=0`, { tags: { name: "setup_flights_list" } });
    must200(fList, "setup flights list");
    const arr = fList.json();
    const flightId = arr?.[0]?.flight_id ?? arr?.[0]?.id ?? arr?.[0]?.ID;
    if (!flightId) fail(`setup: can't pick flight id: ${fList.body}`);

    const aList = http.get(`${V1}/airplanes?limit=1`, { tags: { name: "setup_airplanes_list" } });
    must200(aList, "setup airplanes list");
    const aArr = aList.json();
    const airplane = aArr?.[0]?.Code ?? aArr?.[0]?.code;
    if (!airplane) fail(`setup: can't pick airplane code: ${aList.body}`);

    // Даты обязательны для /stats/routes — берём фикс “последние 7 дней” от текущей даты
    const to = new Date();
    const from = new Date(to.getTime() - 7 * 24 * 60 * 60 * 1000);
    const fromS = from.toISOString().slice(0, 10);
    const toS = to.toISOString().slice(0, 10);

    return { flightId, airplane, fromS, toS };
}

export default function (data) {
    must200(http.get(`${BASE}/health`, { tags: { name: "health" } }), "health");
    must200(http.get(`${BASE}/ready`, { tags: { name: "ready" } }), "ready");
    must200(http.get(`${BASE}/debug/pool`, { tags: { name: "pool" } }), "pool");

    must200(http.get(`${V1}/flights?limit=10&offset=0`, { tags: { name: "flights_list" } }), "flights_list");
    must200(http.get(`${V1}/flights/${data.flightId}`, { tags: { name: "flight_get" } }), "flight_get");
    must200(http.get(`${V1}/flights/${data.flightId}/passengers?limit=10&offset=0`, { tags: { name: "passengers" } }), "passengers");

    must200(http.get(`${V1}/airplanes?limit=10`, { tags: { name: "airplanes_list" } }), "airplanes_list");
    must200(http.get(`${V1}/airplanes/${data.airplane}`, { tags: { name: "airplane_get" } }), "airplane_get");

    must200(
        http.get(`${V1}/stats/routes?limit=10&from=${data.fromS}&to=${data.toS}`, { tags: { name: "stats_routes" } }),
        "stats_routes"
    );

    sleep(0.1);
}