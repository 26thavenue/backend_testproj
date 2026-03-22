import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 20,
  duration: "20s",
};

const baseUrl = __ENV.BASE_URL || "http://localhost:8080/api/v1";

export function setup() {
  const payload = JSON.stringify({
    event_name: "page_view",
    event_description: "Page view events",
  });

  http.post(`${baseUrl}/event-types`, payload, {
    headers: { "Content-Type": "application/json" },
  });
}

export default function () {
  const payload = JSON.stringify({
    event_name: "page_view",
    user_id: "load-user",
    properties: { page: "/home" },
  });

  const res = http.post(`${baseUrl}/events`, payload, {
    headers: { "Content-Type": "application/json" },
  });

  check(res, { "status is 201": (r) => r.status === 201 });
  sleep(0.1);
}
