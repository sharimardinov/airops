import http from "k6/http";
import { check, group, sleep, fail } from "k6";

export const options = {
    scenarios: {
        passengers: {
            executor: "constant-arrival-rate",
            rate: 120,
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
        search: {
            executor: "constant-arrival-rate",
            rate: 5,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 10,
            maxVUs: 100,
            exec: "search",
        },
    },

    thresholds: {
        http_req_failed: ["rate<0.02"],
        // общая p95
        http_req_duration: ["p(95)<800"],

        // по группам (метрики k6 автоматически тегируют group)
        "http_req_duration{group:::passengers}": ["p(95)<800"],
        "http_req_duration{group:::flights}": ["p(95)<600"],
        "http_req_duration{group:::search}": ["p(95)<900"],
    },
};

const BASE = __ENV.BASE_URL || "http://localhost:8080";

function pickFirstFlightId(body) {
    // ожидаем массив объектов с полем id
    if (!Array.isArray(body) || body.length === 0) return null;
    const f = body[0];
    return f && (f.id ?? f.ID ?? null);
}

function toYYYYMMDD(scheduledDeparture) {
    // у тебя это может быть ISO строка: "2017-01-01T..."
    if (!scheduledDeparture || typeof scheduledDeparture !== "string") return null;
    return scheduledDeparture.slice(0, 10);
}

export function setup() {
    // берём 1 рейс, потом дергаем детали рейса
    const listRes = http.get(`${BASE}/api/v1/flights?limit=1&offset=0`, {
        tags: { name: "setup_flights_list" },
    });

    if (listRes.status !== 200) {
        fail(`setup: flights list status=${listRes.status} body=${listRes.body}`);
    }

    let listJson;
    try {
        listJson = listRes.json();
    } catch (e) {
        fail(`setup: flights list is not json: ${listRes.body}`);
    }

    const flightId = pickFirstFlightId(listJson);
    if (!flightId) fail(`setup: can't pick flight id from: ${listRes.body}`);

    const flightRes = http.get(`${BASE}/api/v1/flights/${flightId}`, {
        tags: { name: "setup_flight_by_id" },
    });

    if (flightRes.status !== 200) {
        fail(`setup: flight by id status=${flightRes.status} body=${flightRes.body}`);
    }

    let flightJson;
    try {
        flightJson = flightRes.json();
    } catch (e) {
        fail(`setup: flight by id is not json: ${flightRes.body}`);
    }

    // ПОДСТРОЙ ПОД ТВОЙ JSON: ниже наиболее частые варианты имен полей
    const from =
        flightJson.departure_airport ??
        flightJson.DepartureAirport ??
        flightJson.departureAirport ??
        null;

    const to =
        flightJson.arrival_airport ??
        flightJson.ArrivalAirport ??
        flightJson.arrivalAirport ??
        null;

    const depTS =
        flightJson.scheduled_departure ??
        flightJson.ScheduledDeparture ??
        flightJson.scheduledDeparture ??
        null;

    const date = toYYYYMMDD(depTS);

    if (!from || !to || !date) {
        fail(
            `setup: missing from/to/date. from=${from} to=${to} depTS=${depTS}. body=${flightRes.body}`
        );
    }

    return { flightId, from, to, date };
}

export function passengers(data) {
    group("passengers", () => {
        const res = http.get(
            `${BASE}/api/v1/flights/${data.flightId}/passengers?limit=200&offset=0`,
            { tags: { name: "passengers" } }
        );

        check(res, {
            "passengers 200": (r) => r.status === 200,
        });
    });
}

export function flights() {
    group("flights", () => {
        const res = http.get(`${BASE}/api/v1/flights?limit=200&offset=0`, {
            tags: { name: "flights" },
        });

        check(res, { "flights 200": (r) => r.status === 200 });
    });
}

export function search(data) {
    group("search", () => {
        const passengers = 1;
        const fareClass = "Economy";

        const url =
            `${BASE}/api/v1/flights/search` +
            `?from=${encodeURIComponent(data.from)}` +
            `&to=${encodeURIComponent(data.to)}` +
            `&date=${encodeURIComponent(data.date)}` +
            `&passengers=${passengers}` +
            `&fare_class=${encodeURIComponent(fareClass)}` +
            `&limit=50&offset=0`;

        const res = http.get(url, { tags: { name: "search" } });

        check(res, {
            "search 200": (r) => r.status === 200,
        });
    });
}