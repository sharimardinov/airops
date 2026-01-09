import http from "k6/http";
import { check } from "k6";

export const options = {
    scenarios: {
        health: {
            executor: "constant-arrival-rate",
            rate: __ENV.RPS ? Number(__ENV.RPS) : 2000,
            timeUnit: "1s",
            duration: "30s",
            preAllocatedVUs: 500,
            maxVUs: 5000,
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.001"],
        http_req_duration: ["p(95)<100"],
        dropped_iterations: ["count<1"],
    },
};

export default function () {
    const res = http.get("http://localhost:8080/health");
    check(res, { "200": r => r.status === 200 });
}