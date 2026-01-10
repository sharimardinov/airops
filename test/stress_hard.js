import http from "k6/http";
import { check, sleep } from "k6";

const BASE = __ENV.BASE_URL || "http://localhost:8080";
const TARGET = parseInt(__ENV.TARGET_RPS || "5000", 10);

export const options = {
    scenarios: {
        hard: {
            executor: "constant-arrival-rate",
            rate: TARGET,          // 5000 итераций/сек
            timeUnit: "1s",
            duration: "90s",
            preAllocatedVUs: 2000, // стартовый пул
            maxVUs: 6000,          // потолок
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<1200"],
        dropped_iterations: ["count<1"], // если пересек — значит не смог держать rate
    },
};

export default function () {
    // ОДНА ручка, чтобы чисто пробить потолок (как у тебя mixOneReq)
    const res = http.get(`${BASE}/api/v1/flights?limit=200&offset=0`);
    check(res, { "status is 200": (r) => r.status === 200 });
    // sleep(0) не нужен
}