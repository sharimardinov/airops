import http from "k6/http";
import { check, fail } from "k6";

export const options = {
    scenarios: {
        hard5k: {
            executor: "ramping-arrival-rate",
            timeUnit: "1s",
            startRate: 500,
            stages: [
                { target: 1000, duration: "20s" },
                { target: 2000, duration: "20s" },
                { target: 3000, duration: "20s" },
                { target: 4000, duration: "20s" },
                { target: 5000, duration: "30s" },
                { target: 5000, duration: "60s" },
                { target: 0,    duration: "10s" },
            ],
            preAllocatedVUs: 3000,
            maxVUs: 20000,
            exec: "mixOneReq",
            gracefulStop: "30s",
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.05"],
        http_req_duration: ["p(95)<1200"],
        dropped_iterations: ["count==0"],
    },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";

function pickId(body) {
    if (!Array.isArray(body) || body.length === 0) return null;
    const f = body[0];
    return f?.id ?? f?.ID ?? null;
}
function toYYYYMMDD(ts) {
    if (!ts || typeof ts !== "string") return null;
    return ts.slice(0, 10);
}

export function setup() {
    const listRes = http.get(`${BASE}/api/v1/flights?limit=1&offset=0`);
    if (listRes.status !== 200) fail(`setup flights list ${listRes.status}: ${listRes.body}`);
    let listJson;
    try { listJson = listRes.json(); } catch { fail(`setup flights list not json: ${listRes.body}`); }

    const flightId = pickId(listJson);
    if (!flightId) fail("setup: can't pick flight id");

    const flightRes = http.get(`${BASE}/api/v1/flights/${flightId}`);
    if (flightRes.status !== 200) fail(`setup flight by id ${flightRes.status}: ${flightRes.body}`);

    let f;
    try { f = flightRes.json(); } catch { fail(`setup flight by id not json: ${flightRes.body}`); }

    const from = f.departure_airport ?? f.DepartureAirport ?? f.departureAirport ?? null;
    const to   = f.arrival_airport   ?? f.ArrivalAirport   ?? f.arrivalAirport   ?? null;
    const date = toYYYYMMDD(f.scheduled_departure ?? f.ScheduledDeparture ?? f.scheduledDeparture ?? null);
    if (!from || !to || !date) fail("setup: missing from/to/date");

    return { flightId, from, to, date };
}

export function mixOneReq(data) {
    const x = Math.random();
    let url;

    // 50% passengers (дешево)
    if (x < 0.50) {
        url = `${BASE}/api/v1/flights/${data.flightId}/passengers?limit=200&offset=0`;

        // 25% flights list
    } else if (x < 0.75) {
        url = `${BASE}/api/v1/flights?limit=200&offset=0`;

        // 20% search (тяжелее)
    } else if (x < 0.95) {
        url =
            `${BASE}/api/v1/flights/search` +
            `?from=${encodeURIComponent(data.from)}` +
            `&to=${encodeURIComponent(data.to)}` +
            `&date=${encodeURIComponent(data.date)}` +
            `&passengers=1&fare_class=Economy` +
            `&limit=50&offset=0`;

        // 5% stats
    } else {
        url = `${BASE}/api/v1/stats/routes?from=${data.date}&to=${data.date}&limit=20`;
    }

    const res = http.get(url);
    check(res, { "status 200": (r) => r.status === 200 });
}