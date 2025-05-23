import grpc from 'k6/net/grpc';
import { check, sleep, fail } from 'k6';
import { Trend } from 'k6/metrics';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import encoding from 'k6/encoding';

const GRPC_SERVICE_ADDRESS = 'localhost:13013';
const PROTO_IMPORT_PATH = ['../protos/proto'];

const MAIN_PROTO_FILE = 'client/godiploma/cipher/v1/cipher_api.proto';

const USERNAME = __ENV.K6_USERNAME;
const PASSWORD = __ENV.K6_PASSWORD;

const rawPayload = open('../kreya/client/godiploma/cipher/v1/CipherService/CreateStegoImage-request.json');
const basePayload = JSON.parse(rawPayload);


const client = new grpc.Client();

try {
  client.load(PROTO_IMPORT_PATH, MAIN_PROTO_FILE);
} catch (error) {
  console.error(`Proto loading error: ${error}`);
  fail(`Failed to load proto files: ${error}. Check PROTO_IMPORT_PATH ('${PROTO_IMPORT_PATH.join(', ')}') and MAIN_PROTO_FILE ('${MAIN_PROTO_FILE}'). Ensure paths are correct relative to where you run 'k6 run'.`);
}

const createStegoImageDuration = new Trend('grpc_req_create_stego_image_duration_ms', true);

export const options = {
  stages: [
      { duration: '30s', target: 5 },
      { duration: '90s', target: 20 },
      { duration: '30s', target: 0 },
  ],
  thresholds: {
    'grpc_req_duration': ['p(95)<1500'],
    'grpc_req_create_stego_image_duration_ms': ['p(95)<3000'],
    'checks': ['rate>0.99'],
  },
  ext: {
      loadimpact: {
        prometheus: {
          url: 'http://localhost:9090/api/v1/write',
        }
      }
    }
};

export function setup() {
  console.log(`Connecting to gRPC server at: ${GRPC_SERVICE_ADDRESS}`);
  try {
    client.connect(GRPC_SERVICE_ADDRESS, {
      plaintext: true
    });
  } catch (error) {
    console.error(`gRPC connection error in setup: ${error}`);
    fail(`Failed to connect to gRPC server at ${GRPC_SERVICE_ADDRESS}: ${error}`);
  }

  if (!USERNAME || !PASSWORD) {
    console.warn("Warning: K6_USERNAME or K6_PASSWORD environment variable is not set or is empty. Basic Auth might fail if required by the server.");
  }
  console.log(`Using username: ${USERNAME} for Basic Auth.`);
}

export default function () {
    if (!client.isConnected) {
      try {
        client.connect(GRPC_SERVICE_ADDRESS, { plaintext: true });
        console.log(`VU ${__VU} reconnected to gRPC`);
      } catch (error) {
        console.error(`VU ${__VU} failed to connect to gRPC: ${error}`);
        fail(`Connection failed: ${error}`);
        return;
      }
    }

  const payload = JSON.parse(JSON.stringify(basePayload));
  payload.requestId = uuidv4();
  payload.plaintext = `VU: ${__VU}, ITER: ${__ITER}`;

  const params = {
    metadata: {},
    timeout: "60s"
  };

  if (USERNAME && PASSWORD) {
    const credentials = `${USERNAME}:${PASSWORD}`;
    const encodedCredentials = encoding.b64encode(credentials);
    params.metadata['authorization'] = `Basic ${encodedCredentials}`;
  }

  const rpcMethod = 'client.godiploma.cipher.v1.CipherService/CreateStegoImage';
  let response;
  try {
    response = client.invoke(rpcMethod, payload, params);
  } catch (error) {
    console.error(`Error during client.invoke for ${rpcMethod}: ${error}`);
    fail(`gRPC invoke failed for ${rpcMethod}: ${error}`);
    return;
  }

  if (response && response.timings && typeof response.timings.duration !== 'undefined') {
    createStegoImageDuration.add(response.timings.duration);
  }

  const checkRes = check(response, {
    'RPC call succeeded (status OK)': (r) => r && r.status === grpc.StatusOK,
  });

  if (!checkRes) {
    console.error(
      `Request to ${rpcMethod} failed: Status code ${response ? response.status : 'N/A'}, Error: ${response && response.error ? JSON.stringify(response.error) : 'N/A'}`
    );
  }

  sleep(1);
}

export function teardown() {
  client.close();
  console.log("gRPC client connection closed.");
}