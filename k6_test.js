import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
  stages: [
    // Ramp-up phase 1: Warm-up
    { duration: "30s", target: 100 }, // Ramp up to 100 VUs over 30 seconds
    { duration: "1m", target: 100 }, // Stay at 100 VUs for 1 minute

    // Ramp-up phase 2: Gradual increase to peak
    { duration: "1m", target: 1000 }, // Ramp up to 1000 VUs over 1 minutes

    // Peak load phase
    { duration: "1m", target: 1000 }, // Stay at 1000 VUs for 1 minute

    // Ramp-down phase
    { duration: "1m", target: 0 }, // Ramp down to 0 VUs over 1 minute
  ],
};

export default function () {
  const BASE_URL = "http://localhost:8080";

  const res = http.get(`${BASE_URL}/produk`);
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}
