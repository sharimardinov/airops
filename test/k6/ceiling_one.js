import http from "k6/http";
import { check } from "k6";

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const V1 = `${BASE}/api/v1`;

export const options = {
    scenarios: {
        s: {
            executor: "ramping-arrival-rate",
            timeUnit: "1s",
            startRate: Number(__ENV.START_RPS || 200),
            stages: [
                { target: Number(__ENV.R1 || 500), duration: "20s" },
                { target: Number(__ENV.R2 || 1000), duration: "20s" },
                { target: Number(__ENV.R3 || 2000), duration: "20s" },
                { target: 0, duration: "10s" },
            ],
            preAllocatedVUs: Number(__ENV.PRE_VUS || 500),
            maxVUs: Number(__ENV.MAX_VUS || 5000),
            gracefulStop: "30s",
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        dropped_iterations: ["count==0"],
    },
};

export default function () {
    const res = http.get(`${V1}/flights?limit=10&offset=0`, { tags: { name: "flights_list" } });
    check(res, { "200": (r) => r.status === 200 });
}