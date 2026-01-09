import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    scenarios: {
        passengers: {
            executor: "constant-arrival-rate",
            rate: 120,          // стартуй 120 rps
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 100,
            maxVUs: 500,
            exec: "passengers",
        },
        flights: {
            executor: "constant-arrival-rate",
            rate: 20,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 20,
            maxVUs: 200,
            exec: "flights",
        },
        stats: {
            executor: "constant-arrival-rate",
            rate: 2,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 5,
            maxVUs: 30,
            exec: "stats",
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.02"],
        http_req_duration: ["p(95)<1000"], // потом ужесточишь
    },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";

function randInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function passengers() {
    const id = randInt(1, 20000); // подстрой под диапазон flight_id
    const res = http.get(`${BASE}/flights/${id}/passengers?limit=200`);
    check(res, { "passengers 200": (r) => r.status === 200 });
}

export function flights() {
    const res = http.get(`${BASE}/flights?limit=200`);
    check(res, { "flights 200": (r) => r.status === 200 });
}

export function stats() {
    // лучше всегда передавать диапазон, иначе будешь считать всю таблицу
    const from = encodeURIComponent("2017-01-01T00:00:00Z");
    const to   = encodeURIComponent("2017-02-01T00:00:00Z");
    const res = http.get(`${BASE}/stats/routes?from=${from}&to=${to}&limit=20`);
    check(res, { "stats 200": (r) => r.status === 200 });
}