import './style.css';
import sodium from "libsodium-wrappers";

const defaultToken = "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJYIiwic3ViIjoiMTM3NDg2NDcxODU5MTA5ODg4NiIsImV4cCI6IjE3NjQ0NTU2MDMiLCJuYmYiOiIxNzY0NDU1MzAzIiwidWEiOiJNb3ppbGxhXC81LjAgKE1hY2ludG9zaDsgSW50ZWwgTWFjIE9TIFggMTBfMTVfNykgQXBwbGVXZWJLaXRcLzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZVwvMTQyLjAuMC4wIFNhZmFyaVwvNTM3LjM2In0.cHQXTovky-wueJJaWUEZBf-xXAangd0e4daweiTIKO4";
const defaultKeyB64 = "ST4hdGP1V3A5az3dM4GgjvoPRiaBl9TZUwGkfPy9/0E=";

const app = document.querySelector<HTMLDivElement>('#app');

if (!app) {
	throw new Error("App container missing");
}

app.innerHTML = `
  <main class="min-h-screen bg-slate-950 text-slate-50 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-3xl space-y-6">
      <div class="space-y-1">
        <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">XChat websocket</p>
        <h1 class="text-3xl font-semibold">Reverse XChat</h1>
        <p class="text-slate-300">Provide your token and decryption key, then click connect to start the websocket.</p>
      </div>

      <form id="connect-form" class="space-y-4 rounded-xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-sky-900/30 backdrop-blur">
        <label class="block space-y-2">
          <span class="text-sm font-medium text-slate-200">Token</span>
          <textarea id="token-input" name="token" rows="4" required class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40">${defaultToken}</textarea>
        </label>

        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-slate-200">Decryption key (base64)</span>
            <button
              id="save-key-btn"
              type="button"
              class="rounded-md border border-slate-700 bg-slate-800 px-3 py-1 text-xs font-semibold text-slate-100 hover:border-sky-500 hover:text-white focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-slate-950"
            >
              Save key
            </button>
          </div>
          <input
            id="key-input"
            name="key"
            type="text"
            value="${defaultKeyB64}"
            required
            class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
          />
          <p id="key-preview" class="text-xs text-slate-400 font-mono">Length: - | hex: -</p>
          <label class="block space-y-1">
            <span class="text-xs font-medium text-slate-300">Optional public key (base64)</span>
            <input
              id="pubkey-input"
              name="pubkey"
              type="text"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
            />
          </label>
          <label class="block space-y-1">
            <span class="text-xs font-medium text-slate-300">Saved keys</span>
            <div class="flex items-center gap-2">
              <select
                id="key-select"
                class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
              ></select>
              <button
                id="delete-key-btn"
                type="button"
                class="whitespace-nowrap rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-100 hover:border-rose-500 hover:text-white focus:outline-none focus:ring-2 focus:ring-rose-500 focus:ring-offset-2 focus:ring-offset-slate-950"
              >
                Delete
              </button>
            </div>
          </label>
        </div>

        <div class="flex items-center gap-3">
          <button
            id="connect-btn"
            type="submit"
            class="inline-flex items-center justify-center rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
          >
            Connect & Listen
          </button>
          <p class="text-xs text-slate-400">Click the button to initiate the websocket connection.</p>
        </div>
      </form>

      <pre id="status-log" class="whitespace-pre-wrap rounded-lg border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-200">Waiting to connect...</pre>

      <section class="grid gap-4 lg:grid-cols-2 w-full">
        <div class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20">
          <div class="mb-2 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-slate-100">Latest event JSON</h2>
            <span class="rounded-full bg-slate-800 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-slate-300">Raw</span>
          </div>
          <pre id="event-output" class="h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100">No event received yet.</pre>
        </div>
        <div class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20">
          <div class="mb-2 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-slate-100">Decrypted payload</h2>
            <span class="rounded-full bg-emerald-800/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-emerald-100">Plaintext</span>
          </div>
          <pre id="decrypted-output" class="h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100">No decrypted payload yet.</pre>
        </div>
      </section>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full">
        <div class="mb-3 space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">Manual decode</p>
          <h2 class="text-lg font-semibold text-slate-100">Paste a base64 payload</h2>
          <p class="text-sm text-slate-300">Base64-decodes the payload and parses it with the generic struct parser into JSON (no decryption attempted).</p>
        </div>
        <form id="manual-form" class="space-y-3">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">Ciphertext bundle (base64)</span>
            <textarea
              id="manual-input"
              name="ciphertext"
              rows="3"
              placeholder="Paste base64-encoded ciphertext bundle here"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
            ></textarea>
          </label>
          <div class="flex items-center gap-3">
            <button
              id="manual-decrypt-btn"
              type="submit"
              class="inline-flex items-center justify-center rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Parse message
            </button>
            <p class="text-xs text-slate-400">We will base64-decode and parse to JSON below.</p>
          </div>
        </form>
        <pre id="manual-output" class="mt-4 h-48 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100">Awaiting ciphertext...</pre>
      </section>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full">
        <div class="mb-3 space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-purple-300">Conversation key event</p>
          <h2 class="text-lg font-semibold text-slate-100">Decrypt conversation keys from event JSON</h2>
          <p class="text-sm text-slate-300">Paste the full event JSON (or base64 of it). We will extract blobs at 7→3→2[*]→2 and secretbox-decrypt them with the selected 32-byte key (nonce|ciphertext|mac).</p>
        </div>
        <form id="conv-event-form" class="space-y-3">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">Event JSON or base64-encoded JSON</span>
            <textarea
              id="conv-event-input"
              name="conv-event"
              rows="6"
              placeholder="Paste the event JSON here"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-500/40"
            ></textarea>
          </label>
          <div class="flex items-center gap-3">
            <button
              id="conv-event-btn"
              type="submit"
              class="inline-flex items-center justify-center rounded-lg bg-purple-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Decrypt event blobs
            </button>
            <p class="text-xs text-slate-400">Uses the selected key (32-byte secret or 64-byte priv+pub).</p>
          </div>
        </form>
        <pre id="conv-event-output" class="mt-4 h-64 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100">Awaiting event JSON...</pre>
      </section>
    </div>
  </main>
`;

