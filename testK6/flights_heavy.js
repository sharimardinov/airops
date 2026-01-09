import http from "k6/http";
import { check } from "k6";

export const options = {
    scenarios: {
        flights: {
            executor: "constant-arrival-rate",
            rate: __ENV.RPS ? Number(__ENV.RPS) : 1000,
            timeUnit: "1s",
            duration: "60s",
            preAllocatedVUs: 300,
            maxVUs: 3000,
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<500"],
        dropped_iterations: ["count<10"],
    },
};

export default function () {
    const res = http.get("http://localhost:8080/api/v1/flights?limit=200&offset=0");
    check(res, { "200": r => r.status === 200 });
}