import http from "k6/http";
import { check } from "k6";

export const options = {
    scenarios: {
        s: {
            executor: "constant-arrival-rate",
            rate: 2000,
            timeUnit: "1s",
            duration: "30s",
            preAllocatedVUs: 200,
            maxVUs: 2000,
        },
    },
    noConnectionReuse: false,
    userAgent: "k6",
};

export default function () {
    const r = http.get("http://localhost:8080/api/v1/stats/routes?from=2025-12-01&to=2026-01-01&limit=10");
    check(r, { "200": (x) => x.status === 200 });
}