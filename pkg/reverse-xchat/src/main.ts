import "./style.css";
import "highlight.js/styles/github-dark-dimmed.css";
import hljs from "highlight.js/lib/core";
import jsonLang from "highlight.js/lib/languages/json";
import { Client, Configuration, RecoverErrorReason } from "juicebox-sdk";
import sodium from "libsodium-wrappers";

hljs.registerLanguage("json", jsonLang);

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
	throw new Error("App container missing");
}

type UiRefs = {
	form: HTMLFormElement;
	tokenInput: HTMLTextAreaElement;
	keyInput: HTMLInputElement;
	pubKeyInput: HTMLInputElement;
	keyPreview: HTMLParagraphElement;
	maskKeyToggle: HTMLInputElement;
	statusLog: HTMLPreElement;
	connectButton: HTMLButtonElement;
	keySelect: HTMLSelectElement;
	saveKeyButton: HTMLButtonElement;
	deleteKeyButton: HTMLButtonElement;
	eventOutput: HTMLPreElement;
	decryptedOutput: HTMLPreElement;
	manualForm: HTMLFormElement;
	manualInput: HTMLTextAreaElement;
	manualSchemaSelect: HTMLSelectElement;
	manualOutput: HTMLPreElement;
	manualVerifyBtn: HTMLButtonElement;
	manualVerifyOutput: HTMLPreElement;
	genericForm: HTMLFormElement;
	genericInput: HTMLTextAreaElement;
	genericOutput: HTMLPreElement;
	genericVerifyBtn: HTMLButtonElement;
	genericVerifyOutput: HTMLPreElement;
	convEventForm: HTMLFormElement;
	convEventInput: HTMLTextAreaElement;
	convEventOutput: HTMLPreElement;
	juiceboxForm: HTMLFormElement;
	juiceboxJsonInput: HTMLTextAreaElement;
	juiceboxPinInput: HTMLInputElement;
	juiceboxOutput: HTMLPreElement;
	juiceboxSubmitBtn: HTMLButtonElement;
	base64BlobForm: HTMLFormElement;
	base64BlobInput: HTMLTextAreaElement;
	base64BlobOutput: HTMLPreElement;
};

const getUiRefs = (): UiRefs => {
	const refs: Partial<UiRefs> = {
		form: document.querySelector<HTMLFormElement>("#connect-form") ?? undefined,
		tokenInput: document.querySelector<HTMLTextAreaElement>("#token-input") ?? undefined,
		keyInput: document.querySelector<HTMLInputElement>("#key-input") ?? undefined,
		pubKeyInput: document.querySelector<HTMLInputElement>("#pubkey-input") ?? undefined,
		keyPreview: document.querySelector<HTMLParagraphElement>("#key-preview") ?? undefined,
		maskKeyToggle: document.querySelector<HTMLInputElement>("#mask-key-toggle") ?? undefined,
		statusLog: document.querySelector<HTMLPreElement>("#status-log") ?? undefined,
		connectButton: document.querySelector<HTMLButtonElement>("#connect-btn") ?? undefined,
		keySelect: document.querySelector<HTMLSelectElement>("#key-select") ?? undefined,
		saveKeyButton: document.querySelector<HTMLButtonElement>("#save-key-btn") ?? undefined,
		deleteKeyButton: document.querySelector<HTMLButtonElement>("#delete-key-btn") ?? undefined,
		eventOutput: document.querySelector<HTMLPreElement>("#event-output") ?? undefined,
		decryptedOutput: document.querySelector<HTMLPreElement>("#decrypted-output") ?? undefined,
		manualForm: document.querySelector<HTMLFormElement>("#manual-form") ?? undefined,
		manualInput: document.querySelector<HTMLTextAreaElement>("#manual-input") ?? undefined,
		manualSchemaSelect: document.querySelector<HTMLSelectElement>("#manual-schema-select") ?? undefined,
		manualOutput: document.querySelector<HTMLPreElement>("#manual-output") ?? undefined,
		manualVerifyBtn: document.querySelector<HTMLButtonElement>("#manual-verify-btn") ?? undefined,
		manualVerifyOutput:
			document.querySelector<HTMLPreElement>("#manual-verify-output") ?? undefined,
		genericForm: document.querySelector<HTMLFormElement>("#generic-form") ?? undefined,
		genericInput: document.querySelector<HTMLTextAreaElement>("#generic-input") ?? undefined,
		genericOutput: document.querySelector<HTMLPreElement>("#generic-output") ?? undefined,
		genericVerifyBtn: document.querySelector<HTMLButtonElement>("#generic-verify-btn") ?? undefined,
		genericVerifyOutput:
			document.querySelector<HTMLPreElement>("#generic-verify-output") ?? undefined,
		convEventForm: document.querySelector<HTMLFormElement>("#conv-event-form") ?? undefined,
		convEventInput:
			document.querySelector<HTMLTextAreaElement>("#conv-event-input") ?? undefined,
		convEventOutput: document.querySelector<HTMLPreElement>("#conv-event-output") ?? undefined,
		juiceboxForm: document.querySelector<HTMLFormElement>("#juicebox-form") ?? undefined,
		juiceboxJsonInput:
			document.querySelector<HTMLTextAreaElement>("#juicebox-json-input") ?? undefined,
		juiceboxPinInput:
			document.querySelector<HTMLInputElement>("#juicebox-pin-input") ?? undefined,
		juiceboxOutput: document.querySelector<HTMLPreElement>("#juicebox-output") ?? undefined,
		juiceboxSubmitBtn:
			document.querySelector<HTMLButtonElement>("#juicebox-submit-btn") ?? undefined,
		base64BlobForm: document.querySelector<HTMLFormElement>("#blob-form") ?? undefined,
		base64BlobInput: document.querySelector<HTMLTextAreaElement>("#blob-input") ?? undefined,
		base64BlobOutput: document.querySelector<HTMLPreElement>("#blob-output") ?? undefined,
	};

	const missing = Object.entries(refs)
		.filter(([, el]) => !el)
		.map(([name]) => name);

	if (missing.length) {
		throw new Error(`Failed to initialize app controls: ${missing.join(", ")}`);
	}

	return refs as UiRefs;
};