const form = document.querySelector<HTMLFormElement>('#connect-form');
const tokenInput = document.querySelector<HTMLTextAreaElement>('#token-input');
const keyInput = document.querySelector<HTMLInputElement>('#key-input');
const pubKeyInput = document.querySelector<HTMLInputElement>('#pubkey-input');
const keyPreview = document.querySelector<HTMLParagraphElement>('#key-preview');
const statusLog = document.querySelector<HTMLPreElement>('#status-log');
const connectButton = document.querySelector<HTMLButtonElement>('#connect-btn');
const keySelect = document.querySelector<HTMLSelectElement>('#key-select');
const saveKeyButton = document.querySelector<HTMLButtonElement>('#save-key-btn');
const deleteKeyButton = document.querySelector<HTMLButtonElement>('#delete-key-btn');
const eventOutput = document.querySelector<HTMLPreElement>('#event-output');
const decryptedOutput = document.querySelector<HTMLPreElement>('#decrypted-output');
const manualForm = document.querySelector<HTMLFormElement>('#manual-form');
const manualInput = document.querySelector<HTMLTextAreaElement>('#manual-input');
const manualOutput = document.querySelector<HTMLPreElement>('#manual-output');
const convEventForm = document.querySelector<HTMLFormElement>('#conv-event-form');
const convEventInput = document.querySelector<HTMLTextAreaElement>('#conv-event-input');
const convEventOutput = document.querySelector<HTMLPreElement>('#conv-event-output');

if (!form || !tokenInput || !keyInput || !pubKeyInput || !keyPreview || !statusLog || !connectButton || !keySelect || !saveKeyButton || !deleteKeyButton || !eventOutput || !decryptedOutput || !manualForm || !manualInput || !manualOutput || !convEventForm || !convEventInput || !convEventOutput) {
	throw new Error("Failed to initialize app controls");
}

const logStatus = (message: string) => {
	const timestamp = new Date().toLocaleTimeString();
	statusLog.textContent = `[${timestamp}] ${message}`;
	console.info(message);
};

const STORAGE_KEYS = "xchatSavedKeys";

type SavedKey = {
	priv: string;
	pub?: string;
};

const loadSavedKeys = (): SavedKey[] => {
	const raw = localStorage.getItem(STORAGE_KEYS);
	if (!raw) return [];
	try {
		const parsed = JSON.parse(raw);
		if (!Array.isArray(parsed)) return [];
		return parsed
			.map((item) => {
				if (typeof item === "string") return { priv: item };
				if (item && typeof item.priv === "string") {
					return { priv: String(item.priv), pub: item.pub ? String(item.pub) : undefined };
				}
				return null;
			})
			.filter((v): v is SavedKey => Boolean(v && v.priv));
	} catch {
		return [];
	}
};

const persistKeys = (keys: SavedKey[]) => {
	localStorage.setItem(STORAGE_KEYS, JSON.stringify(keys));
};

let savedKeys = loadSavedKeys();
if (!savedKeys.length) savedKeys = [{ priv: defaultKeyB64 }];

const renderSavedKeyOptions = (activeIndex = 0) => {
	keySelect.innerHTML = "";
	const frag = document.createDocumentFragment();
	const preview = (val: string) => (val.length > 14 ? `${val.slice(0, 8)}…${val.slice(-6)}` : val);

	if (!savedKeys.length) {
		const opt = document.createElement("option");
		opt.value = "";
		opt.textContent = "No saved keys";
		frag.appendChild(opt);
		keySelect.disabled = true;
		deleteKeyButton.disabled = true;
		keyInput.value = "";
		pubKeyInput.value = "";
	} else {
		keySelect.disabled = false;
		deleteKeyButton.disabled = false;
		savedKeys.forEach((k, idx) => {
			const opt = document.createElement("option");
			opt.value = String(idx);
			opt.textContent = k.pub ? `Priv ${preview(k.priv)} | Pub ${preview(k.pub)}` : `Priv ${preview(k.priv)}`;
			frag.appendChild(opt);
		});
		const safeIndex = activeIndex < savedKeys.length ? activeIndex : 0;
		keySelect.value = String(safeIndex);
		const activeKey = savedKeys[safeIndex];
		if (activeKey) {
			keyInput.value = activeKey.priv;
			pubKeyInput.value = activeKey.pub ?? "";
		}
	}

	keySelect.appendChild(frag);
};

const upsertKey = (privB64: string, pubB64?: string) => {
	const trimmed = privB64.trim();
	if (!trimmed) return;
	const entry: SavedKey = { priv: trimmed, pub: pubB64?.trim() || undefined };
	// Replace existing entry with same priv if found
	const idx = savedKeys.findIndex((k) => k.priv === entry.priv);
	if (idx >= 0) {
		savedKeys[idx] = entry;
		renderSavedKeyOptions(idx);
	} else {
		savedKeys = [entry, ...savedKeys];
		renderSavedKeyOptions(0);
	}
	persistKeys(savedKeys);
};

renderSavedKeyOptions(0);

