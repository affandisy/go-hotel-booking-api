import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp up to 200 users
    { duration: '5m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(99)<1000'], // 99% under 1s
  },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  const res = http.get(`${BASE_URL}/hotels`);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
  
  sleep(1);
}