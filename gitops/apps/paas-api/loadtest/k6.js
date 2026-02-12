import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    ramp: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "30s", target: 10 },
        { duration: "30s", target: 30 },
        { duration: "60s", target: 60 }, // should trigger HPA if endpoint costs CPU
        { duration: "30s", target: 0 },
      ],
      gracefulRampDown: "10s",
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"], // adjust later to your reality
  },
};

const BASE_URL = __ENV.BASE_URL; // injected from the Job

export default function () {
  const res = http.get(`${BASE_URL}/work?ms=50`);  // BEST for HPA demo
  check(res, { "status is 200": (r) => r.status === 200 });
  sleep(0.1);
}