const jsonReplacer = (_key: string, value: unknown) => (typeof value === "bigint" ? value.toString() : value);
const renderEventJson = (json: string) => {
	eventOutput.textContent = json;
};
const renderDecryptedPayload = (text: string) => {
	decryptedOutput.textContent = text;
};
const renderManualPayload = (text: string) => {
	manualOutput.textContent = text;
};
const updateKeyPreview = () => {
	const value = keyInput.value.trim();
	if (!value) {
		keyPreview.textContent = "Length: 0 | hex: (empty)";
		return;
	}
	try {
		const bytes = base64ToUint8Array(value);
		keyPreview.textContent = `Length: ${bytes.length} | hex: ${bytesToHex(bytes)}`;
	} catch {
		keyPreview.textContent = "Invalid base64 (cannot decode)";
	}
};
const renderConvEventOutput = (text: string) => {
	convEventOutput.textContent = text;
};

let activeSocket: WebSocket | null = null;

keyInput.addEventListener('input', updateKeyPreview);
updateKeyPreview();

keySelect.addEventListener('change', () => {
	const idx = Number(keySelect.value);
	const entry = savedKeys[idx];
	if (entry) {
		keyInput.value = entry.priv;
		pubKeyInput.value = entry.pub ?? "";
	}
	updateKeyPreview();
});

saveKeyButton.addEventListener('click', () => {
	const keyValue = keyInput.value.trim();
	const pubValue = pubKeyInput.value.trim();
	if (!keyValue) {
		logStatus("Enter a base64 decryption key before saving.");
		return;
	}
	try {
		base64ToUint8Array(keyValue);
	} catch {
		logStatus("Key must be valid base64.");
		return;
	}
	if (pubValue) {
		try {
			base64ToUint8Array(pubValue);
		} catch {
			logStatus("Public key must be valid base64.");
			return;
		}
	}
	upsertKey(keyValue, pubValue || undefined);
	logStatus("Key saved.");
	updateKeyPreview();
});

deleteKeyButton.addEventListener('click', () => {
	const idx = Number(keySelect.value);
	if (Number.isNaN(idx) || idx < 0 || idx >= savedKeys.length) {
		logStatus("No saved key selected to delete.");
		return;
	}
	savedKeys.splice(idx, 1);
	persistKeys(savedKeys);
	const nextIdx = Math.max(0, Math.min(idx, savedKeys.length - 1));
	renderSavedKeyOptions(nextIdx);
	logStatus("Key deleted.");
	updateKeyPreview();
});

form.addEventListener('submit', async (event) => {
	event.preventDefault();

	const token = tokenInput.value.trim();
	const decryptionKeyB64 = keyInput.value.trim();

	if (!token || !decryptionKeyB64) {
		logStatus("Please provide both token and decryption key.");
		return;
	}

	connectButton.disabled = true;
	logStatus("Connecting with provided token and key...");
	renderEventJson("Waiting for event...");
	renderDecryptedPayload("Waiting for decrypted payload...");

	try {
		activeSocket?.close();
		activeSocket = await websocketReverseEngineering(
			token,
			decryptionKeyB64,
			logStatus,
			(evt) => renderEventJson(JSON.stringify(evt, jsonReplacer, 2)),
			(payload) => renderDecryptedPayload(payload),
		);
	} catch (err) {
		const message = err instanceof Error ? err.message : String(err);
		logStatus(`Connection failed: ${message}`);
	} finally {
		connectButton.disabled = false;
	}
});

manualForm.addEventListener('submit', async (event) => {
	event.preventDefault();

	const ciphertextB64 = manualInput.value.trim();

	if (!ciphertextB64) {
		renderManualPayload("Please paste the base64 payload (same as websocket binary, base64-encoded).");
		return;
	}

	const bytes = base64ToUint8Array(ciphertextB64);

	try {
		const decoder = new Decoder(bytes);
		const obj = decoder.readStruct(xchatEventSchema);
		const bodyJson = JSON.stringify(obj, jsonReplacer, 2);
		renderManualPayload(bodyJson);
		return;
	} catch (err) {
	}

	try {
		const decoder = new Decoder(bytes);
		const obj = decoder.readStructGeneric();
		const bodyJson = JSON.stringify(obj, jsonReplacer, 2);
		renderManualPayload(bodyJson);
	} catch (err) {
		renderManualPayload(`Parse Error: ${err}`);
	}
});

convEventForm.addEventListener('submit', async (event) => {
	event.preventDefault();

	const privScalarB64 = keyInput.value.trim();
	const raw = convEventInput.value.trim();
	if (!raw) {
		renderConvEventOutput("Please paste an event JSON or base64-encoded JSON.");
		return;
	}

	if (!privScalarB64) {
		renderConvEventOutput("Provide your 32-byte P-256 private scalar in the key field above.");
		return;
	}

	let eventJson: any = null;
	const tryParseJson = (text: string) => {
		try {
			return JSON.parse(text);
		} catch {
			return null;
		}
	};

	eventJson = tryParseJson(raw);

	if (!eventJson) {
			const bytes = base64ToUint8Array(raw);

				const decoder = new Decoder(bytes);
				eventJson = decoder.readStruct(xchatEventSchema)

			// Fallback: treat base64 as JSON string.
			if (!eventJson) {
				try {
					const decoded = new TextDecoder().decode(bytes);
					eventJson = tryParseJson(decoded);
				} catch {
					/* ignore */
				}
			}
	}

	if (!eventJson) {
		renderConvEventOutput("Could not parse JSON (neither plain nor base64).");
		return;
	}

	const evt = eventJson.event ?? eventJson;

	console.log(evt);

	const payloads = evt?.payload?.encryptedConversationKey?.encryptedKeyPayload;

	if (!Array.isArray(payloads) || payloads.length === 0) {
		renderConvEventOutput("No encrypted conversation key payloads found at payload.encryptedConversationKey.encryptedKeyPayload.");
		return;
	}

	const lines: string[] = [];
	for (const entry of payloads) {
		const label = entry?.userId ?? "(unknown user)";
		const keyB64 = entry?.keyB64;
		if (!keyB64 || typeof keyB64 !== "string") {
			lines.push(`✗ ${label}: missing keyB64`);
			continue;
		}

		try {
			const ck = await unwrapConversationKey(keyB64, privScalarB64);
			const ckB64 = bytesToBase64(ck);
			lines.push(`✓ ${label}: ${ckB64}`);
		} catch (err) {
			console.error(err);
			const msg = err instanceof Error ? err.message : String(err);
			lines.push(`✗ ${label}: ${msg}`);
		}
	}

	renderConvEventOutput(lines.join("\n"));
});

