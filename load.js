import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    scenarios: {
        passengers: {
            executor: "constant-arrival-rate",
            rate: 400,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 100,
            maxVUs: 500,
            exec: "passengers",
        },
        flights: {
            executor: "constant-arrival-rate",
            rate: 100,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 20,
            maxVUs: 200,
            exec: "flights",
        },
        stats: {
            executor: "constant-arrival-rate",
            rate: 10,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 5,
            maxVUs: 30,
            exec: "stats",
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.02"],
        "http_req_duration{ name:GET /flights }": ["p(95)<200"],
        "http_req_duration{ name:GET /flights/:id/passengers }": ["p(95)<300"],
        "http_req_duration{ name:GET /stats/routes }": ["p(95)<500"],
    },
};

const BASE = (__ENV.BASE_URL || "http://localhost:8080").replace(/\/$/, "");
const API = `${BASE}/api/v1`;

function randInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function passengers() {
    const id = randInt(
        parseInt(__ENV.FLIGHT_ID_MIN || "1", 10),
        parseInt(__ENV.FLIGHT_ID_MAX || "20000", 10)
    );

    const url = `${API}/flights/${id}/passengers?limit=200&offset=0`;
    const res = http.get(url, { tags: { name: "GET /flights/:id/passengers" } });

    check(res, { "passengers: 200": (r) => r.status === 200 });
    sleep(0.01);
}

export function flights() {
    const url = `${API}/flights?limit=200&offset=0`;
    const res = http.get(url, { tags: { name: "GET /flights" } });

    check(res, { "flights: 200": (r) => r.status === 200 });
    sleep(0.01);
}

export function stats() {
    const from = __ENV.STATS_FROM || "2017-01-01";
    const to = __ENV.STATS_TO || "2017-02-01";

    const url = `${API}/stats/routes?from=${from}&to=${to}&limit=20`;
    const res = http.get(url, { tags: { name: "GET /stats/routes" } });

    check(res, { "stats: 200": (r) => r.status === 200 });
    sleep(0.01);
}