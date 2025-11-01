import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

// Test configuration
export const options = {
stages: [
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 50 },  // Ramp up to 50 users
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 0 },   // Ramp down to 0 users
],
thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'],   // Less than 1% of requests should fail
    errors: ['rate<0.1'],              // Less than 10% errors
},
};

const BASE_URL = 'http://localhost:8080/api/v1';

let authToken = '';
let hotelId = '';
let roomId = '';

export function setup() {
  // Register and login to get token
  const registerPayload = JSON.stringify({
    name: 'K6 Load Test User',
    email: `k6test-${Date.now()}@example.com`,
    password: 'password123',
    role: 'CUSTOMER'
  });

  const registerRes = http.post(`${BASE_URL}/auth/register`, registerPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (registerRes.status === 201) {
    const registerBody = JSON.parse(registerRes.body);
    
    // Login
    const loginPayload = JSON.stringify({
      email: registerBody.data.email,
      password: 'password123',
    });

    const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, {
      headers: { 'Content-Type': 'application/json' },
    });

    if (loginRes.status === 200) {
      const loginBody = JSON.parse(loginRes.body);
      return { token: loginBody.data.token };
    }
  }

  return { token: '' };
}

export default function (data) {
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${data.token}`,
    },
  };

  // Test 1: Health Check
  group('Health Check', () => {
    const res = http.get('http://localhost:8080/health');
    
    check(res, {
      'health check status is 200': (r) => r.status === 200,
      'health check has status ok': (r) => JSON.parse(r.body).status === 'ok',
    }) || errorRate.add(1);
  });

  sleep(1);

  // Test 2: List Hotels
  group('List Hotels', () => {
    const res = http.get(`${BASE_URL}/hotels`);
    
    check(res, {
      'list hotels status is 200': (r) => r.status === 200,
      'list hotels returns array': (r) => {
        const body = JSON.parse(r.body);
        return body.success === true;
      },
    }) || errorRate.add(1);
  });

  sleep(1);

  // Test 3: Get User Bookings (Authenticated)
  if (data.token) {
    group('Get User Bookings', () => {
      const res = http.get(`${BASE_URL}/bookings`, params);
      
      check(res, {
        'get bookings status is 200': (r) => r.status === 200,
        'get bookings is authenticated': (r) => r.status !== 401,
      }) || errorRate.add(1);
    });
  }

  sleep(1);

  // Test 4: Register User (Stress Test)
  group('Register New User', () => {
    const payload = JSON.stringify({
      name: `Load Test User ${__VU}-${__ITER}`,
      email: `loadtest-${__VU}-${__ITER}-${Date.now()}@example.com`,
      password: 'password123',
      role: 'CUSTOMER'
    });

    const res = http.post(`${BASE_URL}/auth/register`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });
    
    check(res, {
      'register status is 201': (r) => r.status === 201,
      'register returns user data': (r) => {
        const body = JSON.parse(r.body);
        return body.success === true && body.data.email !== undefined;
      },
    }) || errorRate.add(1);
  });

  sleep(1);
}

export function teardown(data) {
  console.log('Load test completed');
}