function bytesToBase64(bytes: Uint8Array) {
	let bin = "";
	for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
	return btoa(bin);
}

async function importP256PrivateKeyFromScalar(scalarB64: string) {
	const privBytes = base64ToUint8Array(scalarB64);
	if (privBytes.length !== 32) {
		throw new Error(`Expected 32-byte scalar, got ${privBytes.length}`);
	}
	const pkcs8 = buildP256Pkcs8FromScalar(privBytes);
	return crypto.subtle.importKey("pkcs8", pkcs8, { name: "ECDH", namedCurve: "P-256" }, false, ["deriveBits"]);
}

function concatBytes(...parts: Uint8Array[]) {
	const len = parts.reduce((n, p) => n + p.length, 0);
	const out = new Uint8Array(len);
	let o = 0;
	for (const p of parts) {
		out.set(p, o);
		o += p.length;
	}
	return out;
}

function asArrayBuffer(view: Uint8Array) {
	return view.buffer.slice(view.byteOffset, view.byteOffset + view.byteLength) as ArrayBuffer;
}

function derLen(len: number): Uint8Array {
	if (len < 128) return Uint8Array.of(len);
	const bytes: number[] = [];
	let n = len;
	while (n > 0) {
		bytes.unshift(n & 0xff);
		n >>>= 8;
	}
	return Uint8Array.of(0x80 | bytes.length, ...bytes);
}

function derNode(tag: number, content: Uint8Array) {
	return concatBytes(Uint8Array.of(tag), derLen(content.length), content);
}

function derSeq(...parts: Uint8Array[]) {
	return derNode(0x30, concatBytes(...parts));
}

function derOctetString(bytes: Uint8Array) {
	return derNode(0x04, bytes);
}

function derInteger(value: number) {
	return Uint8Array.of(0x02, 0x01, value & 0xff);
}

function buildP256Pkcs8FromScalar(priv: Uint8Array) {
	// OIDs
	const oidEcPublicKey = Uint8Array.from([0x06, 0x07, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x02, 0x01]);
	const oidPrime256v1 = Uint8Array.from([0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07]);

	// ECPrivateKey ::= SEQUENCE { version, privateKey, [0] parameters }
	const ecPrivateKey = derSeq(derInteger(1), derOctetString(priv), derNode(0xa0, oidPrime256v1));

	// PrivateKeyInfo ::= SEQUENCE { version, algorithm, privateKey }
	const algorithmId = derSeq(oidEcPublicKey, oidPrime256v1);
	return derSeq(derInteger(0), algorithmId, derOctetString(ecPrivateKey));
}

async function sha256(bytes: Uint8Array) {
	const digest = await crypto.subtle.digest("SHA-256", asArrayBuffer(bytes));
	return new Uint8Array(digest);
}

async function kdf2Sha256(shared: Uint8Array, other: Uint8Array, length: number) {
	const chunks: Uint8Array[] = [];
	let counter = 1;
	let total = 0;
	while (total < length) {
		const counterBytes = new Uint8Array(4);
		new DataView(counterBytes.buffer).setUint32(0, counter, false);
		const digest = await sha256(concatBytes(shared, counterBytes, other));
		chunks.push(digest);
		total += digest.length;
		counter += 1;
	}
	const material = concatBytes(...chunks);
	return material.slice(0, length);
}

async function unwrapConversationKey(keyB64: string, privScalarB64: string) {
	const blob = base64ToUint8Array(keyB64);
	if (blob.length < 65 + 16) {
		throw new Error(`Unexpected keyB64 length=${blob.length}`);
	}

	const ephPub = blob.slice(0, 65);
	const cipherAndTag = blob.slice(65);

	const privKey = await importP256PrivateKeyFromScalar(privScalarB64);
	const pubKey = await crypto.subtle.importKey("raw", ephPub, { name: "ECDH", namedCurve: "P-256" }, false, []);

	const sharedBits = await crypto.subtle.deriveBits({ name: "ECDH", public: pubKey }, privKey, 256);
	const shared = new Uint8Array(sharedBits);

	const keyNonce = await kdf2Sha256(shared, ephPub, 32);
	const aesKeyBytes = keyNonce.slice(0, 16);
	const iv = keyNonce.slice(16);

	const aesKey = await crypto.subtle.importKey("raw", asArrayBuffer(aesKeyBytes), "AES-GCM", false, ["decrypt"]);

	let plaintext: Uint8Array;
	try {
		const plaintextBuf = await crypto.subtle.decrypt({ name: "AES-GCM", iv: asArrayBuffer(iv), tagLength: 128 }, aesKey, asArrayBuffer(cipherAndTag));
		plaintext = new Uint8Array(plaintextBuf);
	} catch (err) {
		const msg = err instanceof Error ? err.message : String(err);
		throw new Error(`AES-GCM decrypt failed (check key/iv/payload): ${msg}`);
	}

	if (plaintext.length !== 32) {
		throw new Error(`Unexpected conversation key length: ${plaintext.length}`);
	}

	return plaintext;
}

