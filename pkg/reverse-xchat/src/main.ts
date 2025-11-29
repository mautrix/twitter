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

        <label class="block space-y-2">
          <span class="text-sm font-medium text-slate-200">Decryption key (base64)</span>
          <input
            id="key-input"
            name="key"
            type="text"
            value="${defaultKeyB64}"
            required
            class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
          />
        </label>

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

      <section class="grid gap-4 lg:grid-cols-2">
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
    </div>
  </main>
`;

const form = document.querySelector<HTMLFormElement>('#connect-form');
const tokenInput = document.querySelector<HTMLTextAreaElement>('#token-input');
const keyInput = document.querySelector<HTMLInputElement>('#key-input');
const statusLog = document.querySelector<HTMLPreElement>('#status-log');
const connectButton = document.querySelector<HTMLButtonElement>('#connect-btn');
const eventOutput = document.querySelector<HTMLPreElement>('#event-output');
const decryptedOutput = document.querySelector<HTMLPreElement>('#decrypted-output');

if (!form || !tokenInput || !keyInput || !statusLog || !connectButton || !eventOutput || !decryptedOutput) {
	throw new Error("Failed to initialize app controls");
}

const logStatus = (message: string) => {
	const timestamp = new Date().toLocaleTimeString();
	statusLog.textContent = `[${timestamp}] ${message}`;
	console.info(message);
};

const jsonReplacer = (_key: string, value: unknown) => (typeof value === "bigint" ? value.toString() : value);
const renderEventJson = (json: string) => {
	eventOutput.textContent = json;
};
const renderDecryptedPayload = (text: string) => {
	decryptedOutput.textContent = text;
};

let activeSocket: WebSocket | null = null;

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

function base64ToUint8Array(b64: string) {
	// Normalize base64url
	b64 = b64.replace(/-/g, "+").replace(/_/g, "/");
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
export enum T {
	STOP = 0,
	BOOL = 2,
	BYTE = 3,
	DOUBLE = 4,
	I16 = 6,
	I32 = 8,
	I64 = 10,
	STRING = 11,
	STRUCT = 12,
	MAP = 13,
	SET = 14,
	LIST = 15,
	UTF8 = 16,
	UTF16 = 17,
}

export type FieldSchema = {
	id: number;
	key: string;
	type: T;
	schema?: FieldSchema[];
	elemType?: T;
	elemSchema?: FieldSchema[];
	keyType?: T;
	keySchema?: FieldSchema[];
	valType?: T;
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

const field = (type: T, id: number, valBytes: Uint8Array) =>
	concat(Uint8Array.of(type), enc.i16(id), valBytes);

const stop = () => Uint8Array.of(T.STOP);

const defaultEncode = (t: T) => (v: any): Uint8Array => {
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

export enum MsgType {
	CALL = 1,
	REPLY = 2,
	EXCEPTION = 3,
	ONEWAY = 4,
}

export type MessageBegin = {
	name: string;
	type: MsgType;
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
		const type = (word & 0x000000ff) as MsgType;
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
				this.skip(t as T);
				continue;
			}
			out[f.key] = this.readValue(t as T, f);
		}
		return out;
	}

	private readValue(t: T, f: FieldSchema): any {
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
				const etWire = this.readByte() as T;
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
				const ktWire = this.readByte() as T;
				const vtWire = this.readByte() as T;
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

	private skip(t: T) {
		this.readValue(t, { id: 0, key: "", type: t } as FieldSchema);
	}

	readStructGeneric(): any {
		const out: any = {};
		while (true) {
			const t = this.readByte();
			if (t === T.STOP) break;
			const id = this.readI16();
			const val = this.readValueAny(t as T);
			out[id] = val;
		}
		return out;
	}

	private readValueAny(t: T): any {
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
				const et = this.readByte() as T;
				const count = this.readI32();
				const arr: any[] = [];
				for (let i = 0; i < count; i++) arr.push(this.readValueAny(et));
				return arr;
			}
			case T.MAP: {
				const kt = this.readByte() as T;
				const vt = this.readByte() as T;
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
		const { event } = decoder.readStruct(xchatRootSchema);
		onEvent?.(event);

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

// Field 7: union-ish payload wrapper (only one of these is present at a time)
export const xchatPayloadSchema: FieldSchema[] = [
	{
		id: 1,
		key: "encryptedMessage",
		type: T.STRUCT,
		schema: xchatEncryptedPayloadSchema,
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
