import http from 'k6/http'
import { check, sleep } from 'k6'
import { SharedArray } from 'k6/data'

export const options = {
  vus: 5,
  duration: '15s',
}

const BASE_URL = 'http://localhost:8080'

const users = new SharedArray('users', () => [
  { id: '11111111-1111-1111-1111-111111111111', username: 'u1' },
  { id: '21111111-1111-1111-1111-111111111111', username: 'u2' },
  { id: '31111111-1111-1111-1111-111111111111', username: 'u3' },
  { id: '41111111-1111-1111-1111-111111111111', username: 'u4' },
  { id: '51111111-1111-1111-1111-111111111111', username: 'u5' },
])

let prCounter = 0
function generateSequentialUid(counter) {
  const vu = String(__VU).padStart(6, '0')
  const suffix = String(counter).padStart(6, '0')
  return `11111111-1111-1111-1111-${vu}${suffix}`
}

export function setup() {
  const payload = JSON.stringify({
    team_name: 'backend',
    members: users.map(u => ({
      user_id: u.id,
      username: u.username,
      is_active: true,
    })),
  })

  http.post(`${BASE_URL}/team/add`, payload, {
    headers: { 'Content-Type': 'application/json' },
  })

  return { teamName: 'backend' }
}

export default function (data) {
  const author = users[0]
  const prId = generateSequentialUid(prCounter++)

  const createPayload = JSON.stringify({
    pull_request_id: prId,
    pull_request_name: `Test PR ${prId}`,
    author_id: author.id,
  })

  const resCreate = http.post(`${BASE_URL}/pullRequest/create`, createPayload, {
    headers: { 'Content-Type': 'application/json' },
  })

  check(resCreate, {
    'PR created': r => r.status === 201,
  })

  const reviewer = users[1]
  const resReview = http.get(
      `${BASE_URL}/users/getReview?user_id=${reviewer.id}`
  )

  check(resReview, {
    'getReview OK': r => r.status === 200 || r.status === 404,
  })

  sleep(1)
}