function base64ToUint8Array(b64: string) {
	// Normalize base64url, strip whitespace, and fix padding.
	b64 = b64.trim().replace(/\s+/g, "").replace(/-/g, "+").replace(/_/g, "/");
	const pad = b64.length % 4;
	if (pad > 0) b64 += "=".repeat(4 - pad); // be permissive; let atob validate
	const binary = atob(b64);
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i) & 0xff;
	return bytes;
}

function utf8ToBase64(s: string) {
	const bytes = new TextEncoder().encode(s);
	let binary = "";
	for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i]);
	return btoa(binary);
}

function bytesToHex(bytes: Uint8Array) {
	return Array.from(bytes)
		.map((b) => b.toString(16).padStart(2, "0"))
		.join("");
}

function decryptAfterExtractingNonce(nc: Uint8Array, key: Uint8Array) {
	if (nc.length < sodium.crypto_secretbox_NONCEBYTES + sodium.crypto_secretbox_MACBYTES) {
		throw new Error("Short message");
	}
	const nonce = nc.slice(0, sodium.crypto_secretbox_NONCEBYTES);
	const ciphertext = nc.slice(sodium.crypto_secretbox_NONCEBYTES);
	return sodium.crypto_secretbox_open_easy(ciphertext, nonce, key);
}

// Strip control/framing bytes that occasionally wrap the plaintext
function trimBinaryEdges(bytes: Uint8Array) {
	const isControl = (b: number) => b < 0x20 && b !== 0x09 && b !== 0x0a && b !== 0x0d;
	let start = 0;
	let end = bytes.length;

	while (start < end && isControl(bytes[start])) start += 1;
	while (end > start && isControl(bytes[end - 1])) end -= 1;

	return bytes.slice(start, end);
}

// ---- Thrift decoder helpers (unchanged structure, but binary-safe strings) ----
export const T = {
	STOP: 0,
	BOOL: 2,
	BYTE: 3,
	DOUBLE: 4,
	I16: 6,
	I32: 8,
	I64: 10,
	STRING: 11,
	STRUCT: 12,
	MAP: 13,
	SET: 14,
	LIST: 15,
	UTF8: 16,
	UTF16: 17,
} as const;
export type TCode = (typeof T)[keyof typeof T];

export type FieldSchema = {
	id: number;
	key: string;
	type: TCode;
	schema?: FieldSchema[];
	elemType?: TCode;
	elemSchema?: FieldSchema[];
	keyType?: TCode;
	keySchema?: FieldSchema[];
	valType?: TCode;
	valSchema?: FieldSchema[];
	encodeElem?: (v: any) => Uint8Array;
	encodeKey?: (v: any) => Uint8Array;
	encodeVal?: (v: any) => Uint8Array;
};

const concat = (...parts: Uint8Array[]) => {
	const len = parts.reduce((n, p) => n + p.length, 0);
	const out = new Uint8Array(len);
	let o = 0;
	for (const p of parts) {
		out.set(p, o);
		o += p.length;
	}
	return out;
};

const enc = {
	i16: (v: number) => {
		const b = new Uint8Array(2);
		new DataView(b.buffer).setInt16(0, v, false);
		return b;
	},
	i32: (v: number) => {
		const b = new Uint8Array(4);
		new DataView(b.buffer).setInt32(0, v, false);
		return b;
	},
	i64: (v: bigint | number) => {
		const b = new Uint8Array(8);
		new DataView(b.buffer).setBigInt64(0, BigInt(v), false);
		return b;
	},
	dbl: (v: number) => {
		const b = new Uint8Array(8);
		new DataView(b.buffer).setFloat64(0, v, false);
		return b;
	},
	bool: (v: boolean) => Uint8Array.of(v ? 1 : 0),
	str: (s: string) => {
		const utf = new TextEncoder().encode(s);
		return concat(enc.i32(utf.length), utf);
	},
};

const field = (type: TCode, id: number, valBytes: Uint8Array) =>
	concat(Uint8Array.of(type), enc.i16(id), valBytes);

const stop = () => Uint8Array.of(T.STOP);

const defaultEncode = (t: TCode) => (v: any): Uint8Array => {
	switch (t) {
		case T.STRING:
		case T.UTF8:
		case T.UTF16:
			return enc.str(String(v));
		case T.BOOL:
			return enc.bool(!!v);
		case T.BYTE: {
			const n = Number(v);
			const b = (n | 0) & 0xff;
			return Uint8Array.of(b);
		}
		case T.I16:
			return enc.i16(Number(v));
		case T.I32:
			return enc.i32(Number(v));
		case T.I64:
			return enc.i64(v);
		case T.DOUBLE:
			return enc.dbl(Number(v));
		default:
			throw new Error(`Provide encoder for complex type ${t}`);
	}
};