app.innerHTML = `
  <main class="min-h-screen bg-slate-950 text-slate-50 flex items-center justify-center px-4 py-10 w-full">
    <div class="w-full max-w-3xl space-y-6">
      <div class="space-y-1">
        <h1 class="text-3xl font-semibold">Reverse XChat</h1>
        <p class="text-slate-300">Utilities to recover keys, decrypt events, and inspect blobs.</p>
      </div>
      <div class="flex items-center justify-end">
        <label class="inline-flex items-center gap-2 text-xs text-slate-300 rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2">
          <input id="mask-key-toggle" type="checkbox" class="h-4 w-4 rounded border-slate-600 bg-slate-900/80 text-sky-500 focus:ring-sky-500" />
          <span>Hide private key everywhere</span>
        </label>
      </div>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full space-y-3">
        <div class="space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-amber-300">Juicebox</p>
          <h2 class="text-lg font-semibold text-slate-100">Import key from GetPublicKeysResult</h2>
          <p class="text-sm text-slate-300">Paste the GraphQL result (or base64 of it), enter your PIN, and we will recover the secret key with juicebox-sdk and stash it in your saved keys.</p>
        </div>
        <form id="juicebox-form" class="space-y-3">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">GetPublicKeysResult JSON</span>
            <textarea
              id="juicebox-json-input"
              name="juicebox-json"
              rows="7"
              placeholder="Paste the JSON (or base64-encoded JSON) here"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500/40"
            ></textarea>
          </label>
          <div class="grid gap-3 md:grid-cols-2">
            <label class="block space-y-1">
              <span class="text-xs font-medium text-slate-300">PIN (4 digits)</span>
              <input
                id="juicebox-pin-input"
                name="juicebox-pin"
                type="password"
                inputmode="numeric"
                pattern="\\d{4}"
                maxlength="4"
                placeholder="1234"
                class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500/40"
              />
            </label>
          </div>
          <div class="flex items-center gap-3">
            <button
              id="juicebox-submit-btn"
              type="submit"
              class="inline-flex items-center justify-center rounded-lg bg-amber-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Recover & Save Key
            </button>
            <p class="text-xs text-slate-400">We extract key_store_token_map_json + token_map → juicebox → saved key.</p>
          </div>
        </form>
        <pre id="juicebox-output" class="json-output h-44 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Awaiting Juicebox input...</pre>
      </section>

      <form id="connect-form" class="space-y-4 rounded-xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-sky-900/30 backdrop-blur">
        <div class="space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">XChat websocket</p>
          <p class="text-sm text-slate-300">Provide your token and decryption key, then click connect to start the websocket.</p>
        </div>
        <label class="block space-y-2">
          <span class="text-sm font-medium text-slate-200">Token</span>
          <textarea id="token-input" name="token" rows="4" required class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"></textarea>
        </label>

        <div class="space-y-2">
          <div class="flex items-center justify-between gap-3">
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
            value=""
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
          <pre id="event-output" class="json-output h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">No event received yet.</pre>
        </div>
        <div class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20">
          <div class="mb-2 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-slate-100">Decrypted payload</h2>
            <span class="rounded-full bg-emerald-800/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-emerald-100">Plaintext</span>
          </div>
          <pre id="decrypted-output" class="json-output h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">No decrypted payload yet.</pre>
        </div>
      </section>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full">
        <div class="mb-3 space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-400">Manual decode</p>
          <h2 class="text-lg font-semibold text-slate-100">Paste a base64 payload</h2>
          <p class="text-sm text-slate-300">Base64-decodes the payload and parses it with a selected Thrift schema. Use the root schema for websocket binaries; use the event schema for direct message events. Generic is available as a fallback.</p>
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
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">Decode using schema</span>
            <select
              id="manual-schema-select"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500/40"
            >
              <option value="auto">Auto (Root → Event → Generic)</option>
              <option value="root">Root schema (built-in)</option>
              <option value="event">Event schema (built-in messageEventSchema)</option>
              <option value="generic">Generic reader only</option>
            </select>
          </label>
          <p class="text-xs text-slate-400">Root wraps the websocket message; event is the message event payload. Generic will attempt field IDs only.</p>
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
        <pre id="manual-output" class="json-output mt-4 h-48 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Awaiting ciphertext...</pre>
        <div class="mt-3 flex flex-col gap-2">
          <div class="flex items-center gap-3">
            <button
              id="manual-verify-btn"
              type="button"
              class="inline-flex items-center justify-center rounded-lg bg-amber-600 px-3 py-2 text-xs font-semibold text-white shadow hover:bg-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Verify signature
            </button>
            <p class="text-xs text-slate-400">Paste a public key above, decode an event, then verify.</p>
          </div>
          <pre id="manual-verify-output" class="h-24 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Decode an event above to enable verification.</pre>
        </div>
      </section>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full">
        <div class="mb-3 space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">Generic decode</p>
          <h2 class="text-lg font-semibold text-slate-100">Paste a base64 payload (generic)</h2>
          <p class="text-sm text-slate-300">Base64-decodes the payload and parses it with the generic struct parser into JSON using <code>readStructGeneric</code> (no schema assumed).</p>
        </div>
        <form id="generic-form" class="space-y-3">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">Ciphertext bundle (base64)</span>
            <textarea
              id="generic-input"
              name="generic-ciphertext"
              rows="3"
              placeholder="Paste base64-encoded ciphertext bundle here"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-cyan-500 focus:outline-none focus:ring-2 focus:ring-cyan-500/40"
            ></textarea>
          </label>
          <div class="flex items-center gap-3">
            <button
              id="generic-decrypt-btn"
              type="submit"
              class="inline-flex items-center justify-center rounded-lg bg-cyan-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-cyan-500 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Parse message (generic)
            </button>
            <p class="text-xs text-slate-400">We will base64-decode and parse with the generic reader below.</p>
          </div>
        </form>
        <pre id="generic-output" class="json-output mt-4 h-48 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Awaiting ciphertext...</pre>
        <div class="mt-3 flex flex-col gap-2">
          <div class="flex items-center gap-3">
            <button
              id="generic-verify-btn"
              type="button"
              class="inline-flex items-center justify-center rounded-lg bg-amber-600 px-3 py-2 text-xs font-semibold text-white shadow hover:bg-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Verify signature
            </button>
            <p class="text-xs text-slate-400">Paste a public key above, decode an event, then verify.</p>
          </div>
          <pre id="generic-verify-output" class="h-24 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Decode an event above to enable verification.</pre>
        </div>
      </section>

      <section class="rounded-xl border border-slate-800 bg-slate-900/70 p-4 shadow-sm shadow-sky-900/20 w-full">
        <div class="mb-3 space-y-1">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-indigo-300">Blob decrypt</p>
          <h2 class="text-lg font-semibold text-slate-100">Decrypt a base64 blob</h2>
          <p class="text-sm text-slate-300">UI only: paste any base64 blob and click decrypt. Hook your decryption logic in the handler later.</p>
        </div>
        <form id="blob-form" class="space-y-3">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-slate-200">Base64 blob</span>
            <textarea
              id="blob-input"
              name="blob"
              rows="3"
              placeholder="Paste the base64 blob here"
              class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/40"
            ></textarea>
          </label>
          <div class="flex items-center gap-3">
            <button
              id="blob-submit-btn"
              type="submit"
              class="inline-flex items-center justify-center rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-slate-950 disabled:opacity-50"
            >
              Decrypt blob
            </button>
            <p class="text-xs text-slate-400">Outputs a stub message; add your decryption later.</p>
          </div>
        </form>
        <pre id="blob-output" class="mt-4 h-40 overflow-auto whitespace-pre-wrap rounded-lg bg-slate-950/70 p-3 text-xs text-slate-100 font-mono leading-relaxed">Awaiting blob...</pre>
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

const {
	form,
	tokenInput,
	keyInput,
	pubKeyInput,
	keyPreview,
	maskKeyToggle,
	statusLog,
	connectButton,
	keySelect,
	saveKeyButton,
	deleteKeyButton,
	eventOutput,
	decryptedOutput,
	manualForm,
	manualInput,
	manualSchemaSelect,
	manualOutput,
	manualVerifyBtn,
	manualVerifyOutput,
	genericForm,
	genericInput,
	genericOutput,
	genericVerifyBtn,
	genericVerifyOutput,
	convEventForm,
	convEventInput,
	convEventOutput,
	juiceboxForm,
	juiceboxJsonInput,
	juiceboxPinInput,
	juiceboxOutput,
	juiceboxSubmitBtn,
	base64BlobForm,
	base64BlobInput,
	base64BlobOutput,
} = getUiRefs();

const logStatus = (message: string) => {
	const timestamp = new Date().toLocaleTimeString();
	statusLog.textContent = `[${timestamp}] ${message}`;
	console.info(message);
};
// ---- Persistence helpers ----
const STORAGE_KEYS = "xchatSavedKeys";

type SavedKey = {
	priv: string;
	pub?: string;
};
type DecodedPayload = { parsed: any; bytes: Uint8Array };
type AppState = {
	savedKeys: SavedKey[];
	hidePrivateKeys: boolean;
	lastJuiceboxPayload: string;
	lastManualDecoded: DecodedPayload | null;
	lastGenericDecoded: DecodedPayload | null;
};

declare global {
	interface Window {
		JuiceboxGetAuthToken?: (realmId: Uint8Array) => Promise<string> | string;
	}
}
// ---- App state ----

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
					return {
						priv: String(item.priv),
						pub: item.pub ? String(item.pub) : undefined,
					};
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

const state: AppState = {
	savedKeys: (() => {
		const stored = loadSavedKeys();
		return stored.length ? stored : [];
	})(),
	hidePrivateKeys: false,
	lastJuiceboxPayload: "",
	lastManualDecoded: null,
	lastGenericDecoded: null,
};

// ---- Saved key rendering / management ----
const renderSavedKeyOptions = (activeIndex = 0) => {
	keySelect.innerHTML = "";
	const frag = document.createDocumentFragment();
	const preview = (val: string) =>
		val.length > 14 ? `${val.slice(0, 8)}…${val.slice(-6)}` : val;
	const saved = state.savedKeys;

	if (!saved.length) {
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
		saved.forEach((k, idx) => {
			const opt = document.createElement("option");
			opt.value = String(idx);
			if (state.hidePrivateKeys) {
				opt.textContent = k.pub
					? `Pub ${preview(k.pub)} | Key hidden`
					: `Key ${idx + 1} (hidden)`;
			} else {
				opt.textContent = k.pub
					? `Priv ${preview(k.priv)} | Pub ${preview(k.pub)}`
					: `Priv ${preview(k.priv)}`;
			}
			frag.appendChild(opt);
		});
		const safeIndex = activeIndex < saved.length ? activeIndex : 0;
		keySelect.value = String(safeIndex);
		const activeKey = saved[safeIndex];
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
	const idx = state.savedKeys.findIndex((k) => k.priv === entry.priv);
	if (idx >= 0) {
		state.savedKeys[idx] = entry;
		renderSavedKeyOptions(idx);
	} else {
		state.savedKeys = [entry, ...state.savedKeys];
		renderSavedKeyOptions(0);
	}
	persistKeys(state.savedKeys);
};

renderSavedKeyOptions(0);

// ---- JSON helpers ----
const jsonReplacer = (_key: string, value: unknown) =>
	typeof value === "bigint" ? value.toString() : value;
const tryParseJson = (text: string) => {
	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
};
const tryParseJsonOrBase64Json = (text: string) => {
	const parsed = tryParseJson(text);
	if (parsed !== null) return parsed;
	try {
		const decoded = new TextDecoder().decode(base64ToUint8Array(text));
		return tryParseJson(decoded);
	} catch {
		return null;
	}
};
// ---- Juicebox helpers ----
const extractJuiceboxBundle = (root: any) => {
	const pkList = root?.public_keys_with_token_map;
	if (!Array.isArray(pkList) || !pkList.length) {
		throw new Error("public_keys_with_token_map missing or empty");
	}
	const tokenMapNode = pkList[0]?.token_map ?? pkList[0];
	const rawConfig =
		tokenMapNode?.key_store_token_map_json ??
		tokenMapNode?.key_store_token_map_json_string ??
		tokenMapNode?.key_store_token_map_json_json;
	let configObj: any = null;
	if (typeof rawConfig === "string") {
		configObj = tryParseJson(rawConfig);
	}
	if (!configObj && rawConfig && typeof rawConfig === "object") {
		configObj = rawConfig;
	}
	if (!configObj) {
		throw new Error("key_store_token_map_json missing or invalid");
	}

	const tokenEntries = tokenMapNode?.token_map;
	const tokenMap: Record<string, string> = {};
	if (Array.isArray(tokenEntries)) {
		for (const entry of tokenEntries) {
			const k = entry?.key;
			const token = entry?.value?.token;
			if (typeof k === "string" && typeof token === "string") {
				tokenMap[k.toLowerCase()] = token;
			}
		}
	}
	if (!Object.keys(tokenMap).length) {
		throw new Error("token_map missing or empty");
	}

	return { configObj, tokenMap };
};
const extractPublicKeyVersion = (root: any) => {
	const meta = root?.public_keys_with_token_map?.[0];
	const version = meta?.public_key_with_metadata?.version;
	console.log(root);
	if (typeof version === "string") {
		const trimmed = version.trim();
		return trimmed || null;
	}
	if (typeof version === "number" && Number.isFinite(version)) return version;
	return null;
};
// ---- Rendering helpers ----
const highlightJson = (target: HTMLPreElement, text: string) => {
	target.classList.add("hljs");
	target.classList.add("language-json");
	const highlighted = hljs.highlight(text, { language: "json", ignoreIllegals: true }).value;
	target.innerHTML = highlighted;
};
const renderJsonSection = (target: HTMLPreElement, text: string) => {
	if (!text) {
		target.textContent = "";
		return;
	}
	const parsed = tryParseJson(text);
	if (parsed !== null) {
		const pretty = JSON.stringify(parsed, jsonReplacer, 2);
		highlightJson(target, pretty);
		return;
	}
	target.textContent = text;
};
const renderEventJson = (json: string) => renderJsonSection(eventOutput, json);
const renderDecryptedPayload = (text: string) => renderJsonSection(decryptedOutput, text);
const renderManualPayload = (text: string) => renderJsonSection(manualOutput, text);
const renderManualVerifyPayload = (text: string) => {
	manualVerifyOutput.textContent = text;
};
const renderGenericPayload = (text: string) => renderJsonSection(genericOutput, text);
const renderGenericVerifyPayload = (text: string) => {
	genericVerifyOutput.textContent = text;
};
const renderJuiceboxPayload = (text: string) => {
	state.lastJuiceboxPayload = text;
	const parsed = tryParseJson(text);
	if (state.hidePrivateKeys && parsed && typeof parsed === "object") {
		const clone = JSON.parse(JSON.stringify(parsed));
		if ("savedPrivKeyB64" in clone) clone.savedPrivKeyB64 = "[hidden]";
		if ("savedKeyB64" in clone) clone.savedKeyB64 = "[hidden]";
		renderJsonSection(juiceboxOutput, JSON.stringify(clone, jsonReplacer, 2));
		return;
	}
	renderJsonSection(juiceboxOutput, text);
};
const updateKeyPreview = () => {
	const value = keyInput.value.trim();
	if (state.hidePrivateKeys) {
		keyPreview.textContent = value
			? 'Hidden while "Hide private key" is on.'
			: "Length: 0 | hidden";
		return;
	}
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
const renderBlobOutput = (text: string) => {
	base64BlobOutput.textContent = text;
};
const describeRecoverError = (err: unknown) => {
	if (err && typeof err === "object" && "reason" in err) {
		const reason = (err as any).reason as RecoverErrorReason;
		const reasonLabel = RecoverErrorReason[reason] ?? `Reason ${reason}`;
		const guesses = (err as any).guesses_remaining;
		return guesses !== undefined
			? `${reasonLabel} (guesses remaining: ${guesses})`
			: reasonLabel;
	}
	return err instanceof Error ? err.message : String(err);
};

const tryThrift = <T>(fn: () => T): T | null => {
	try {
		return fn();
	} catch {
		return null;
	}
};

const parseConversationEventInput = (raw: string) => {
	const parsed = tryParseJson(raw);
	if (parsed) return parsed;

	const bytes = base64ToUint8Array(raw);

	const decodedStruct = tryThrift(() => new Decoder(bytes).readStruct(xchatRootSchema));
	if (decodedStruct) return decodedStruct;

	const decodedGeneric = tryThrift(() => new Decoder(bytes).readStructGeneric());
	if (decodedGeneric) return decodedGeneric;

	try {
		const decodedText = new TextDecoder().decode(bytes);
		return tryParseJson(decodedText);
	} catch {
		return null;
	}
};

type DecodeSource = "manual" | "generic";
const setDecoded = (source: DecodeSource, parsed: any | null, bytes?: Uint8Array) => {
	const assign =
		source === "manual"
			? (payload: DecodedPayload | null) => {
				state.lastManualDecoded = payload;
			}
			: (payload: DecodedPayload | null) => {
				state.lastGenericDecoded = payload;
			};
	const button = source === "manual" ? manualVerifyBtn : genericVerifyBtn;
	const render =
		source === "manual" ? renderManualVerifyPayload : renderGenericVerifyPayload;

	if (parsed && bytes) {
		assign({ parsed, bytes });
		button.disabled = false;
		render("Decoded. Click verify to check signature.");
	} else {
		assign(null);
		button.disabled = true;
		render("Decode an event above to enable verification.");
	}
};
const setManualDecoded = (parsed: any | null, bytes?: Uint8Array) =>
	setDecoded("manual", parsed, bytes);
const setGenericDecoded = (parsed: any | null, bytes?: Uint8Array) =>
	setDecoded("generic", parsed, bytes);

let activeSocket: WebSocket | null = null;
// ---- Event handlers ----

const handleMaskToggle = () => {
	state.hidePrivateKeys = maskKeyToggle.checked;
	keyInput.type = state.hidePrivateKeys ? "password" : "text";
	const currentIdx = Number(keySelect.value) || 0;
	renderSavedKeyOptions(currentIdx);
	updateKeyPreview();
	if (state.lastJuiceboxPayload) renderJuiceboxPayload(state.lastJuiceboxPayload);
};

const handleKeySelectChange = () => {
	const idx = Number(keySelect.value);
	const entry = state.savedKeys[idx];
	if (entry) {
		keyInput.value = entry.priv;
		pubKeyInput.value = entry.pub ?? "";
	}
	updateKeyPreview();
};

const handlePinInputChange = () => {
	const digits = juiceboxPinInput.value.replace(/\D/g, "").slice(0, 4);
	if (digits !== juiceboxPinInput.value) {
		juiceboxPinInput.value = digits;
	}
};

const handleSaveKeyClick = () => {
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
};

const handleDeleteKeyClick = () => {
	const idx = Number(keySelect.value);
	if (Number.isNaN(idx) || idx < 0 || idx >= state.savedKeys.length) {
		logStatus("No saved key selected to delete.");
		return;
	}
	state.savedKeys.splice(idx, 1);
	persistKeys(state.savedKeys);
	const nextIdx = Math.max(0, Math.min(idx, state.savedKeys.length - 1));
	renderSavedKeyOptions(nextIdx);
	logStatus("Key deleted.");
	updateKeyPreview();
};

const handleConnectSubmit = async (event: Event) => {
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
};

const handleManualSubmit = async (event: Event) => {
	event.preventDefault();

	const ciphertextB64 = manualInput.value.trim();

	if (!ciphertextB64) {
		renderManualPayload(
			"Please paste the base64 payload (same as websocket binary, base64-encoded).",
		);
		setManualDecoded(null);
		return;
	}

	let bytes: Uint8Array;
	try {
		bytes = base64ToUint8Array(ciphertextB64);
	} catch (err) {
		renderManualPayload(`Invalid base64: ${err instanceof Error ? err.message : String(err)}`);
		setManualDecoded(null);
		return;
	}

	const errors: string[] = [];
	const tryDecode = (schema: FieldSchema[], label: string) => {
		try {
			const decoder = new Decoder(bytes);
			const obj = decoder.readStruct(schema);
			const bodyJson = JSON.stringify(obj, jsonReplacer, 2);
			renderManualPayload(bodyJson);
			setManualDecoded(obj, bytes);
			return true;
		} catch (err) {
			const msg = err instanceof Error ? err.message : String(err);
			errors.push(`${label} failed: ${msg}`);
			return false;
		}
	};
	const tryGeneric = () => {
		try {
			const decoder = new Decoder(bytes);
			const obj = decoder.readStructGeneric();
			const bodyJson = JSON.stringify(obj, jsonReplacer, 2);
			renderManualPayload(bodyJson);
			setManualDecoded(obj, bytes);
			return true;
		} catch (err) {
			const msg = err instanceof Error ? err.message : String(err);
			errors.push(`Generic decode failed: ${msg}`);
			return false;
		}
	};

	const choice = manualSchemaSelect.value as "auto" | "root" | "event" | "generic";

	const decodeWith = (schema: FieldSchema[] | null, label: string) =>
		schema ? tryDecode(schema, label) : false;

	switch (choice) {
		case "root":
			if (decodeWith(xchatRootSchema, "Root schema")) return;
			break;
		case "event":
			if (decodeWith(messageEventSchema, "Event schema")) return;
			break;
		case "generic":
			if (tryGeneric()) return;
			break;
		case "auto":
		default:
			if (tryDecode(xchatRootSchema, "Built-in root schema")) return;
			if (tryDecode(messageEventSchema, "Built-in event schema")) return;
			break;
	}

	if (tryGeneric()) return;

	renderManualPayload(`Parse Error:\n${errors.join("\n")}`);
	setManualDecoded(null);
};

const handleGenericSubmit = async (event: Event) => {
	event.preventDefault();

	const ciphertextB64 = genericInput.value.trim();

	if (!ciphertextB64) {
		renderGenericPayload(
			"Please paste the base64 payload (same as websocket binary, base64-encoded).",
		);
		setGenericDecoded(null);
		return;
	}

	let bytes: Uint8Array;
	try {
		bytes = base64ToUint8Array(ciphertextB64);
	} catch (err) {
		renderGenericPayload(`Invalid base64: ${err instanceof Error ? err.message : String(err)}`);
		setGenericDecoded(null);
		return;
	}

	try {
		const decoder = new Decoder(bytes);
		const obj = decoder.readStructGeneric();
		const bodyJson = JSON.stringify(obj, jsonReplacer, 2);
		renderGenericPayload(bodyJson);
		setGenericDecoded(obj, bytes);
	} catch (err) {
		renderGenericPayload(`Parse Error: ${err}`);
		setGenericDecoded(null);
	}
};

const handleJuiceboxSubmit = async (event: Event) => {
	event.preventDefault();

	const raw = juiceboxJsonInput.value.trim();
	const pinDigits = (juiceboxPinInput.value ?? "").replace(/\D/g, "").slice(0, 4);
	const info = "";

	juiceboxPinInput.value = pinDigits;

	if (!raw) {
		renderJuiceboxPayload("Please paste the GetPublicKeysResult JSON or base64-encoded JSON.");
		return;
	}

	if (!/^\d{4}$/.test(pinDigits)) {
		renderJuiceboxPayload("Enter your 4-digit PIN to recover the secret.");
		return;
	}

	const parsed = tryParseJsonOrBase64Json(raw);
	if (!parsed) {
		renderJuiceboxPayload("Could not parse JSON (neither plain nor base64).");
		return;
	}

	const publicKeyVersion = extractPublicKeyVersion(parsed);

	let bundle: { configObj: any; tokenMap: Record<string, string> };
	try {
		bundle = extractJuiceboxBundle(parsed);
	} catch (err) {
		renderJuiceboxPayload(`Parse error: ${err instanceof Error ? err.message : String(err)}`);
		return;
	}

	try {
		juiceboxSubmitBtn.disabled = true;
		renderJuiceboxPayload("Recovering key via Juicebox…");

		const configuration = new Configuration(bundle.configObj);

		window.JuiceboxGetAuthToken = async (realmId: Uint8Array) => {
			const realmIdHex = bytesToHex(realmId).toLowerCase();
			const token = bundle.tokenMap[realmIdHex];
			if (!token) {
				throw new Error(`No token found for realm ${realmIdHex}`);
			}
			return token;
		};

		const client = new Client(configuration, []);
		const encoder = new TextEncoder();
		const recovered = new Uint8Array(
			await client.recover(encoder.encode(pinDigits), encoder.encode(info)),
		);

		let privBytes: Uint8Array;
		let pubBytes: Uint8Array | undefined;

		if (recovered.length >= 64 && recovered.length % 2 === 0) {
			const half = recovered.length / 2;
			privBytes = recovered.slice(0, half);
			pubBytes = recovered.slice(half);
		} else if (recovered.length > 32) {
			privBytes = recovered.slice(0, 32);
			pubBytes = recovered.slice(32);
		} else {
			privBytes = recovered;
		}

		const privKeyB64 = bytesToBase64(privBytes);
		const pubKeyB64 = pubBytes ? bytesToBase64(pubBytes) : undefined;

		upsertKey(privKeyB64, pubKeyB64);
		updateKeyPreview();

		renderJuiceboxPayload(
			JSON.stringify(
				{
					publicKeyVersion: publicKeyVersion ?? "(not found)",
					savedPrivKeyB64: privKeyB64,
					savedPubKeyB64: pubKeyB64,
					config: bundle.configObj,
					tokenMap: bundle.tokenMap,
				},
				jsonReplacer,
				2,
			),
		);
		logStatus("Recovered keypair from Juicebox and saved to your saved keys.");
	} catch (err) {
		renderJuiceboxPayload(`Juicebox recovery failed: ${describeRecoverError(err)}`);
	} finally {
		juiceboxSubmitBtn.disabled = false;
	}
};

const getDecodedState = (source: DecodeSource) =>
	source === "manual" ? state.lastManualDecoded : state.lastGenericDecoded;

const handleVerifyClick = async (source: DecodeSource) => {
	await sodium.ready;

	const render =
		source === "manual" ? renderManualVerifyPayload : renderGenericVerifyPayload;
	const button = source === "manual" ? manualVerifyBtn : genericVerifyBtn;
	const last = getDecodedState(source);
	if (!last) {
		render("Decode an event above before verifying.");
		return;
	}

	let evt = (last.parsed as any)?.event ?? last.parsed;
	if ((!evt || typeof evt !== "object") && last.bytes) {
		const decodedWithSchema = tryThrift(() => new Decoder(last.bytes).readStruct(xchatRootSchema));
		if (decodedWithSchema && typeof decodedWithSchema === "object") {
			evt = (decodedWithSchema as any)?.event ?? decodedWithSchema;
		}
	}
	if (!evt || typeof evt !== "object") {
		render("Decoded payload does not look like an event object.");
		return;
	}

	const signatureB64 = evt?.keyBundle?.sodiumKeyBlobB64;
	if (typeof signatureB64 !== "string") {
		render("Missing signature at keyBundle.sodiumKeyBlobB64.");
		return;
	}

	const publicKeyB64FromInput = pubKeyInput.value.trim();
	const eventPubKeyB64 =
		typeof evt?.keyBundle?.ecP256SpkiB64 === "string" ? evt.keyBundle.ecP256SpkiB64 : "";
	const publicKeyB64 = publicKeyB64FromInput || eventPubKeyB64;
	const publicKeySource = publicKeyB64FromInput ? "Public key input" : "keyBundle.ecP256SpkiB64";

	if (!publicKeyB64) {
		render(
			"Provide a P-256 public key in the public key field or ensure keyBundle.ecP256SpkiB64 is present.",
		);
		return;
	}

	let publicKey: CryptoKey;
	try {
		const spkiBytes = base64ToUint8Array(publicKeyB64);
		publicKey = await importP256PublicKey(spkiBytes);
	} catch (err) {
		render(
			`Public key import failed (expecting base64 SPKI or 65-byte uncompressed raw P-256): ${err instanceof Error ? err.message : String(err)}`,
		);
		return;
	}

	let signatureBytes: Uint8Array;
	try {
		signatureBytes = base64ToUint8Array(signatureB64);
	} catch (err) {
		render(
			`Signature is not valid base64 (keyBundle.sodiumKeyBlobB64): ${err instanceof Error ? err.message : String(err)
			}`,
		);
		return;
	}

	const clientToken = typeof evt.cryptoContext === "string" ? evt.cryptoContext : null;
	if (!clientToken) {
		render("Missing cryptoContext/clientToken in event.");
		return;
	}

	const userId = typeof evt.actorUserId === "string" ? evt.actorUserId : null;
	const conversationId = typeof evt.conversationId === "string" ? evt.conversationId : null;
	if (!userId || !conversationId) {
		render("Missing actorUserId or conversationId in event.");
		return;
	}

	const encryptedMessage = evt?.payload?.encryptedMessage;
	const cipherField = encryptedMessage?.ciphertext;
	let ciphertext: Uint8Array | null = null;
	if (cipherField instanceof Uint8Array) {
		ciphertext = cipherField;
	} else if (typeof cipherField === "string") {
		try {
			ciphertext = base64ToUint8Array(cipherField);
		} catch {
			ciphertext = null;
		}
	}

	if (!ciphertext) {
		render("Missing ciphertext at payload.encryptedMessage.ciphertext.");
		return;
	}

	const convKeyVersion =
		parseDecimalString(evt?.keyBundle?.keyCreatedAtMs ?? encryptedMessage?.sessionKeyCreatedAtMs) ??
		null;
	if (!convKeyVersion) {
		render("keyBundle.keyCreatedAtMs missing or not numeric (key version).");
		return;
	}

	const sigVersionRaw = parseDecimalString(evt?.keyBundle?.keyVersion);
	const signatureVersion = sigVersionRaw ? Number(sigVersionRaw) : 2;

	button.disabled = true;
	render("Verifying signature…");

	try {
		const payload = buildSignaturePayloadForEvent(
			signatureVersion,
			clientToken,
			userId,
			conversationId,
			encryptedMessage,
			convKeyVersion,
		);

		if (!payload) {
			render("Could not build signing payload (missing fields or a field contains a comma).");
			return;
		}

		const ok = await verifyEcdsaSignature(publicKey, payload, signatureBytes);

		const lines = [
			ok ? "✓ Signature valid" : "✗ Signature invalid",
			`signature_version: ${signatureVersion}`,
			`keyCreatedAtMs (key version): ${convKeyVersion}`,
			`payload bytes: ${payload.length}`,
			`ciphertext bytes: ${ciphertext.length}`,
			`public key source: ${publicKeySource}`,
		];

		render(lines.join("\n"));
	} catch (err) {
		render(`Verification error: ${err instanceof Error ? err.message : String(err)}`);
	} finally {
		button.disabled = false;
	}
};

const handleBlobSubmit = (event: Event) => {
	event.preventDefault();

	const blob = base64BlobInput.value.trim();
	if (!blob) {
		renderBlobOutput("Paste a base64 blob above.");
		return;
	}

	try {
		const bytes = base64ToUint8Array(blob);
		renderBlobOutput(`Received ${bytes.length} bytes. Add your decryption logic here.`);
	} catch (err) {
		renderBlobOutput(`Invalid base64: ${err instanceof Error ? err.message : String(err)}`);
	}
};

const handleConvEventSubmit = async (event: Event) => {
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

	const eventJson = parseConversationEventInput(raw);

	if (!eventJson) {
		renderConvEventOutput("Could not parse JSON (neither plain nor base64).");
		return;
	}

	const evt = eventJson.event ?? eventJson;

	console.log(evt);

	const payloads = evt?.payload?.encryptedConversationKey?.encryptedKeyPayload;

	if (!Array.isArray(payloads) || payloads.length === 0) {
		renderConvEventOutput(
			"No encrypted conversation key payloads found at payload.encryptedConversationKey.encryptedKeyPayload.",
		);
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
};

keyInput.addEventListener("input", updateKeyPreview);
maskKeyToggle.addEventListener("change", handleMaskToggle);
keySelect.addEventListener("change", handleKeySelectChange);
juiceboxPinInput.addEventListener("input", handlePinInputChange);
saveKeyButton.addEventListener("click", handleSaveKeyClick);
deleteKeyButton.addEventListener("click", handleDeleteKeyClick);
form.addEventListener("submit", handleConnectSubmit);
manualForm.addEventListener("submit", handleManualSubmit);
genericForm.addEventListener("submit", handleGenericSubmit);
manualVerifyBtn.addEventListener("click", () => handleVerifyClick("manual"));
genericVerifyBtn.addEventListener("click", () => handleVerifyClick("generic"));
juiceboxForm.addEventListener("submit", handleJuiceboxSubmit);
base64BlobForm.addEventListener("submit", handleBlobSubmit);
convEventForm.addEventListener("submit", handleConvEventSubmit);

updateKeyPreview();
setManualDecoded(null);
setGenericDecoded(null);

// ---- Encoding & crypto helpers ----
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
	return crypto.subtle.importKey("pkcs8", pkcs8, { name: "ECDH", namedCurve: "P-256" }, false, [
		"deriveBits",
	]);
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

const extractRawP256FromSpki = (spki: Uint8Array): Uint8Array | null => {
	// Minimal DER walk: SEQUENCE -> SEQUENCE (alg) -> BIT STRING
	let o = 0;
	if (spki[o++] !== 0x30) return null;
	const readLen = () => {
		const first = spki[o++];
		if (first === undefined) throw new Error("Unexpected end of SPKI");
		if (first < 0x80) return first;
		const bytes = first & 0x7f;
		if (!bytes || bytes > 4) throw new Error("Unsupported SPKI length form");
		let val = 0;
		for (let i = 0; i < bytes; i++) {
			val = (val << 8) | spki[o++];
		}
		return val;
	};
	try {
		const seqLen = readLen();
		const seqEnd = o + seqLen;
		if (spki[o++] !== 0x30) return null;
		const algLen = readLen();
		o += algLen; // skip algorithm identifiers
		if (spki[o++] !== 0x03) return null;
		const bitLen = readLen();
		if (o >= spki.length) return null;
		const unused = spki[o++];
		if (unused !== 0) return null; // expect aligned bits
		const bitBytes = spki.slice(o, o + bitLen - 1);
		if (bitBytes.length === 65 && bitBytes[0] === 0x04) return bitBytes;
		// Guard against malformed length
		if (seqEnd !== 0 && seqEnd > spki.length) return null;
		return null;
	} catch {
		return null;
	}
};

async function importP256PublicKey(spkiOrRaw: Uint8Array) {
	if (!(globalThis as any).crypto || !(globalThis as any).crypto.subtle) {
		throw new Error("WebCrypto 'crypto.subtle' is not available (serve over https/localhost in a modern browser).");
	}
	const tryImportSpki = () =>
		(globalThis as any).crypto.subtle.importKey(
			"spki",
			asArrayBuffer(spkiOrRaw),
			{ name: "ECDSA", namedCurve: "P-256" },
			true,
			["verify"],
		);

	const tryImportRaw = (raw: Uint8Array) =>
		(globalThis as any).crypto.subtle.importKey(
			"raw",
			asArrayBuffer(raw),
			{ name: "ECDSA", namedCurve: "P-256" },
			true,
			["verify"],
		);

	try {
		return await tryImportSpki();
	} catch (err) {
		if (spkiOrRaw.length === 65 && spkiOrRaw[0] === 0x04) {
			return tryImportRaw(spkiOrRaw);
		}
		const raw = extractRawP256FromSpki(spkiOrRaw);
		if (raw) {
			return tryImportRaw(raw);
		}
		throw err;
	}
}

const rawEcdsaToDer = (raw: Uint8Array) => {
	if (raw.length % 2 !== 0) {
		throw new Error(`Unexpected raw signature length=${raw.length} (want even)`);
	}
	const half = raw.length / 2;
	const r = raw.slice(0, half);
	const s = raw.slice(half);
	const encodeInt = (v: Uint8Array) => {
		let i = 0;
		while (i < v.length && v[i] === 0) i += 1;
		let bytes = v.slice(i);
		if (!bytes.length) bytes = Uint8Array.of(0);
		if (bytes[0] & 0x80) bytes = concatBytes(Uint8Array.of(0), bytes);
		return concatBytes(Uint8Array.of(0x02, bytes.length), bytes);
	};
	const rEnc = encodeInt(r);
	const sEnc = encodeInt(s);
	return concatBytes(Uint8Array.of(0x30, rEnc.length + sEnc.length), rEnc, sEnc);
};

const verifyEcdsaSignature = async (
	publicKey: CryptoKey,
	payload: Uint8Array,
	signature: Uint8Array,
) => {
	const tryVerify = async (sig: Uint8Array) =>
		crypto.subtle.verify({ name: "ECDSA", hash: { name: "SHA-256" } }, publicKey, sig, payload);

	if (await tryVerify(signature)) return true;

	if (signature.length === 64) {
		try {
			const der = rawEcdsaToDer(signature);
			if (await tryVerify(der)) return true;
		} catch {
			/* ignore */
		}
	}

	return false;
};

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
	const oidPrime256v1 = Uint8Array.from([
		0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07,
	]);

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
	const pubKey = await crypto.subtle.importKey(
		"raw",
		ephPub,
		{ name: "ECDH", namedCurve: "P-256" },
		false,
		[],
	);

	const sharedBits = await crypto.subtle.deriveBits(
		{ name: "ECDH", public: pubKey },
		privKey,
		256,
	);
	const shared = new Uint8Array(sharedBits);

	const keyNonce = await kdf2Sha256(shared, ephPub, 32);
	const aesKeyBytes = keyNonce.slice(0, 16);
	const iv = keyNonce.slice(16);

	const aesKey = await crypto.subtle.importKey(
		"raw",
		asArrayBuffer(aesKeyBytes),
		"AES-GCM",
		false,
		["decrypt"],
	);

	let plaintext: Uint8Array;
	try {
		const plaintextBuf = await crypto.subtle.decrypt(
			{ name: "AES-GCM", iv: asArrayBuffer(iv), tagLength: 128 },
			aesKey,
			asArrayBuffer(cipherAndTag),
		);
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
	// Normalize base64url, strip whitespace/PEM headers, and fix padding.
	b64 = b64
		.trim()
		.replace(/-----BEGIN [^-]+-----/gi, "")
		.replace(/-----END [^-]+-----/gi, "")
		.replace(/\s+/g, "")
		.replace(/-/g, "+")
		.replace(/_/g, "/");
	const pad = b64.length % 4;
	if (pad > 0) b64 += "=".repeat(4 - pad); // be permissive; let atob validate
	const binary = atob(b64);
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i) & 0xff;
	return bytes;
}

function base64UrlNoPad(data: Uint8Array) {
	return sodium.to_base64(data, sodium.base64_variants.URLSAFE_NO_PADDING);
}

const parseDecimalString = (raw: unknown): string | null => {
	if (typeof raw !== "string") return null;
	const trimmed = raw.trim();
	if (!/^\d+$/.test(trimmed)) return null;
	try {
		return BigInt(trimmed).toString(10);
	} catch {
		return null;
	}
};

const joinFieldsOrNull = (fields: Array<string | null | undefined>): Uint8Array | null => {
	const parts: string[] = [];
	for (const f of fields) {
		if (f === null || f === undefined) continue;
		if (f.includes(",")) return null;
		parts.push(f);
	}
	return new TextEncoder().encode(parts.join(","));
};

const joinEventBaseAndExtrasOrNull = (
	eventType: string,
	clientToken: string,
	userId: string,
	conversationId: string | null,
	extras: string[],
): Uint8Array | null => {
	const base = [eventType, clientToken, userId, conversationId ?? null];
	return joinFieldsOrNull([...base, ...extras]);
};

const buildMessageCreatePayload = (
	signatureVersion: number,
	clientToken: string,
	userId: string,
	conversationId: string,
	ciphertext: Uint8Array,
	conversationKeyVersion: string | null,
): Uint8Array | null => {
	if (!ciphertext?.length) return null;
	const ciphertextB64 = base64UrlNoPad(ciphertext);

	if (signatureVersion === 1) {
		return joinFieldsOrNull([clientToken, userId, conversationId, ciphertextB64]);
	}

	const convKeyVersionNorm = parseDecimalString(conversationKeyVersion);
	if (!convKeyVersionNorm) return null;

	const extras = [convKeyVersionNorm, ciphertextB64];
	return joinEventBaseAndExtrasOrNull(
		"MessageCreateEvent",
		clientToken,
		userId,
		conversationId,
		extras,
	);
};

const buildSignaturePayloadForEvent = (
	signatureVersion: number,
	clientToken: string,
	userId: string,
	conversationId: string | null,
	eventDetail: any,
	conversationKeyVersion: string | null,
): Uint8Array | null => {
	if (!eventDetail) return null;

	const cipherField = eventDetail?.ciphertext;
	let ciphertext: Uint8Array | null = null;
	if (cipherField instanceof Uint8Array) {
		ciphertext = cipherField;
	} else if (typeof cipherField === "string") {
		try {
			ciphertext = base64ToUint8Array(cipherField);
		} catch {
			ciphertext = null;
		}
	}

	if (!ciphertext || !conversationId) return null;

	return buildMessageCreatePayload(
		signatureVersion,
		clientToken,
		userId,
		conversationId,
		ciphertext,
		conversationKeyVersion,
	);
};

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

const defaultEncode =
	(t: TCode) =>
		(v: any): Uint8Array => {
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
				const header = concat(
					Uint8Array.of(kt),
					Uint8Array.of(vt),
					enc.i32(entries.length / 2),
				);
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
					throw new Error(
						`Element type mismatch: schema=${f.elemType}, wire=${etWire} (field=${f.key})`,
					);
				}
				const arr: any[] = [];
				for (let i = 0; i < count; i++) {
					if (etWire === T.STRUCT) {
						const elemSchema = f.elemSchema;
						arr.push(
							elemSchema ? this.readStruct(elemSchema) : this.readStructGeneric(),
						);
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
					throw new Error(
						`Map key type mismatch: schema=${f.keyType}, wire=${ktWire} (field=${f.key})`,
					);
				}
				if (f.valType !== undefined && f.valType !== vtWire) {
					throw new Error(
						`Map val type mismatch: schema=${f.valType}, wire=${vtWire} (field=${f.key})`,
					);
				}
				const obj: any = {};
				for (let i = 0; i < count; i++) {
					const k =
						ktWire === T.STRUCT
							? this.readStruct(f.keySchema || [])
							: this.readValue(ktWire, { id: 0, key: "", type: ktWire });
					const v =
						vtWire === T.STRUCT
							? this.readStruct(f.valSchema || [])
							: this.readValue(vtWire, { id: 0, key: "", type: vtWire });
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

// Thrift types generated from twittermeow/data/payload/thrift.go (dmv2.thriftjava)
export enum AdditionalAction {
	AdditionalActionFetchConvIfMissingCkey = 0,
}

export type AdditionalActionCode = AdditionalAction;

export enum DeleteMessageAction {
	DeleteMessageActionDeleteForSelf = 1,
	DeleteMessageActionDeleteForAll = 2,
}

export type DeleteMessageActionCode = DeleteMessageAction;

export enum FailureType {
	FailureTypeEmptyDetail = 1,
	FailureTypeInternalError = 2,
	FailureTypeContentsTooLarge = 3,
	FailureTypeTooManyMessages = 4,
	FailureTypeInvalidSenderSignature = 5,
	FailureTypeNonLatestCkeyVersion = 6,
	FailureTypeRecipientHasNotTrustedConversation = 7,
	FailureTypeRecipientKeyHasChanged = 8,
}

export type FailureTypeCode = FailureType;

export enum MediaType {
	MediaTypeImage = 1,
	MediaTypeGif = 2,
	MediaTypeVideo = 3,
	MediaTypeAudio = 4,
	MediaTypeFile = 5,
	MediaTypeSvg = 6,
}

export type MediaTypeCode = MediaType;

export enum ScreenCaptureType {
	ScreenCaptureTypeScreenshot = 1,
	ScreenCaptureTypeRecording = 2,
}

export type ScreenCaptureTypeCode = ScreenCaptureType;

export enum SentFromSurface {
	SentFromSurfaceConversationScreenComposer = 1,
	SentFromSurfaceNotificationReply = 2,
	SentFromSurfaceShareSheet = 3,
	SentFromSurfacePaymentsSupportComposer = 4,
	SentFromSurfaceMessageForwardSheet = 5,
}

export type SentFromSurfaceCode = SentFromSurface;

export type AVCallEnded = {
	sent_at_millis?: bigint;
	duration_seconds?: bigint;
	is_audio_only?: boolean;
	broadcast_id?: string;
};

export type AVCallMissed = {
	sent_at_millis?: bigint;
	is_audio_only?: boolean;
};

export type AVCallStarted = {
	is_audio_only?: boolean;
	broadcast_id?: string;
};

export type AcceptMessageRequest = {
};

export type AddressRichTextContent = {
};

export type BatchedMessageEvents = {
	message_events?: MessageEvent[];
};

export type CallToAction = {
	label?: string;
	url?: string;
};

export type CashtagRichTextContent = {
};

export type ConversationDeleteEvent = {
	conversation_id?: string;
};

export type ConversationKeyChangeEvent = {
	conversation_key_version?: string;
	conversation_participant_keys?: ConversationParticipantKey[];
	ratchet_tree?: KeyRotation;
};

export type ConversationMetadataChange = {
	message_duration_change?: MessageDurationChange;
	message_duration_remove?: MessageDurationRemove;
	mute_conversation?: MuteConversation;
	unmute_conversation?: UnmuteConversation;
	enable_screen_capture_detection?: EnableScreenCaptureDetection;
	disable_screen_capture_detection?: DisableScreenCaptureDetection;
	enable_screen_capture_blocking?: EnableScreenCaptureBlocking;
	disable_screen_capture_blocking?: DisableScreenCaptureBlocking;
};

export type ConversationMetadataChangeEvent = {
	conversation_metadata_change?: ConversationMetadataChange;
};

export type ConversationParticipantKey = {
	user_id?: string;
	encrypted_conversation_key?: string;
	public_key_version?: string;
};

export type DisableScreenCaptureBlocking = {
	placeholder?: string;
};

export type DisableScreenCaptureDetection = {
	placeholder?: string;
};

export type DisplayTemporaryPasscodeInstruction = {
	token?: string;
	latest_public_key_version?: string;
};

export type DraftMessage = {
	conversation_id?: string;
	draft_text?: string;
};

export type EmailRichTextContent = {
};

export type EmptyNode = {
	description?: string;
};

export type EnableScreenCaptureBlocking = {
	placeholder?: string;
};

export type EnableScreenCaptureDetection = {
	placeholder?: string;
};

export type EventQueuePriority = {
};

export type ForwardedMessage = {
	message_text?: string;
	entities?: RichTextEntity[];
};

export type GrokSearchResponseEvent = {
	search_response_id?: string;
};

export type GroupAdminAddChange = {
	admin_ids?: string[];
};

export type GroupAdminRemoveChange = {
	admin_ids?: string[];
};

export type GroupAvatarUrlChange = {
	custom_avatar_url?: string;
	conversation_key_version?: string;
};

export type GroupChange = {
	group_create?: GroupCreate;
	group_title_change?: GroupTitleChange;
	group_avatar_change?: GroupAvatarUrlChange;
	group_admin_add?: GroupAdminAddChange;
	group_member_add?: GroupMemberAddChange;
	group_admin_remove?: GroupAdminRemoveChange;
	group_member_remove?: GroupMemberRemoveChange;
	group_invite_enable?: GroupInviteEnable;
	group_invite_disable?: GroupInviteDisable;
	group_join_request?: GroupJoinRequest;
	group_join_reject?: GroupJoinReject;
};

export type GroupChangeEvent = {
	group_change?: GroupChange;
};

export type GroupCreate = {
	member_ids?: string[];
	admin_ids?: string[];
	title?: string;
	avatar_url?: string;
	conversation_key_version?: string;
};

export type GroupInviteDisable = {
	disabled_by_member_id?: string;
};

export type GroupInviteEnable = {
	expires_at_msec?: bigint;
	invite_url?: string;
	affiliate_id?: string;
};

export type GroupJoinReject = {
	rejected_user_ids?: string[];
};

export type GroupJoinRequest = {
	requesting_user_id?: string;
};

export type GroupMemberAddChange = {
	member_ids?: string[];
	current_member_ids?: string[];
	current_admin_ids?: string[];
	current_title?: string;
	current_avatar_url?: string;
	conversation_key_version?: string;
	current_ttl_msec?: bigint;
	current_pending_member_ids?: string[];
};

export type GroupMemberRemoveChange = {
	member_ids?: string[];
};

export type GroupTitleChange = {
	custom_title?: string;
	conversation_key_version?: string;
};

export type HashtagRichTextContent = {
};

export type KeepAliveInstruction = {
};

export type KeyRotation = {
	previous_version?: string;
	ratchet_tree?: RatchetTree;
	nodes?: UpdatePathNode[];
	encrypted_private_key?: string;
};

export type LeafNode = {
	subtree_encryption_public_key?: string;
	signature_public_key?: string;
	keypair_id?: string;
	max_supported_protocol_version?: number;
	parent_hash?: string;
	signature?: string;
};

export type MarkConversationRead = {
	seen_until_sequence_id?: string;
	seen_at_millis?: bigint;
};

export type MarkConversationReadEvent = {
	seen_until_sequence_id?: string;
	seen_at_millis?: bigint;
};

export type MarkConversationUnread = {
	seen_until_sequence_id?: string;
};

export type MarkConversationUnreadEvent = {
	seen_until_sequence_id?: string;
};

export type MaybeKeypair = {
	empty?: string;
	keypair?: StoredKeypair;
};

export type MediaAttachment = {
	media_hash_key?: string;
	dimensions?: MediaDimensions;
	type?: number;
	duration_millis?: bigint;
	filesize_bytes?: bigint;
	filename?: string;
	attachment_id?: string;
	legacy_media_url_https?: string;
	legacy_media_preview_url?: string;
};

export type MediaDimensions = {
	width?: bigint;
	height?: bigint;
};

export type MemberAccountDeleteEvent = {
	member_id?: string;
};

export type MentionRichTextContent = {
};

export type Message = {
	messageEvent?: MessageEvent;
	messageInstruction?: MessageInstruction;
	batchedMessageEvents?: BatchedMessageEvents;
};

export type MessageAttachment = {
	media?: MediaAttachment;
	post?: PostAttachment;
	url?: UrlAttachment;
	unified_card?: UnifiedCardAttachment;
	money?: MoneyAttachment;
};

export type MessageContents = {
	message_text?: string;
	entities?: RichTextEntity[];
	attachments?: MessageAttachment[];
	replying_to_preview?: ReplyingToPreview;
	forwarded_message?: ForwardedMessage;
	sent_from?: number;
	quick_reply?: QuickReply;
	ctas?: CallToAction[];
};

export type MessageCreateEvent = {
	contents?: string;
	conversation_key_version?: string;
	should_notify?: boolean;
	ttl_msec?: bigint;
	delivered_at_msec?: bigint;
	is_pending_public_key?: boolean;
	priority?: number;
	additional_action_list?: number[];
};

export type MessageDeleteEvent = {
	sequence_ids?: string[];
	delete_message_action?: number;
};

export type MessageDurationChange = {
	ttl_msec?: bigint;
};

export type MessageDurationRemove = {
	current_ttl_msec?: bigint;
};

export type MessageEdit = {
	message_sequence_id?: string;
	updated_text?: string;
	entities?: RichTextEntity[];
};

export type MessageEntryContents = {
	reaction_add?: MessageReactionAdd;
	reaction_remove?: MessageReactionRemove;
	message_edit?: MessageEdit;
	mark_conversation_read?: MarkConversationRead;
	mark_conversation_unread?: MarkConversationUnread;
	pin_conversation?: PinConversation;
	unpin_conversation?: UnpinConversation;
	screen_capture_detected?: ScreenCaptureDetected;
	av_call_ended?: AVCallEnded;
	av_call_missed?: AVCallMissed;
	draft_message?: DraftMessage;
	accept_message_request?: AcceptMessageRequest;
	nickname_message?: NicknameMessage;
	set_verified_status?: SetVerifiedStatus;
	av_call_started?: AVCallStarted;
};

export type MessageEntryHolder = {
	contents?: MessageEntryContents;
};

export type MessageEvent = {
	sequence_id?: string;
	message_id?: string;
	sender_id?: string;
	conversation_id?: string;
	conversation_token?: string;
	created_at_msec?: string;
	detail?: MessageEventDetail;
	relay_source?: number;
	message_event_signature?: MessageEventSignature;
	previous_sequence_id?: string;
	is_trusted?: boolean;
};

export type MessageEventDetail = {
	messageCreateEvent?: MessageCreateEvent;
	conversationKeyChangeEvent?: ConversationKeyChangeEvent;
	groupChangeEvent?: GroupChangeEvent;
	messageFailureEvent?: MessageFailureEvent;
	messageTypingEvent?: MessageTypingEvent;
	messageDeleteEvent?: MessageDeleteEvent;
	conversationDeleteEvent?: ConversationDeleteEvent;
	conversationMetadataChangeEvent?: ConversationMetadataChangeEvent;
	grokSearchResponseEvent?: GrokSearchResponseEvent;
	requestForEncryptedResendEvent?: RequestForEncryptedResendEvent;
	markConversationReadEvent?: MarkConversationReadEvent;
	markConversationUnreadEvent?: MarkConversationUnreadEvent;
	memberAccountDeleteEvent?: MemberAccountDeleteEvent;
};

export type MessageEventRelaySource = {
};

export type MessageEventSignature = {
	signature?: string;
	public_key_version?: string;
	signature_version?: string;
	signing_public_key?: string;
};

export type MessageFailureEvent = {
	failure_type?: number;
};

export type MessageInstruction = {
	pullMessagesInstruction?: PullMessagesInstruction;
	keepAliveInstruction?: KeepAliveInstruction;
	pullMessagesFinishedInstruction?: PullMessagesFinishedInstruction;
	pinReminderInstruction?: PinReminderInstruction;
	switchToHybridPullInstruction?: SwitchToHybridPullInstruction;
	displayTemporaryPasscodeInstruction?: DisplayTemporaryPasscodeInstruction;
};

export type MessageReactionAdd = {
	message_sequence_id?: string;
	emoji?: string;
};

export type MessageReactionRemove = {
	message_sequence_id?: string;
	emoji?: string;
};

export type MessageTypingEvent = {
	conversation_id?: string;
};

export type MoneyAttachment = {
	fallbackText?: string;
	payload?: string;
};

export type MuteConversation = {
	muted_conversation_ids?: string[];
};

export type NicknameMessage = {
	user_id?: bigint;
	nickname_text?: string;
};

export type ParentNode = {
	subtree_encryption_public_key?: string;
	parent_hash?: string;
};

export type PhoneNumberRichTextContent = {
};

export type PinConversation = {
	conversation_id?: string;
};

export type PinReminderInstruction = {
	should_register?: boolean;
	should_generate?: boolean;
};

export type PostAttachment = {
	rest_id?: string;
	post_url?: string;
	attachment_id?: string;
};

export type PullMessagePageDetails = {
	min_sequence_id?: string;
	max_sequence_id?: string;
	is_batched_pull?: boolean;
};

export type PullMessagesFinishedInstruction = {
	finished_pull?: boolean;
	sequence_continue?: string;
	pull_message_page_details?: PullMessagePageDetails;
};

export type PullMessagesInstruction = {
	sequence_start?: string;
	sender_id?: string;
	is_batched_pull?: boolean;
};

export type QuickReply = {
	request?: QuickReplyRequest;
	response?: QuickReplyResponse;
};

export type QuickReplyOption = {
	id?: string;
	label?: string;
	metadata?: string;
	description?: string;
};

export type QuickReplyOptionsRequest = {
	id?: string;
	options?: QuickReplyOption[];
};

export type QuickReplyOptionsResponse = {
	request_id?: string;
	metadata?: string;
	selected_option_id?: string;
};

export type QuickReplyRequest = {
	options?: QuickReplyOptionsRequest;
};

export type QuickReplyResponse = {
	options?: QuickReplyOptionsResponse;
};

export type RatchetTree = {
	leaves?: RatchetTreeLeaf[];
	parents?: RatchetTreeParent[];
};

export type RatchetTreeLeaf = {
	empty?: EmptyNode;
	leaf?: LeafNode;
};

export type RatchetTreeParent = {
	empty?: EmptyNode;
	parent?: ParentNode;
};

export type ReplyingToPreview = {
	sender_id?: bigint;
	message_text?: string;
	entities?: RichTextEntity[];
	attachments?: MessageAttachment[];
	sender_display_name?: string;
	replying_to_message_sequence_id?: string;
	replying_to_message_id?: string;
};

export type RequestForEncryptedResendEvent = {
	min_sequence_id?: string;
	max_sequence_id?: string;
};

export type RichTextContent = {
	hashtag?: HashtagRichTextContent;
	cashtag?: CashtagRichTextContent;
	mention?: MentionRichTextContent;
	url?: UrlRichTextContent;
	email?: EmailRichTextContent;
	phoneNumber?: PhoneNumberRichTextContent;
};

export type RichTextEntity = {
	start_index?: number;
	end_index?: number;
	content?: RichTextContent;
};

export type ScreenCaptureDetected = {
	type?: number;
};

export type SetVerifiedStatus = {
	user_id?: bigint;
	verified_status?: boolean;
};

export type StoredGroupState = {
	keypairs?: MaybeKeypair[];
	ratchet_tree?: RatchetTree;
};

export type StoredKeypair = {
	public_key?: string;
	private_key?: string;
};

export type SwitchToHybridPullInstruction = {
	requesting_user_agent?: string;
};

export type UnifiedCardAttachment = {
	url?: string;
	attachment_id?: string;
};

export type UnmuteConversation = {
	unmuted_conversation_ids?: string[];
};

export type UnpinConversation = {
	conversation_id?: string;
};

export type UpdatePathNode = {
	encrypted_secrets?: string[];
	encrypted_private_key?: string;
};

export type UrlAttachment = {
	url?: string;
	banner_image_media_hash_key?: UrlAttachmentImage;
	favicon_image_media_hash_key?: UrlAttachmentImage;
	display_title?: string;
	attachment_id?: string;
};

export type UrlAttachmentImage = {
	media_hash_key?: string;
	filesize_bytes?: bigint;
	filename?: string;
	dimensions?: MediaDimensions;
};

export type UrlRichTextContent = {
};

export const aVCallEndedSchema: FieldSchema[] = [];
export const aVCallMissedSchema: FieldSchema[] = [];
export const aVCallStartedSchema: FieldSchema[] = [];
export const acceptMessageRequestSchema: FieldSchema[] = [];
export const addressRichTextContentSchema: FieldSchema[] = [];
export const batchedMessageEventsSchema: FieldSchema[] = [];
export const callToActionSchema: FieldSchema[] = [];
export const cashtagRichTextContentSchema: FieldSchema[] = [];
export const conversationDeleteEventSchema: FieldSchema[] = [];
export const conversationKeyChangeEventSchema: FieldSchema[] = [];
export const conversationMetadataChangeSchema: FieldSchema[] = [];
export const conversationMetadataChangeEventSchema: FieldSchema[] = [];
export const conversationParticipantKeySchema: FieldSchema[] = [];
export const disableScreenCaptureBlockingSchema: FieldSchema[] = [];
export const disableScreenCaptureDetectionSchema: FieldSchema[] = [];
export const displayTemporaryPasscodeInstructionSchema: FieldSchema[] = [];
export const draftMessageSchema: FieldSchema[] = [];
export const emailRichTextContentSchema: FieldSchema[] = [];
export const emptyNodeSchema: FieldSchema[] = [];
export const enableScreenCaptureBlockingSchema: FieldSchema[] = [];
export const enableScreenCaptureDetectionSchema: FieldSchema[] = [];
export const eventQueuePrioritySchema: FieldSchema[] = [];
export const forwardedMessageSchema: FieldSchema[] = [];
export const grokSearchResponseEventSchema: FieldSchema[] = [];
export const groupAdminAddChangeSchema: FieldSchema[] = [];
export const groupAdminRemoveChangeSchema: FieldSchema[] = [];
export const groupAvatarUrlChangeSchema: FieldSchema[] = [];
export const groupChangeSchema: FieldSchema[] = [];
export const groupChangeEventSchema: FieldSchema[] = [];
export const groupCreateSchema: FieldSchema[] = [];
export const groupInviteDisableSchema: FieldSchema[] = [];
export const groupInviteEnableSchema: FieldSchema[] = [];
export const groupJoinRejectSchema: FieldSchema[] = [];
export const groupJoinRequestSchema: FieldSchema[] = [];
export const groupMemberAddChangeSchema: FieldSchema[] = [];
export const groupMemberRemoveChangeSchema: FieldSchema[] = [];
export const groupTitleChangeSchema: FieldSchema[] = [];
export const hashtagRichTextContentSchema: FieldSchema[] = [];
export const keepAliveInstructionSchema: FieldSchema[] = [];
export const keyRotationSchema: FieldSchema[] = [];
export const leafNodeSchema: FieldSchema[] = [];
export const markConversationReadSchema: FieldSchema[] = [];
export const markConversationReadEventSchema: FieldSchema[] = [];
export const markConversationUnreadSchema: FieldSchema[] = [];
export const markConversationUnreadEventSchema: FieldSchema[] = [];
export const maybeKeypairSchema: FieldSchema[] = [];
export const mediaAttachmentSchema: FieldSchema[] = [];
export const mediaDimensionsSchema: FieldSchema[] = [];
export const memberAccountDeleteEventSchema: FieldSchema[] = [];
export const mentionRichTextContentSchema: FieldSchema[] = [];
export const messageSchema: FieldSchema[] = [];
export const messageAttachmentSchema: FieldSchema[] = [];
export const messageContentsSchema: FieldSchema[] = [];
export const messageCreateEventSchema: FieldSchema[] = [];
export const messageDeleteEventSchema: FieldSchema[] = [];
export const messageDurationChangeSchema: FieldSchema[] = [];
export const messageDurationRemoveSchema: FieldSchema[] = [];
export const messageEditSchema: FieldSchema[] = [];
export const messageEntryContentsSchema: FieldSchema[] = [];
export const messageEntryHolderSchema: FieldSchema[] = [];
export const messageEventSchema: FieldSchema[] = [];
export const messageEventDetailSchema: FieldSchema[] = [];
export const messageEventRelaySourceSchema: FieldSchema[] = [];
export const messageEventSignatureSchema: FieldSchema[] = [];
export const messageFailureEventSchema: FieldSchema[] = [];
export const messageInstructionSchema: FieldSchema[] = [];
export const messageReactionAddSchema: FieldSchema[] = [];
export const messageReactionRemoveSchema: FieldSchema[] = [];
export const messageTypingEventSchema: FieldSchema[] = [];
export const moneyAttachmentSchema: FieldSchema[] = [];
export const muteConversationSchema: FieldSchema[] = [];
export const nicknameMessageSchema: FieldSchema[] = [];
export const parentNodeSchema: FieldSchema[] = [];
export const phoneNumberRichTextContentSchema: FieldSchema[] = [];
export const pinConversationSchema: FieldSchema[] = [];
export const pinReminderInstructionSchema: FieldSchema[] = [];
export const postAttachmentSchema: FieldSchema[] = [];
export const pullMessagePageDetailsSchema: FieldSchema[] = [];
export const pullMessagesFinishedInstructionSchema: FieldSchema[] = [];
export const pullMessagesInstructionSchema: FieldSchema[] = [];
export const quickReplySchema: FieldSchema[] = [];
export const quickReplyOptionSchema: FieldSchema[] = [];
export const quickReplyOptionsRequestSchema: FieldSchema[] = [];
export const quickReplyOptionsResponseSchema: FieldSchema[] = [];
export const quickReplyRequestSchema: FieldSchema[] = [];
export const quickReplyResponseSchema: FieldSchema[] = [];
export const ratchetTreeSchema: FieldSchema[] = [];
export const ratchetTreeLeafSchema: FieldSchema[] = [];
export const ratchetTreeParentSchema: FieldSchema[] = [];
export const replyingToPreviewSchema: FieldSchema[] = [];
export const requestForEncryptedResendEventSchema: FieldSchema[] = [];
export const richTextContentSchema: FieldSchema[] = [];
export const richTextEntitySchema: FieldSchema[] = [];
export const screenCaptureDetectedSchema: FieldSchema[] = [];
export const setVerifiedStatusSchema: FieldSchema[] = [];
export const storedGroupStateSchema: FieldSchema[] = [];
export const storedKeypairSchema: FieldSchema[] = [];
export const switchToHybridPullInstructionSchema: FieldSchema[] = [];
export const unifiedCardAttachmentSchema: FieldSchema[] = [];
export const unmuteConversationSchema: FieldSchema[] = [];
export const unpinConversationSchema: FieldSchema[] = [];
export const updatePathNodeSchema: FieldSchema[] = [];
export const urlAttachmentSchema: FieldSchema[] = [];
export const urlAttachmentImageSchema: FieldSchema[] = [];
export const urlRichTextContentSchema: FieldSchema[] = [];

aVCallEndedSchema.push(
	{ id: 1, key: "sent_at_millis", type: T.I64 },
	{ id: 2, key: "duration_seconds", type: T.I64 },
	{ id: 3, key: "is_audio_only", type: T.BOOL },
	{ id: 5, key: "broadcast_id", type: T.STRING },
);

aVCallMissedSchema.push(
	{ id: 1, key: "sent_at_millis", type: T.I64 },
	{ id: 2, key: "is_audio_only", type: T.BOOL },
);

aVCallStartedSchema.push(
	{ id: 1, key: "is_audio_only", type: T.BOOL },
	{ id: 3, key: "broadcast_id", type: T.STRING },
);

acceptMessageRequestSchema.push(
);

addressRichTextContentSchema.push(
);

batchedMessageEventsSchema.push(
	{ id: 1, key: "message_events", type: T.LIST, elemType: T.STRUCT, elemSchema: messageEventSchema },
);

callToActionSchema.push(
	{ id: 1, key: "label", type: T.STRING },
	{ id: 2, key: "url", type: T.STRING },
);

cashtagRichTextContentSchema.push(
);

conversationDeleteEventSchema.push(
	{ id: 1, key: "conversation_id", type: T.STRING },
);

conversationKeyChangeEventSchema.push(
	{ id: 1, key: "conversation_key_version", type: T.STRING },
	{ id: 2, key: "conversation_participant_keys", type: T.LIST, elemType: T.STRUCT, elemSchema: conversationParticipantKeySchema },
	{ id: 3, key: "ratchet_tree", type: T.STRUCT, schema: keyRotationSchema },
);

conversationMetadataChangeSchema.push(
	{ id: 1, key: "message_duration_change", type: T.STRUCT, schema: messageDurationChangeSchema },
	{ id: 2, key: "message_duration_remove", type: T.STRUCT, schema: messageDurationRemoveSchema },
	{ id: 3, key: "mute_conversation", type: T.STRUCT, schema: muteConversationSchema },
	{ id: 4, key: "unmute_conversation", type: T.STRUCT, schema: unmuteConversationSchema },
	{ id: 5, key: "enable_screen_capture_detection", type: T.STRUCT, schema: enableScreenCaptureDetectionSchema },
	{ id: 6, key: "disable_screen_capture_detection", type: T.STRUCT, schema: disableScreenCaptureDetectionSchema },
	{ id: 7, key: "enable_screen_capture_blocking", type: T.STRUCT, schema: enableScreenCaptureBlockingSchema },
	{ id: 8, key: "disable_screen_capture_blocking", type: T.STRUCT, schema: disableScreenCaptureBlockingSchema },
);

conversationMetadataChangeEventSchema.push(
	{ id: 1, key: "conversation_metadata_change", type: T.STRUCT, schema: conversationMetadataChangeSchema },
);

conversationParticipantKeySchema.push(
	{ id: 1, key: "user_id", type: T.STRING },
	{ id: 2, key: "encrypted_conversation_key", type: T.STRING },
	{ id: 3, key: "public_key_version", type: T.STRING },
);

disableScreenCaptureBlockingSchema.push(
	{ id: 1, key: "placeholder", type: T.STRING },
);

disableScreenCaptureDetectionSchema.push(
	{ id: 1, key: "placeholder", type: T.STRING },
);

displayTemporaryPasscodeInstructionSchema.push(
	{ id: 1, key: "token", type: T.STRING },
	{ id: 2, key: "latest_public_key_version", type: T.STRING },
);

draftMessageSchema.push(
	{ id: 1, key: "conversation_id", type: T.STRING },
	{ id: 2, key: "draft_text", type: T.STRING },
);

emailRichTextContentSchema.push(
);

emptyNodeSchema.push(
	{ id: 1, key: "description", type: T.STRING },
);

enableScreenCaptureBlockingSchema.push(
	{ id: 1, key: "placeholder", type: T.STRING },
);

enableScreenCaptureDetectionSchema.push(
	{ id: 1, key: "placeholder", type: T.STRING },
);

eventQueuePrioritySchema.push(
);

forwardedMessageSchema.push(
	{ id: 1, key: "message_text", type: T.STRING },
	{ id: 2, key: "entities", type: T.LIST, elemType: T.STRUCT, elemSchema: richTextEntitySchema },
);

grokSearchResponseEventSchema.push(
	{ id: 1, key: "search_response_id", type: T.STRING },
);

groupAdminAddChangeSchema.push(
	{ id: 1, key: "admin_ids", type: T.LIST, elemType: T.STRING },
);

groupAdminRemoveChangeSchema.push(
	{ id: 1, key: "admin_ids", type: T.LIST, elemType: T.STRING },
);

groupAvatarUrlChangeSchema.push(
	{ id: 1, key: "custom_avatar_url", type: T.STRING },
	{ id: 2, key: "conversation_key_version", type: T.STRING },
);

groupChangeSchema.push(
	{ id: 1, key: "group_create", type: T.STRUCT, schema: groupCreateSchema },
	{ id: 2, key: "group_title_change", type: T.STRUCT, schema: groupTitleChangeSchema },
	{ id: 3, key: "group_avatar_change", type: T.STRUCT, schema: groupAvatarUrlChangeSchema },
	{ id: 4, key: "group_admin_add", type: T.STRUCT, schema: groupAdminAddChangeSchema },
	{ id: 5, key: "group_member_add", type: T.STRUCT, schema: groupMemberAddChangeSchema },
	{ id: 6, key: "group_admin_remove", type: T.STRUCT, schema: groupAdminRemoveChangeSchema },
	{ id: 7, key: "group_member_remove", type: T.STRUCT, schema: groupMemberRemoveChangeSchema },
	{ id: 8, key: "group_invite_enable", type: T.STRUCT, schema: groupInviteEnableSchema },
	{ id: 9, key: "group_invite_disable", type: T.STRUCT, schema: groupInviteDisableSchema },
	{ id: 10, key: "group_join_request", type: T.STRUCT, schema: groupJoinRequestSchema },
	{ id: 11, key: "group_join_reject", type: T.STRUCT, schema: groupJoinRejectSchema },
);

groupChangeEventSchema.push(
	{ id: 1, key: "group_change", type: T.STRUCT, schema: groupChangeSchema },
);

groupCreateSchema.push(
	{ id: 1, key: "member_ids", type: T.LIST, elemType: T.STRING },
	{ id: 2, key: "admin_ids", type: T.LIST, elemType: T.STRING },
	{ id: 3, key: "title", type: T.STRING },
	{ id: 4, key: "avatar_url", type: T.STRING },
	{ id: 5, key: "conversation_key_version", type: T.STRING },
);

groupInviteDisableSchema.push(
	{ id: 1, key: "disabled_by_member_id", type: T.STRING },
);

groupInviteEnableSchema.push(
	{ id: 1, key: "expires_at_msec", type: T.I64 },
	{ id: 2, key: "invite_url", type: T.STRING },
	{ id: 3, key: "affiliate_id", type: T.STRING },
);

groupJoinRejectSchema.push(
	{ id: 1, key: "rejected_user_ids", type: T.LIST, elemType: T.STRING },
);

groupJoinRequestSchema.push(
	{ id: 1, key: "requesting_user_id", type: T.STRING },
);

groupMemberAddChangeSchema.push(
	{ id: 1, key: "member_ids", type: T.LIST, elemType: T.STRING },
	{ id: 2, key: "current_member_ids", type: T.LIST, elemType: T.STRING },
	{ id: 3, key: "current_admin_ids", type: T.LIST, elemType: T.STRING },
	{ id: 4, key: "current_title", type: T.STRING },
	{ id: 5, key: "current_avatar_url", type: T.STRING },
	{ id: 6, key: "conversation_key_version", type: T.STRING },
	{ id: 7, key: "current_ttl_msec", type: T.I64 },
	{ id: 8, key: "current_pending_member_ids", type: T.LIST, elemType: T.STRING },
);

groupMemberRemoveChangeSchema.push(
	{ id: 1, key: "member_ids", type: T.LIST, elemType: T.STRING },
);

groupTitleChangeSchema.push(
	{ id: 1, key: "custom_title", type: T.STRING },
	{ id: 2, key: "conversation_key_version", type: T.STRING },
);

hashtagRichTextContentSchema.push(
);

keepAliveInstructionSchema.push(
);

keyRotationSchema.push(
	{ id: 1, key: "previous_version", type: T.STRING },
	{ id: 2, key: "ratchet_tree", type: T.STRUCT, schema: ratchetTreeSchema },
	{ id: 3, key: "nodes", type: T.LIST, elemType: T.STRUCT, elemSchema: updatePathNodeSchema },
	{ id: 4, key: "encrypted_private_key", type: T.STRING },
);

leafNodeSchema.push(
	{ id: 1, key: "subtree_encryption_public_key", type: T.STRING },
	{ id: 2, key: "signature_public_key", type: T.STRING },
	{ id: 3, key: "keypair_id", type: T.STRING },
	{ id: 4, key: "max_supported_protocol_version", type: T.I32 },
	{ id: 5, key: "parent_hash", type: T.STRING },
	{ id: 6, key: "signature", type: T.STRING },
);

markConversationReadSchema.push(
	{ id: 1, key: "seen_until_sequence_id", type: T.STRING },
	{ id: 2, key: "seen_at_millis", type: T.I64 },
);

markConversationReadEventSchema.push(
	{ id: 1, key: "seen_until_sequence_id", type: T.STRING },
	{ id: 2, key: "seen_at_millis", type: T.I64 },
);

markConversationUnreadSchema.push(
	{ id: 1, key: "seen_until_sequence_id", type: T.STRING },
);

markConversationUnreadEventSchema.push(
	{ id: 1, key: "seen_until_sequence_id", type: T.STRING },
);

maybeKeypairSchema.push(
	{ id: 1, key: "empty", type: T.STRING },
	{ id: 2, key: "keypair", type: T.STRUCT, schema: storedKeypairSchema },
);

mediaAttachmentSchema.push(
	{ id: 1, key: "media_hash_key", type: T.STRING },
	{ id: 2, key: "dimensions", type: T.STRUCT, schema: mediaDimensionsSchema },
	{ id: 3, key: "type", type: T.I32 },
	{ id: 4, key: "duration_millis", type: T.I64 },
	{ id: 5, key: "filesize_bytes", type: T.I64 },
	{ id: 6, key: "filename", type: T.STRING },
	{ id: 7, key: "attachment_id", type: T.STRING },
	{ id: 8, key: "legacy_media_url_https", type: T.STRING },
	{ id: 9, key: "legacy_media_preview_url", type: T.STRING },
);

mediaDimensionsSchema.push(
	{ id: 1, key: "width", type: T.I64 },
	{ id: 2, key: "height", type: T.I64 },
);

memberAccountDeleteEventSchema.push(
	{ id: 1, key: "member_id", type: T.STRING },
);

mentionRichTextContentSchema.push(
);

messageSchema.push(
	{ id: 1, key: "messageEvent", type: T.STRUCT, schema: messageEventSchema },
	{ id: 2, key: "messageInstruction", type: T.STRUCT, schema: messageInstructionSchema },
	{ id: 3, key: "batchedMessageEvents", type: T.STRUCT, schema: batchedMessageEventsSchema },
);

messageAttachmentSchema.push(
	{ id: 1, key: "media", type: T.STRUCT, schema: mediaAttachmentSchema },
	{ id: 2, key: "post", type: T.STRUCT, schema: postAttachmentSchema },
	{ id: 3, key: "url", type: T.STRUCT, schema: urlAttachmentSchema },
	{ id: 4, key: "unified_card", type: T.STRUCT, schema: unifiedCardAttachmentSchema },
	{ id: 5, key: "money", type: T.STRUCT, schema: moneyAttachmentSchema },
);

messageContentsSchema.push(
	{ id: 1, key: "message_text", type: T.STRING },
	{ id: 2, key: "entities", type: T.LIST, elemType: T.STRUCT, elemSchema: richTextEntitySchema },
	{ id: 3, key: "attachments", type: T.LIST, elemType: T.STRUCT, elemSchema: messageAttachmentSchema },
	{ id: 4, key: "replying_to_preview", type: T.STRUCT, schema: replyingToPreviewSchema },
	{ id: 6, key: "forwarded_message", type: T.STRUCT, schema: forwardedMessageSchema },
	{ id: 7, key: "sent_from", type: T.I32 },
	{ id: 8, key: "quick_reply", type: T.STRUCT, schema: quickReplySchema },
	{ id: 9, key: "ctas", type: T.LIST, elemType: T.STRUCT, elemSchema: callToActionSchema },
);

messageCreateEventSchema.push(
	{ id: 100, key: "contents", type: T.STRING },
	{ id: 101, key: "conversation_key_version", type: T.STRING },
	{ id: 102, key: "should_notify", type: T.BOOL },
	{ id: 103, key: "ttl_msec", type: T.I64 },
	{ id: 104, key: "delivered_at_msec", type: T.I64 },
	{ id: 105, key: "is_pending_public_key", type: T.BOOL },
	{ id: 106, key: "priority", type: T.I32 },
	{ id: 107, key: "additional_action_list", type: T.LIST, elemType: T.I32 },
);

messageDeleteEventSchema.push(
	{ id: 1, key: "sequence_ids", type: T.LIST, elemType: T.STRING },
	{ id: 2, key: "delete_message_action", type: T.I32 },
);

messageDurationChangeSchema.push(
	{ id: 1, key: "ttl_msec", type: T.I64 },
);

messageDurationRemoveSchema.push(
	{ id: 1, key: "current_ttl_msec", type: T.I64 },
);

messageEditSchema.push(
	{ id: 1, key: "message_sequence_id", type: T.STRING },
	{ id: 2, key: "updated_text", type: T.STRING },
	{ id: 3, key: "entities", type: T.LIST, elemType: T.STRUCT, elemSchema: richTextEntitySchema },
);

messageEntryContentsSchema.push(
	{ id: 2, key: "reaction_add", type: T.STRUCT, schema: messageReactionAddSchema },
	{ id: 3, key: "reaction_remove", type: T.STRUCT, schema: messageReactionRemoveSchema },
	{ id: 4, key: "message_edit", type: T.STRUCT, schema: messageEditSchema },
	{ id: 5, key: "mark_conversation_read", type: T.STRUCT, schema: markConversationReadSchema },
	{ id: 6, key: "mark_conversation_unread", type: T.STRUCT, schema: markConversationUnreadSchema },
	{ id: 7, key: "pin_conversation", type: T.STRUCT, schema: pinConversationSchema },
	{ id: 8, key: "unpin_conversation", type: T.STRUCT, schema: unpinConversationSchema },
	{ id: 9, key: "screen_capture_detected", type: T.STRUCT, schema: screenCaptureDetectedSchema },
	{ id: 10, key: "av_call_ended", type: T.STRUCT, schema: aVCallEndedSchema },
	{ id: 11, key: "av_call_missed", type: T.STRUCT, schema: aVCallMissedSchema },
	{ id: 12, key: "draft_message", type: T.STRUCT, schema: draftMessageSchema },
	{ id: 13, key: "accept_message_request", type: T.STRUCT, schema: acceptMessageRequestSchema },
	{ id: 14, key: "nickname_message", type: T.STRUCT, schema: nicknameMessageSchema },
	{ id: 15, key: "set_verified_status", type: T.STRUCT, schema: setVerifiedStatusSchema },
	{ id: 16, key: "av_call_started", type: T.STRUCT, schema: aVCallStartedSchema },
);

messageEntryHolderSchema.push(
	{ id: 1, key: "contents", type: T.STRUCT, schema: messageEntryContentsSchema },
);

messageEventSchema.push(
	{ id: 1, key: "sequence_id", type: T.STRING },
	{ id: 2, key: "message_id", type: T.STRING },
	{ id: 3, key: "sender_id", type: T.STRING },
	{ id: 4, key: "conversation_id", type: T.STRING },
	{ id: 5, key: "conversation_token", type: T.STRING },
	{ id: 6, key: "created_at_msec", type: T.STRING },
	{ id: 7, key: "detail", type: T.STRUCT, schema: messageEventDetailSchema },
	{ id: 8, key: "relay_source", type: T.I32 },
	{ id: 9, key: "message_event_signature", type: T.STRUCT, schema: messageEventSignatureSchema },
	{ id: 10, key: "previous_sequence_id", type: T.STRING },
	{ id: 11, key: "is_trusted", type: T.BOOL },
);

messageEventDetailSchema.push(
	{ id: 1, key: "messageCreateEvent", type: T.STRUCT, schema: messageCreateEventSchema },
	{ id: 3, key: "conversationKeyChangeEvent", type: T.STRUCT, schema: conversationKeyChangeEventSchema },
	{ id: 4, key: "groupChangeEvent", type: T.STRUCT, schema: groupChangeEventSchema },
	{ id: 5, key: "messageFailureEvent", type: T.STRUCT, schema: messageFailureEventSchema },
	{ id: 6, key: "messageTypingEvent", type: T.STRUCT, schema: messageTypingEventSchema },
	{ id: 7, key: "messageDeleteEvent", type: T.STRUCT, schema: messageDeleteEventSchema },
	{ id: 8, key: "conversationDeleteEvent", type: T.STRUCT, schema: conversationDeleteEventSchema },
	{ id: 9, key: "conversationMetadataChangeEvent", type: T.STRUCT, schema: conversationMetadataChangeEventSchema },
	{ id: 10, key: "grokSearchResponseEvent", type: T.STRUCT, schema: grokSearchResponseEventSchema },
	{ id: 11, key: "requestForEncryptedResendEvent", type: T.STRUCT, schema: requestForEncryptedResendEventSchema },
	{ id: 12, key: "markConversationReadEvent", type: T.STRUCT, schema: markConversationReadEventSchema },
	{ id: 13, key: "markConversationUnreadEvent", type: T.STRUCT, schema: markConversationUnreadEventSchema },
	{ id: 14, key: "memberAccountDeleteEvent", type: T.STRUCT, schema: memberAccountDeleteEventSchema },
);

messageEventRelaySourceSchema.push(
);

messageEventSignatureSchema.push(
	{ id: 1, key: "signature", type: T.STRING },
	{ id: 2, key: "public_key_version", type: T.STRING },
	{ id: 3, key: "signature_version", type: T.STRING },
	{ id: 4, key: "signing_public_key", type: T.STRING },
);

messageFailureEventSchema.push(
	{ id: 1, key: "failure_type", type: T.I32 },
);

messageInstructionSchema.push(
	{ id: 1, key: "pullMessagesInstruction", type: T.STRUCT, schema: pullMessagesInstructionSchema },
	{ id: 2, key: "keepAliveInstruction", type: T.STRUCT, schema: keepAliveInstructionSchema },
	{ id: 3, key: "pullMessagesFinishedInstruction", type: T.STRUCT, schema: pullMessagesFinishedInstructionSchema },
	{ id: 4, key: "pinReminderInstruction", type: T.STRUCT, schema: pinReminderInstructionSchema },
	{ id: 5, key: "switchToHybridPullInstruction", type: T.STRUCT, schema: switchToHybridPullInstructionSchema },
	{ id: 6, key: "displayTemporaryPasscodeInstruction", type: T.STRUCT, schema: displayTemporaryPasscodeInstructionSchema },
);

messageReactionAddSchema.push(
	{ id: 1, key: "message_sequence_id", type: T.STRING },
	{ id: 2, key: "emoji", type: T.STRING },
);

messageReactionRemoveSchema.push(
	{ id: 1, key: "message_sequence_id", type: T.STRING },
	{ id: 2, key: "emoji", type: T.STRING },
);

messageTypingEventSchema.push(
	{ id: 1, key: "conversation_id", type: T.STRING },
);

moneyAttachmentSchema.push(
	{ id: 1, key: "fallbackText", type: T.STRING },
	{ id: 2, key: "payload", type: T.STRING },
);

muteConversationSchema.push(
	{ id: 1, key: "muted_conversation_ids", type: T.LIST, elemType: T.STRING },
);

nicknameMessageSchema.push(
	{ id: 1, key: "user_id", type: T.I64 },
	{ id: 2, key: "nickname_text", type: T.STRING },
);

parentNodeSchema.push(
	{ id: 1, key: "subtree_encryption_public_key", type: T.STRING },
	{ id: 2, key: "parent_hash", type: T.STRING },
);

phoneNumberRichTextContentSchema.push(
);

pinConversationSchema.push(
	{ id: 1, key: "conversation_id", type: T.STRING },
);

pinReminderInstructionSchema.push(
	{ id: 1, key: "should_register", type: T.BOOL },
	{ id: 2, key: "should_generate", type: T.BOOL },
);

postAttachmentSchema.push(
	{ id: 1, key: "rest_id", type: T.STRING },
	{ id: 2, key: "post_url", type: T.STRING },
	{ id: 3, key: "attachment_id", type: T.STRING },
);

pullMessagePageDetailsSchema.push(
	{ id: 3, key: "min_sequence_id", type: T.STRING },
	{ id: 4, key: "max_sequence_id", type: T.STRING },
	{ id: 7, key: "is_batched_pull", type: T.BOOL },
);

pullMessagesFinishedInstructionSchema.push(
	{ id: 1, key: "finished_pull", type: T.BOOL },
	{ id: 2, key: "sequence_continue", type: T.STRING },
	{ id: 3, key: "pull_message_page_details", type: T.STRUCT, schema: pullMessagePageDetailsSchema },
);

pullMessagesInstructionSchema.push(
	{ id: 1, key: "sequence_start", type: T.STRING },
	{ id: 2, key: "sender_id", type: T.STRING },
	{ id: 6, key: "is_batched_pull", type: T.BOOL },
);

quickReplySchema.push(
	{ id: 1, key: "request", type: T.STRUCT, schema: quickReplyRequestSchema },
	{ id: 2, key: "response", type: T.STRUCT, schema: quickReplyResponseSchema },
);

quickReplyOptionSchema.push(
	{ id: 1, key: "id", type: T.STRING },
	{ id: 2, key: "label", type: T.STRING },
	{ id: 3, key: "metadata", type: T.STRING },
	{ id: 4, key: "description", type: T.STRING },
);

quickReplyOptionsRequestSchema.push(
	{ id: 1, key: "id", type: T.STRING },
	{ id: 2, key: "options", type: T.LIST, elemType: T.STRUCT, elemSchema: quickReplyOptionSchema },
);

quickReplyOptionsResponseSchema.push(
	{ id: 1, key: "request_id", type: T.STRING },
	{ id: 2, key: "metadata", type: T.STRING },
	{ id: 3, key: "selected_option_id", type: T.STRING },
);

quickReplyRequestSchema.push(
	{ id: 1, key: "options", type: T.STRUCT, schema: quickReplyOptionsRequestSchema },
);

quickReplyResponseSchema.push(
	{ id: 1, key: "options", type: T.STRUCT, schema: quickReplyOptionsResponseSchema },
);

ratchetTreeSchema.push(
	{ id: 1, key: "leaves", type: T.LIST, elemType: T.STRUCT, elemSchema: ratchetTreeLeafSchema },
	{ id: 2, key: "parents", type: T.LIST, elemType: T.STRUCT, elemSchema: ratchetTreeParentSchema },
);

ratchetTreeLeafSchema.push(
	{ id: 1, key: "empty", type: T.STRUCT, schema: emptyNodeSchema },
	{ id: 2, key: "leaf", type: T.STRUCT, schema: leafNodeSchema },
);

ratchetTreeParentSchema.push(
	{ id: 1, key: "empty", type: T.STRUCT, schema: emptyNodeSchema },
	{ id: 2, key: "parent", type: T.STRUCT, schema: parentNodeSchema },
);

replyingToPreviewSchema.push(
	{ id: 1, key: "sender_id", type: T.I64 },
	{ id: 2, key: "message_text", type: T.STRING },
	{ id: 3, key: "entities", type: T.LIST, elemType: T.STRUCT, elemSchema: richTextEntitySchema },
	{ id: 4, key: "attachments", type: T.LIST, elemType: T.STRUCT, elemSchema: messageAttachmentSchema },
	{ id: 5, key: "sender_display_name", type: T.STRING },
	{ id: 6, key: "replying_to_message_sequence_id", type: T.STRING },
	{ id: 7, key: "replying_to_message_id", type: T.STRING },
);

requestForEncryptedResendEventSchema.push(
	{ id: 1, key: "min_sequence_id", type: T.STRING },
	{ id: 2, key: "max_sequence_id", type: T.STRING },
);

richTextContentSchema.push(
	{ id: 1, key: "hashtag", type: T.STRUCT, schema: hashtagRichTextContentSchema },
	{ id: 2, key: "cashtag", type: T.STRUCT, schema: cashtagRichTextContentSchema },
	{ id: 3, key: "mention", type: T.STRUCT, schema: mentionRichTextContentSchema },
	{ id: 4, key: "url", type: T.STRUCT, schema: urlRichTextContentSchema },
	{ id: 5, key: "email", type: T.STRUCT, schema: emailRichTextContentSchema },
	{ id: 7, key: "phoneNumber", type: T.STRUCT, schema: phoneNumberRichTextContentSchema },
);

richTextEntitySchema.push(
	{ id: 1, key: "start_index", type: T.I32 },
	{ id: 2, key: "end_index", type: T.I32 },
	{ id: 3, key: "content", type: T.STRUCT, schema: richTextContentSchema },
);

screenCaptureDetectedSchema.push(
	{ id: 1, key: "type", type: T.I32 },
);

setVerifiedStatusSchema.push(
	{ id: 1, key: "user_id", type: T.I64 },
	{ id: 2, key: "verified_status", type: T.BOOL },
);

storedGroupStateSchema.push(
	{ id: 1, key: "keypairs", type: T.LIST, elemType: T.STRUCT, elemSchema: maybeKeypairSchema },
	{ id: 2, key: "ratchet_tree", type: T.STRUCT, schema: ratchetTreeSchema },
);

storedKeypairSchema.push(
	{ id: 1, key: "public_key", type: T.STRING },
	{ id: 2, key: "private_key", type: T.STRING },
);

switchToHybridPullInstructionSchema.push(
	{ id: 1, key: "requesting_user_agent", type: T.STRING },
);

unifiedCardAttachmentSchema.push(
	{ id: 1, key: "url", type: T.STRING },
	{ id: 2, key: "attachment_id", type: T.STRING },
);

unmuteConversationSchema.push(
	{ id: 1, key: "unmuted_conversation_ids", type: T.LIST, elemType: T.STRING },
);

unpinConversationSchema.push(
	{ id: 1, key: "conversation_id", type: T.STRING },
);

updatePathNodeSchema.push(
	{ id: 1, key: "encrypted_secrets", type: T.LIST, elemType: T.STRING },
	{ id: 2, key: "encrypted_private_key", type: T.STRING },
);

urlAttachmentSchema.push(
	{ id: 1, key: "url", type: T.STRING },
	{ id: 2, key: "banner_image_media_hash_key", type: T.STRUCT, schema: urlAttachmentImageSchema },
	{ id: 3, key: "favicon_image_media_hash_key", type: T.STRUCT, schema: urlAttachmentImageSchema },
	{ id: 4, key: "display_title", type: T.STRING },
	{ id: 5, key: "attachment_id", type: T.STRING },
);

urlAttachmentImageSchema.push(
	{ id: 1, key: "media_hash_key", type: T.STRING },
	{ id: 2, key: "filesize_bytes", type: T.I64 },
	{ id: 3, key: "filename", type: T.STRING },
	{ id: 4, key: "dimensions", type: T.STRUCT, schema: mediaDimensionsSchema },
);

urlRichTextContentSchema.push(
);
export const xchatRootSchema = messageSchema;

// ---- Websocket pipeline ----
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
		const message = decoder.readStruct(xchatRootSchema);

		if (!message) {
			const genericDecoder = new Decoder(new Uint8Array(d));

			onEvent?.(genericDecoder.readStructGeneric());
		} else {
			onEvent?.(message);
		}

		const eventJson = JSON.stringify(message, jsonReplacer, 2);

		console.log(message, utf8ToBase64(eventJson));

		const cipherField = message?.messageEvent?.detail?.messageCreateEvent?.contents;
		let ciphertext: Uint8Array | null = null;
		if (cipherField instanceof Uint8Array) {
			ciphertext = cipherField;
		} else if (typeof cipherField === "string") {
			const trimmed = cipherField.trim();
			if (!trimmed) {
				onDecrypted?.("No contents present in messageCreateEvent.");
				return;
			}
			try {
				ciphertext = base64ToUint8Array(trimmed);
			} catch {
				onDecrypted?.(trimmed);
				return;
			}
		}

		if (!ciphertext) {
			onDecrypted?.("No ciphertext present in messageCreateEvent.contents.");
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
