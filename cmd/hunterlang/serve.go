package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"hunterlang/lang/eval"
	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
)

type serveState struct {
	mu     sync.Mutex
	last   string
	snap   eval.Snapshot
	hasRun bool
}

func runProgramCapture(src string) (out string, snap eval.Snapshot, err error) {
	// Capture stdout produced by the interpreter.
	origStdout := os.Stdout
	r, w, errPipe := os.Pipe()
	if errPipe != nil {
		return "", snap, errPipe
	}
	os.Stdout = w

	lx := lexer.New(src)
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		_ = w.Close()
		os.Stdout = origStdout
		_ = r.Close()
		return "", snap, err
	}

	interpreter := eval.NewInterpreter()
	runErr := interpreter.Run(prog)
	snap = interpreter.Snapshot()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	out = buf.String()
	if runErr != nil {
		return out, snap, runErr
	}
	return out, snap, nil
}

func runProgramCaptureWithEvent(src string, eventType string, right *string) (out string, snap eval.Snapshot, err error) {
	// Capture stdout produced by the interpreter.
	origStdout := os.Stdout
	r, w, errPipe := os.Pipe()
	if errPipe != nil {
		return "", snap, errPipe
	}
	os.Stdout = w

	lx := lexer.New(src)
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		_ = w.Close()
		os.Stdout = origStdout
		_ = r.Close()
		return "", snap, err
	}

	interpreter := eval.NewInterpreter()
	interpreter.SetEvent(eventType, right)
	runErr := interpreter.Run(prog)
	snap = interpreter.Snapshot()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	out = buf.String()
	if runErr != nil {
		return out, snap, runErr
	}
	return out, snap, nil
}

func startServe(addr string, filePath string, siteMode bool, siteFilePath string, st *serveState) error {
	mux := http.NewServeMux()

	siteHTML := htmlPage
	if siteMode {
		// Run sitefile once and capture stdout as the full HTML.
		b, err := os.ReadFile(siteFilePath)
		if err != nil {
			return fmt.Errorf("read sitefile: %w", err)
		}

		if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
			b = b[3:]
		}

		out, _, err := runProgramCapture(string(b))
		if err != nil {
			return fmt.Errorf("run sitefile: %w", err)
		}
		siteHTML = out
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, siteHTML)
	})

	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		b, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Strip UTF-8 BOM if present: EF BB BF
		if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
			b = b[3:]
		}
		src := string(b)

		out, snap, err := runProgramCapture(src)
		if err != nil {
			// Still return output we captured to aid debugging.
			st.mu.Lock()
			st.last = out
			st.snap = snap
			st.hasRun = true
			st.mu.Unlock()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		st.mu.Lock()
		st.last = out
		st.snap = snap
		st.hasRun = true
		st.mu.Unlock()

		resp := map[string]any{"ok": true}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		type eventPayload struct {
			Type  string  `json:"type"`
			Right *string `json:"right"`
		}

		var ep eventPayload
		if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		if ep.Type == "" {
			http.Error(w, "missing event type", http.StatusBadRequest)
			return
		}

		b, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Strip UTF-8 BOM if present: EF BB BF
		if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
			b = b[3:]
		}
		src := string(b)

		out, snap, runErr := runProgramCaptureWithEvent(src, ep.Type, ep.Right)
		if runErr != nil {
			st.mu.Lock()
			st.last = out
			st.snap = snap
			st.hasRun = true
			st.mu.Unlock()
			http.Error(w, runErr.Error(), http.StatusInternalServerError)
			return
		}

		st.mu.Lock()
		st.last = out
		st.snap = snap
		st.hasRun = true
		st.mu.Unlock()

		resp := map[string]any{"ok": true}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		st.mu.Lock()
		defer st.mu.Unlock()

		resp := map[string]any{
			"hasRun": st.hasRun,
			"output": st.last,
			"snapshot": map[string]any{
				"platform": st.snap.Platform,
				"env":      st.snap.Env,
				"objects":  st.snap.Objects,
			},
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("HE Dev Server running at:", addr)
	return http.ListenAndServe(addr, mux)
}

const htmlPage = `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>HE Runtime</title>
  <style>
    body { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace; margin: 16px; background: #fff; color: #000; }
    .row { display: flex; gap: 16px; flex-wrap: wrap; }
    .card { border: 3px solid #000; padding: 12px; box-shadow: 6px 6px 0 #000; background: #f8f8f8; }
    .card h2 { margin: 0 0 10px 0; font-size: 18px; }
    .col { flex: 1 1 320px; min-width: 320px; }
    textarea { width: 100%; height: 260px; font-family: inherit; font-size: 13px; }
    pre { white-space: pre-wrap; word-break: break-word; }
    button { font-family: inherit; font-weight: 700; border: 3px solid #000; padding: 8px 12px; background: #fff; cursor: pointer; box-shadow: 4px 4px 0 #000; }
    button:active { transform: translate(2px,2px); box-shadow: 2px 2px 0 #000; }
    .muted { opacity: 0.7; }
  </style>
</head>
<body>
  <div class="row">
    <div class="card col">
      <h2>Console</h2>
      <button onclick="run()">Run program</button>
      <div class="muted" style="margin-top:10px;" id="status">Idle.</div>
      <textarea id="console" readonly></textarea>
    </div>

    <div class="card col">
      <h2>Runtime Output</h2>
      <div class="muted" style="margin-bottom:8px;">(UI HTML printed by HE)</div>
      <iframe id="ui" style="width:100%; height:520px; border:3px solid #000; background:#fff;"></iframe>
    </div>

    <div class="card col">
      <h2>Runtime Inspector</h2>
      <div class="muted">Platform:</div>
      <pre id="platform">unknown</pre>
      <div class="muted">Env:</div>
      <pre id="env">{}</pre>
      <div class="muted">Objects:</div>
      <pre id="objects">{}</pre>
    </div>
  </div>

<script>
async function run() {
  const status = document.getElementById('status');
  status.textContent = 'Running...';
  try {
    const resp = await fetch('/run', { method: 'POST' });
    if (!resp.ok) {
      const txt = await resp.text();
      status.textContent = 'Error: ' + txt;
      await refresh();
      return;
    }
    status.textContent = 'Done.';
    await refresh();
  } catch (e) {
    status.textContent = 'Error: ' + e;
  }
}

async function heClick(label) {
  // label is used as the trigger "right" payload.
  try {
    await fetch('/event', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: 'click', right: label })
    });
  } finally {
    await refresh();
  }
}

async function refresh() {
  const r = await fetch('/snapshot');
  const data = await r.json();
  document.getElementById('console').value = data.output || '';

  const ui = document.getElementById('ui');
  ui.srcdoc = data.output || '';

  document.getElementById('platform').textContent = data.snapshot.platform || 'unknown';
  document.getElementById('env').textContent = JSON.stringify(data.snapshot.env || {}, null, 2);
  document.getElementById('objects').textContent = JSON.stringify(data.snapshot.objects || {}, null, 2);
}

refresh();
</script>
</body>
</html>
`