// Encode JS object to Thrift struct (binary protocol, no message envelope)
export function encodeStruct(obj: any, schema: FieldSchema[]): Uint8Array {
	const chunks: Uint8Array[] = [];

	for (const f of schema) {
		const v = obj[f.key];
		if (v === undefined || v === null) continue;

		let vb: Uint8Array;

		switch (f.type) {
			case T.BOOL:
				vb = enc.bool(!!v);
				break;
			case T.BYTE: {
				const n = Number(v);
				vb = Uint8Array.of((n | 0) & 0xff);
				break;
			}
			case T.I16:
				vb = enc.i16(Number(v));
				break;
			case T.I32:
				vb = enc.i32(Number(v));
				break;
			case T.I64:
				vb = enc.i64(v);
				break;
			case T.DOUBLE:
				vb = enc.dbl(Number(v));
				break;
			case T.STRING:
			case T.UTF8:
			case T.UTF16:
				vb = enc.str(String(v));
				break;
			case T.STRUCT:
				vb = encodeStruct(v, f.schema || []);
				break;
			case T.LIST:
			case T.SET: {
				const elemType = f.elemType!;
				const encodeElem =
					f.encodeElem ||
					(elemType === T.STRUCT
						? (x: any) => encodeStruct(x, f.elemSchema || [])
						: defaultEncode(elemType));
				const elems = (v as any[]).map((x: any) => encodeElem(x));
				const header = concat(Uint8Array.of(elemType), enc.i32(elems.length));
				vb = concat(header, ...elems);
				break;
			}
			case T.MAP: {
				const kt = f.keyType!;
				const vt = f.valType!;
				const encodeKey =
					f.encodeKey ||
					(kt === T.STRUCT
						? (x: any) => encodeStruct(x, f.keySchema || [])
						: defaultEncode(kt));
				const encodeVal =
					f.encodeVal ||
					(vt === T.STRUCT
						? (x: any) => encodeStruct(x, f.valSchema || [])
						: defaultEncode(vt));

				const entries: Uint8Array[] = [];
				for (const [k, val] of Object.entries(v)) {
					entries.push(encodeKey(k));
					entries.push(encodeVal(val));
				}
				const header = concat(Uint8Array.of(kt), Uint8Array.of(vt), enc.i32(entries.length / 2));
				vb = concat(header, ...entries);
				break;
			}
			default:
				throw new Error(`Unsupported type ${f.type}`);
		}

		chunks.push(field(f.type, f.id, vb));
	}

	chunks.push(stop());
	return concat(...chunks);
}

export const MsgType = {
	CALL: 1,
	REPLY: 2,
	EXCEPTION: 3,
	ONEWAY: 4,
} as const;
export type MsgTypeCode = (typeof MsgType)[keyof typeof MsgType];

export type MessageBegin = {
	name: string;
	type: MsgTypeCode;
	seqid: number;
};

export class Decoder {
	private dv: DataView;
	private o = 0;

	constructor(buf: Uint8Array) {
		this.dv = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
	}

	private ensure(n: number) {
		if (this.o + n > this.dv.byteLength) {
			throw new RangeError(`Need ${n} bytes at ${this.o}, len ${this.dv.byteLength}`);
		}
	}

	private readByte = () => {
		this.ensure(1);
		const v = this.dv.getUint8(this.o);
		this.o += 1;
		return v;
	};

	private readI8 = () => {
		this.ensure(1);
		const v = this.dv.getInt8(this.o);
		this.o += 1;
		return v;
	};

	private readI16 = () => {
		this.ensure(2);
		const v = this.dv.getInt16(this.o, false);
		this.o += 2;
		return v;
	};

	private readI32 = () => {
		this.ensure(4);
		const v = this.dv.getInt32(this.o, false);
		this.o += 4;
		return v;
	};

	private readI64 = () => {
		this.ensure(8);
		const v = this.dv.getBigInt64(this.o, false);
		this.o += 8;
		return v;
	};

	private readDbl = () => {
		this.ensure(8);
		const v = this.dv.getFloat64(this.o, false);
		this.o += 8;
		return v;
	};

	private readBytes = (len: number) => {
		this.ensure(len);
		const bytes = new Uint8Array(this.dv.buffer, this.dv.byteOffset + this.o, len);
		this.o += len;
		return bytes;
	};

	private readStr = (raw = false) => {
		const len = this.readI32();
		const bytes = this.readBytes(len);
		return raw ? bytes : new TextDecoder("utf-8").decode(bytes);
	};

	readMessageBeginStrict(): MessageBegin {
		const word = this.readI32();
		const version = word & 0xffff0000;
		if (version !== 0x80010000) {
			throw new Error(`Not a strict Thrift message (version word: 0x${word.toString(16)})`);
		}
		const type = (word & 0x000000ff) as MsgTypeCode;
		const nameLen = this.readI32();
		if (nameLen < 0) throw new Error(`Negative method name length: ${nameLen}`);
		this.ensure(nameLen);
		const bytes = new Uint8Array(this.dv.buffer, this.dv.byteOffset + this.o, nameLen);
		this.o += nameLen;
		const name = new TextDecoder().decode(bytes);
		const seqid = this.readI32();
		return { name, type, seqid };
	}

	readMessageEndStrict(): void {
		/* no-op */
	}

	readStruct(schema: FieldSchema[]): any {
		const out: any = {};
		while (true) {
			const t = this.readByte();
			if (t === T.STOP) break;
			const id = this.readI16();
			const f = schema.find((x) => x.id === id);
			if (!f) {
				this.skip(t as TCode);
				continue;
			}
			out[f.key] = this.readValue(t as TCode, f);
		}
		return out;
	}

