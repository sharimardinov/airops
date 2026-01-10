import http from "k6/http";
import { check, fail } from "k6";

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const V1 = `${BASE}/api/v1`;

export const options = {
    scenarios: {
        s: {
            executor: "constant-arrival-rate",
            rate: Number(__ENV.RPS || 20),
            timeUnit: "1s",
            duration: __ENV.DURATION || "30s",
            preAllocatedVUs: Number(__ENV.PRE_VUS || 50),
            maxVUs: Number(__ENV.MAX_VUS || 500),
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        dropped_iterations: ["count<10"],
        http_req_duration: ["p(95)<1200"],
    },
};

export function setup() {
    const from = __ENV.FROM;
    const to = __ENV.TO;
    if (!from || !to) fail("set FROM and TO as YYYY-MM-DD, example: -e FROM=2025-12-12 -e TO=2026-01-11");
    return { from, to };
}

export default function (data) {
    const url = `${V1}/stats/routes?limit=10&from=${data.from}&to=${data.to}`;
    const res = http.get(url, { tags: { name: "stats_routes" } });
    check(res, { "200": (r) => r.status === 200 });
}