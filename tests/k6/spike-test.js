import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 },   // Warm up
    { duration: '1m', target: 10 },    // Normal load
    { duration: '10s', target: 500 },  // Spike!
    { duration: '3m', target: 500 },   // Stay at spike
    { duration: '10s', target: 10 },   // Scale down
    { duration: '3m', target: 10 },    // Recovery
    { duration: '10s', target: 0 },    // Ramp down
  ],
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  const res = http.get(`${BASE_URL}/hotels`);
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}