	private readValue(t: TCode, f: FieldSchema): any {
		switch (t) {
			case T.BOOL:
				return this.readByte() !== 0;
			case T.BYTE:
				return this.readI8();
			case T.I16:
				return this.readI16();
			case T.I32:
				return this.readI32();
			case T.I64:
				return this.readI64();
			case T.DOUBLE:
				return this.readDbl();
			case T.STRING:
			case T.UTF8:
			case T.UTF16:
				return f.key === "ciphertext" ? this.readStr(true) : this.readStr();
			case T.STRUCT:
				return this.readStruct(f.schema || []);
			case T.LIST:
			case T.SET: {
				const etWire = this.readByte() as TCode;
				const count = this.readI32();
				if (f.elemType !== undefined && f.elemType !== etWire) {
					throw new Error(`Element type mismatch: schema=${f.elemType}, wire=${etWire} (field=${f.key})`);
				}
				const arr: any[] = [];
				for (let i = 0; i < count; i++) {
					if (etWire === T.STRUCT) {
						const elemSchema = f.elemSchema;
						arr.push(elemSchema ? this.readStruct(elemSchema) : this.readStructGeneric());
					} else {
						arr.push(this.readValue(etWire, { id: 0, key: "", type: etWire }));
					}
				}
				return arr;
			}
			case T.MAP: {
				const ktWire = this.readByte() as TCode;
				const vtWire = this.readByte() as TCode;
				const count = this.readI32();
				if (f.keyType !== undefined && f.keyType !== ktWire) {
					throw new Error(`Map key type mismatch: schema=${f.keyType}, wire=${ktWire} (field=${f.key})`);
				}
				if (f.valType !== undefined && f.valType !== vtWire) {
					throw new Error(`Map val type mismatch: schema=${f.valType}, wire=${vtWire} (field=${f.key})`);
				}
				const obj: any = {};
				for (let i = 0; i < count; i++) {
					const k = ktWire === T.STRUCT ? this.readStruct(f.keySchema || []) : this.readValue(ktWire, { id: 0, key: "", type: ktWire });
					const v = vtWire === T.STRUCT ? this.readStruct(f.valSchema || []) : this.readValue(vtWire, { id: 0, key: "", type: vtWire });
					obj[k as any] = v;
				}
				return obj;
			}
			default:
				throw new Error(`Unsupported type ${t}`);
		}
	}

	private skip(t: TCode) {
		this.readValue(t, { id: 0, key: "", type: t } as FieldSchema);
	}

	readStructGeneric(): any {
		const out: any = {};
		while (true) {
			const t = this.readByte();
			if (t === T.STOP) break;
			const id = this.readI16();
			const val = this.readValueAny(t as TCode);
			out[id] = val;
		}
		return out;
	}

	private readValueAny(t: TCode): any {
		switch (t) {
			case T.BOOL:
				return this.readByte() !== 0;
			case T.BYTE:
				return this.readI8();
			case T.I16:
				return this.readI16();
			case T.I32:
				return this.readI32();
			case T.I64:
				return this.readI64();
			case T.DOUBLE:
				return this.readDbl();
			case T.STRING:
			case T.UTF8:
			case T.UTF16:
				return this.readStr();
			case T.STRUCT:
				return this.readStructGeneric();
			case T.LIST:
			case T.SET: {
				const et = this.readByte() as TCode;
				const count = this.readI32();
				const arr: any[] = [];
				for (let i = 0; i < count; i++) arr.push(this.readValueAny(et));
				return arr;
			}
			case T.MAP: {
				const kt = this.readByte() as TCode;
				const vt = this.readByte() as TCode;
				const count = this.readI32();
				const obj: any = {};
				for (let i = 0; i < count; i++) {
					const k = this.readValueAny(kt);
					const v = this.readValueAny(vt);
					obj[k as any] = v;
				}
				return obj;
			}
			default:
				throw new Error(`Unsupported type ${t}`);
		}
	}
}

// Schemas remain unchanged (omitted here for brevity in diff)
// ... paste your existing schema definitions below ...

export async function websocketReverseEngineering(
	token: string,
	decryptionKeyB64: string,
	log?: (message: string) => void,
	onEvent?: (event: any) => void,
	onDecrypted?: (text: string) => void,
) {
	await sodium.ready;

	let decryptionKey: Uint8Array;
	try {
		decryptionKey = base64ToUint8Array(decryptionKeyB64);
	} catch (err) {
		throw new Error("Invalid decryption key. Please provide a base64-encoded key.");
	}

	const socket = new WebSocket(`wss://chat-ws.x.com/ws?token=${encodeURIComponent(token)}`);

	socket.onopen = () => {
		log?.("Socket opened");
		console.info("Socket opened");
	};

	socket.onmessage = async (m) => {
		const mdata = m.data as Blob;
		const d = await mdata.arrayBuffer();


		const decoder = new Decoder(new Uint8Array(d));
		const event = decoder.readStruct(xchatRootSchema)?.event;

		if (!event) {
			const genericDecoder = new Decoder(new Uint8Array(d));

			onEvent?.(genericDecoder.readStructGeneric());
		} else {
			onEvent?.(event);
		}

		const eventJson = JSON.stringify(event, jsonReplacer, 2);

		console.log(event, utf8ToBase64(eventJson));

		const ciphertext = event.payload?.encryptedMessage?.ciphertext as Uint8Array | undefined;
		if (!ciphertext) {
			onDecrypted?.("No ciphertext present in payload.");
			return;
		}

		let decrypted: Uint8Array | null;
		try {
			decrypted = decryptAfterExtractingNonce(ciphertext, decryptionKey);
		} catch (err) {
			const msg = err instanceof Error ? err.message : "Failed to decrypt ciphertext";
			console.warn(msg);
			log?.(msg);
			onDecrypted?.(msg);
			return;
		}

		if (!decrypted) {
			const msg = "Failed to decrypt ciphertext";
			console.warn(msg);
			log?.(msg);
			onDecrypted?.(msg);
			return;
		}

		const cleaned = trimBinaryEdges(decrypted);
		const dc = new TextDecoder();
		const decoded = dc.decode(cleaned);

		console.log("Decrypted", decoded, {
			prefix: Array.from(decrypted.slice(0, 6)),
			suffix: Array.from(decrypted.slice(-6)),
		});
		log?.("Decrypted payload logged to console");
		onDecrypted?.(decoded);

		decoder.readMessageEndStrict?.();
	};

	socket.onerror = (err) => {
		console.error("Socket Error: ", err);
		log?.("Socket error; see console for details.");
	};
	socket.onclose = () => {
		console.info("Socket closed");
		log?.("Socket closed");
	};

	return socket;
}


