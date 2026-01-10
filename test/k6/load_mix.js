import http from "k6/http";
import { check, fail, sleep } from "k6";

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const V1 = `${BASE}/api/v1`;

export const options = {
    scenarios: {
        mix: {
            executor: "constant-arrival-rate",
            rate: Number(__ENV.RPS || 200),
            timeUnit: "1s",
            duration: __ENV.DURATION || "60s",
            preAllocatedVUs: Number(__ENV.PRE_VUS || 200),
            maxVUs: Number(__ENV.MAX_VUS || 2000),
        },
        stats: {
            executor: "constant-arrival-rate",
            rate: Number(__ENV.RPS_STATS || 20),
            timeUnit: "1s",
            duration: __ENV.DURATION || "60s",
            preAllocatedVUs: Number(__ENV.PRE_VUS_STATS || 50),
            maxVUs: Number(__ENV.MAX_VUS_STATS || 500),
            exec: "statsOnly",
            startTime: "0s",
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        dropped_iterations: ["count<10"],

        "http_req_duration{name:flights_list}": ["p(95)<300"],
        "http_req_duration{name:flight_get}": ["p(95)<300"],
        "http_req_duration{name:passengers}": ["p(95)<500"],
        "http_req_duration{name:airplanes_list}": ["p(95)<250"],
        "http_req_duration{name:airplane_get}": ["p(95)<250"],
        "http_req_duration{name:stats_routes}": ["p(95)<800"],   // подстрой под реальность
    },
};

function must200(res, name) {
    check(res, { [`${name} 200`]: (r) => r.status === 200 }) || fail(`${name} status=${res.status} body=${res.body}`);
}

export function setup() {
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

    const to = new Date();
    const from = new Date(to.getTime() - 7 * 24 * 60 * 60 * 1000);
    const fromS = from.toISOString().slice(0, 10);
    const toS = to.toISOString().slice(0, 10);

    return { flightId, airplane, fromS, toS };
}

function get(url, name) {
    const res = http.get(url, { tags: { name } });
    check(res, { "200": (r) => r.status === 200 });
    return res;
}

export default function (data) {
    // веса — без фанатизма
    const x = Math.random();

    if (x < 0.35) {
        get(`${V1}/flights?limit=10&offset=0`, "flights_list");
    } else if (x < 0.55) {
        get(`${V1}/flights/${data.flightId}`, "flight_get");
    } else if (x < 0.75) {
        get(`${V1}/flights/${data.flightId}/passengers?limit=50&offset=0`, "passengers");
    } else if (x < 0.90) {
        get(`${V1}/airplanes?limit=10`, "airplanes_list");
    } else {
        get(`${V1}/airplanes/${data.airplane}`, "airplane_get");
    }

    sleep(0.01);
}

export function statsOnly(data) {
    // отдельная тяжёлая ручка — всегда с датами (ты её обязал)
    get(`${V1}/stats/routes?limit=10&from=${data.fromS}&to=${data.toS}`, "stats_routes");
    sleep(0.01);
}