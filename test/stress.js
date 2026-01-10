import http from "k6/http";
import { check, group, fail, sleep } from "k6";

export const options = {
    scenarios: {
        // основной сценарий — микс ручек
        mix: {
            executor: "ramping-arrival-rate",
            timeUnit: "1s",

            // стартуем мягко
            startRate: 20,

            stages: [
                { target: 100, duration: "30s" },  // разгон
                { target: 200, duration: "60s" },  // рабочая зона
                { target: 400, duration: "60s" },  // стресс
                { target: 600, duration: "30s" },  // перегруз (здесь уже может сыпаться)
                { target: 0, duration: "20s" },    // остывание
            ],

            preAllocatedVUs: 200,
            maxVUs: 1200,
            exec: "mix",
            gracefulStop: "30s",
        },
    },

    thresholds: {
        http_req_failed: ["rate<0.05"],      // на стрессе допускаем до 5% ошибок
        http_req_duration: ["p(95)<1200"],   // p95 до 1.2s
        "http_req_duration{group:::passengers}": ["p(95)<1500"],
        "http_req_duration{group:::flights}": ["p(95)<800"],
        "http_req_duration{group:::search}": ["p(95)<1500"],
        "http_req_duration{group:::stats}": ["p(95)<1200"],
    },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";

function pickFirstFlightId(body) {
    if (!Array.isArray(body) || body.length === 0) return null;
    const f = body[0];
    return f && (f.id ?? f.ID ?? null);
}

function toYYYYMMDD(ts) {
    if (!ts || typeof ts !== "string") return null;
    return ts.slice(0, 10);
}

export function setup() {
    // 1) flightId + from/to/date для search
    const listRes = http.get(`${BASE}/api/v1/flights?limit=1&offset=0`, {
        tags: { name: "setup_flights_list" },
    });
    if (listRes.status !== 200) fail(`setup flights list status=${listRes.status} body=${listRes.body}`);

    let listJson;
    try { listJson = listRes.json(); } catch { fail(`setup flights list not json: ${listRes.body}`); }

    const flightId = pickFirstFlightId(listJson);
    if (!flightId) fail(`setup: can't pick flight id from: ${listRes.body}`);

    const flightRes = http.get(`${BASE}/api/v1/flights/${flightId}`, {
        tags: { name: "setup_flight_by_id" },
    });
    if (flightRes.status !== 200) fail(`setup flight by id status=${flightRes.status} body=${flightRes.body}`);

    let f;
    try { f = flightRes.json(); } catch { fail(`setup flight by id not json: ${flightRes.body}`); }

    const from = f.departure_airport ?? f.DepartureAirport ?? f.departureAirport ?? null;
    const to = f.arrival_airport ?? f.ArrivalAirport ?? f.arrivalAirport ?? null;
    const depTS = f.scheduled_departure ?? f.ScheduledDeparture ?? f.scheduledDeparture ?? null;
    const date = toYYYYMMDD(depTS);

    if (!from || !to || !date) {
        fail(`setup: missing from/to/date. from=${from} to=${to} depTS=${depTS} body=${flightRes.body}`);
    }

    // 2) диапазон для stats (у тебя строго YYYY-MM-DD)
    // берём “рядом” с найденной датой: [date, date+7d]
    // (без сторонних либ — просто примитивно)
    const d0 = new Date(`${date}T00:00:00Z`);
    const d1 = new Date(d0.getTime() + 7 * 24 * 60 * 60 * 1000);
    const to2 = d1.toISOString().slice(0, 10);

    return { flightId, from, to, date, statsFrom: date, statsTo: to2 };
}

export function mix(data) {
    // простой рандом-микс по весам
    const x = Math.random();

    if (x < 0.45) {
        // 45% passengers
        group("passengers", () => {
            const res = http.get(
                `${BASE}/api/v1/flights/${data.flightId}/passengers?limit=200&offset=0`,
                { tags: { name: "passengers" } }
            );
            check(res, { "passengers 200": (r) => r.status === 200 });
        });
        return;
    }

    if (x < 0.75) {
        // 30% flights list
        group("flights", () => {
            const res = http.get(`${BASE}/api/v1/flights?limit=200&offset=0`, {
                tags: { name: "flights" },
            });
            check(res, { "flights 200": (r) => r.status === 200 });
        });
        return;
    }

    if (x < 0.93) {
        // 18% search (довольно тяжелее)
        group("search", () => {
            const url =
                `${BASE}/api/v1/flights/search` +
                `?from=${encodeURIComponent(data.from)}` +
                `&to=${encodeURIComponent(data.to)}` +
                `&date=${encodeURIComponent(data.date)}` +
                `&passengers=1` +
                `&fare_class=Economy` +
                `&limit=50&offset=0`;

            const res = http.get(url, { tags: { name: "search" } });
            check(res, { "search 200": (r) => r.status === 200 });
        });
        return;
    }

    // 7% stats
    group("stats", () => {
        const url =
            `${BASE}/api/v1/stats/routes` +
            `?from=${encodeURIComponent(data.statsFrom)}` +
            `&to=${encodeURIComponent(data.statsTo)}` +
            `&limit=20`;

        const res = http.get(url, { tags: { name: "stats" } });
        check(res, { "stats 200": (r) => r.status === 200 });
    });
}