// Nested payload: encrypted message content (7.1)
export const xchatEncryptedPayloadSchema: FieldSchema[] = [
	// 100: opaque libsodium ciphertext bundle (nonce+ct+mac or similar)
	{ id: 100, key: "ciphertext", type: T.STRING },
	// 101: when this session/key epoch was created (ms, as string)
	{ id: 101, key: "sessionKeyCreatedAtMs", type: T.STRING },
	// 102: key active/usable flag
	{ id: 102, key: "keyActive", type: T.BOOL },
	// 104: message send/encrypt timestamp (ms, as string)
	{ id: 104, key: "sentAtMs", type: T.STRING },
	// 105: key revoked/deleted flag (false in your samples)
	{ id: 105, key: "keyRevoked", type: T.BOOL },
	// 106: encryption format / version (1 in your samples)
	{ id: 106, key: "encVersion", type: T.I32 },
];

// Nested payload: conversation/member state (7.6)
export const xchatConvStateSchema: FieldSchema[] = [
	// 1: conversation id again ("g" + snowflake)
	{ id: 1, key: "conversationId", type: T.STRING },
];

// Nested payload: delivery/read receipt (7.12)
export const xchatReceiptSchema: FieldSchema[] = [
	// 1: message id being acknowledged
	{ id: 1, key: "messageId", type: T.STRING },
	// 2: ack timestamp (ms, as string)
	{ id: 2, key: "ackTimestampMs", type: T.STRING },
];

export const xchatEncryptedKeyPayload: FieldSchema[] = [
	{ id: 1, key: "userId", type: T.STRING },
	{ id: 2, key: "keyB64", type: T.STRING },
	{ id: 3, key: "keyCreatedAtMs", type: T.STRING }
];

export const xchatEncryptedConversationKeySchema: FieldSchema[] = [
	{ id: 1, key: "eventTimestamp", type: T.STRING },

	{ id: 2, key: "encryptedKeyPayload", type: T.LIST, elemType: T.STRUCT, elemSchema: xchatEncryptedKeyPayload },
];

// Field 7: union-ish payload wrapper (only one of these is present at a time)
export const xchatPayloadSchema: FieldSchema[] = [
	{
		id: 1,
		key: "encryptedMessage",
		type: T.STRUCT,
		schema: xchatEncryptedPayloadSchema,
	},

	{
		id: 3,
		key: "encryptedConversationKey",
		type: T.STRUCT,
		schema: xchatEncryptedConversationKeySchema
	},

	{
		id: 6,
		key: "conversationState",
		type: T.STRUCT,
		schema: xchatConvStateSchema,
	},
	{
		id: 12,
		key: "receipt",
		type: T.STRUCT,
		schema: xchatReceiptSchema,
	},
];

// Field 9: actor’s public key bundle
export const xchatKeyBundleSchema: FieldSchema[] = [
	// 1: libsodium public key blob (64 bytes when base64-decoded)
	{ id: 1, key: "sodiumKeyBlobB64", type: T.STRING },
	// 2: key bundle created/valid-since timestamp (ms, as string)
	{ id: 2, key: "keyCreatedAtMs", type: T.STRING },
	// 3: key/protocol version ("3" in all samples)
	{ id: 3, key: "keyVersion", type: T.STRING },
	// 4: P-256 EC public key in SPKI (base64)
	{ id: 4, key: "ecP256SpkiB64", type: T.STRING },
];

// Inner event struct: the thing under top-level "1"
export const xchatEventSchema: FieldSchema[] = [
	// 1: primary event/message id (snowflake, optional on some state-only events)
	{ id: 1, key: "eventId", type: T.STRING },

	// 2: random UUID for this event
	{ id: 2, key: "eventUuid", type: T.STRING },

	// 3: actor user id (sender / participant)
	{ id: 3, key: "actorUserId", type: T.STRING },

	// 4: conversation id ("g" + conversation snowflake)
	{ id: 4, key: "conversationId", type: T.STRING },

	// 5: opaque conversation crypto context blob (constant per conv)
	{ id: 5, key: "cryptoContext", type: T.STRING },

	// 6: event timestamp in ms (string form in all your examples)
	{ id: 6, key: "eventTimestampMs", type: T.STRING },

	// 7: union payload (encrypted msg / receipt / conv state)
	{
		id: 7,
		key: "payload",
		type: T.STRUCT,
		schema: xchatPayloadSchema,
	},

	// 8: event kind enum (0 in all samples so far)
	{ id: 8, key: "eventKind", type: T.I32 },

	// 9: actor’s public key bundle (present on message + receipt events)
	{
		id: 9,
		key: "keyBundle",
		type: T.STRUCT,
		schema: xchatKeyBundleSchema,
	},

	// 10: pointer to previous event in this conversation (snowflake id)
	{ id: 10, key: "prevEventId", type: T.STRING },

	// 11: extra boolean flag (seen in earlier DM examples only)
	{ id: 11, key: "flag11", type: T.BOOL },
];

// Outer wrapper: what your generic decode shows as top-level { "1": { ... } }
export const xchatRootSchema: FieldSchema[] = [
	{
		id: 1,
		key: "event",
		type: T.STRUCT,
		schema: xchatEventSchema,
	},
];
