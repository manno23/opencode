// src/template.ts
function renderHTML(content, title = "OpenCode Share") {
  return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #fff;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }

        .error {
            background: #fee;
            border: 1px solid #fcc;
            padding: 20px;
            border-radius: 8px;
            color: #c33;
        }

        .header {
            border-bottom: 1px solid #eee;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }

        .header h1 {
            font-size: 24px;
            margin-bottom: 10px;
        }

        .meta {
            color: #666;
            font-size: 14px;
        }

        .message {
            margin-bottom: 30px;
            padding: 20px;
            border: 1px solid #eee;
            border-radius: 8px;
        }

        .message.user {
            background: #f8f9ff;
            border-color: #e1e5ff;
        }

        .message.assistant {
            background: #f8fff8;
            border-color: #e1ffe1;
        }

        .role {
            font-weight: bold;
            margin-bottom: 10px;
            text-transform: capitalize;
        }

        .part {
            margin-bottom: 15px;
            padding: 10px;
            background: white;
            border-radius: 4px;
        }

        pre {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 4px;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 14px;
        }

        code {
            background: #f0f0f0;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', monospace;
        }
    </style>
</head>
<body>
    <div class="container">
        ${content}
    </div>

    <script>
        // Add some basic interactivity for future TSX components
        console.log('OpenCode Web Frontend - Initial Load');

        // Future: This is where we'll add Solid.js hydration
        // window.__INITIAL_DATA__ = ${JSON.stringify({})};
    </script>
</body>
</html>`;
}

// src/index.tsx
var index_default = {
  async fetch(request, env, _ctx) {
    const url = new URL(request.url);
    if (url.pathname === "/") {
      return handleHome();
    }
    if (url.pathname === "/health") {
      return new Response("OK", { status: 200 });
    }
    const shareMatch = url.pathname.match(/^\/s\/(.+)$/);
    if (shareMatch) {
      const shareId = shareMatch[1];
      return handleShare(shareId, env);
    }
    if (url.pathname.startsWith("/static/")) {
      return handleStatic(request);
    }
    return new Response("Not Found", { status: 404 });
  }
};
async function handleHome() {
  const content = `
    <div class="header">
      <h1>OpenCode Web Frontend</h1>
      <p class="meta">Cloudflare Worker-based session sharing</p>
    </div>

    <div class="message">
      <h2>\u{1F680} Worker is running!</h2>
      <p>This is the new Cloudflare Worker frontend for OpenCode session sharing.</p>
      <p>Try accessing a shared session at <code>/s/{shareId}</code></p>

      <h3>Features:</h3>
      <ul>
        <li>\u2705 TypeScript support</li>
        <li>\u2705 TSX component rendering (coming next)</li>
        <li>\u2705 API integration ready</li>
        <li>\u2705 Real-time WebSocket connections (coming next)</li>
      </ul>
    </div>
  `;
  const html = renderHTML(content, "OpenCode Web Frontend");
  return new Response(html, {
    headers: { "Content-Type": "text/html" }
  });
}
async function handleShare(shareId, env) {
  try {
    const apiUrl = `${env.API_URL}/share_data?id=${shareId}`;
    console.log("Fetching session data from:", apiUrl);
    const response = await fetch(apiUrl);
    if (!response.ok) {
      if (response.status === 404) {
        return renderNotFound(shareId);
      }
      throw new Error(`API responded with ${response.status}: ${response.statusText}`);
    }
    const data = await response.json();
    if (!data.info) {
      return renderNotFound(shareId);
    }
    return renderSession(data, shareId, env);
  } catch (error) {
    console.error("Error fetching session:", error);
    return renderError(shareId, error);
  }
}
function renderNotFound(shareId) {
  const content = `
    <div class="error">
      <h2>Session Not Found</h2>
      <p>The shared session <code>${shareId}</code> could not be found.</p>
      <p>It may have been deleted or the link may be incorrect.</p>
    </div>
  `;
  const html = renderHTML(content, "Session Not Found");
  return new Response(html, {
    status: 404,
    headers: { "Content-Type": "text/html" }
  });
}
function renderError(shareId, error) {
  const errorMessage = error instanceof Error ? error.message : "Unknown error";
  const content = `
    <div class="error">
      <h2>Error Loading Session</h2>
      <p>Failed to load shared session <code>${shareId}</code>.</p>
      <p><strong>Error:</strong> ${errorMessage}</p>
      <p>Please try again later or contact support if the issue persists.</p>
    </div>
  `;
  const html = renderHTML(content, "Error Loading Session");
  return new Response(html, {
    status: 500,
    headers: { "Content-Type": "text/html" }
  });
}
function renderSession(data, _shareId, _env) {
  const content = renderShareComponent(data, "connected");
  const html = renderHTML(content, data.info.title);
  return new Response(html, {
    headers: { "Content-Type": "text/html" }
  });
}
function renderShareComponent(data, status, error) {
  const messages = Object.values(data.messages).sort((a, b) => a.id?.localeCompare(b.id));
  const totals = messages.reduce(
    (acc, msg) => {
      if (msg.role === "assistant") {
        acc.cost += msg.cost || 0;
        acc.tokens.input += msg.tokens?.input || 0;
        acc.tokens.output += msg.tokens?.output || 0;
        acc.tokens.reasoning += msg.tokens?.reasoning || 0;
      }
      return acc;
    },
    {
      cost: 0,
      tokens: { input: 0, output: 0, reasoning: 0 }
    }
  );
  const formatDateTime = (timestamp) => new Date(timestamp).toLocaleString();
  const formatCost = (cost) => `$${cost.toFixed(2)}`;
  const formatTokens = (tokens) => tokens.toLocaleString();
  return `
    <main data-component="share">
      <style>
        [data-component="share"] {
          display: flex;
          flex-direction: column;
          gap: 2.5rem;
          padding: 1.5rem;
          max-width: 1200px;
          margin: 0 auto;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
          line-height: 1.6;
          color: #1f2937;
        }

        @media (prefers-color-scheme: dark) {
          [data-component="share"] {
            background-color: #111827;
            color: #f9fafb;
          }
        }

        .header {
          display: flex;
          flex-direction: column;
          gap: 1rem;
        }

        .header-title {
          font-size: 2.75rem;
          font-weight: 500;
          line-height: 1.2;
          letter-spacing: -0.05em;
          margin: 0;
          color: #1f2937;
        }

        @media (prefers-color-scheme: dark) {
          .header-title {
            color: #f9fafb;
          }
        }

        .header-stats {
          list-style: none;
          padding: 0;
          margin: 0;
          display: flex;
          gap: 0.5rem 0.875rem;
          flex-wrap: wrap;
        }

        .header-stats li {
          display: flex;
          align-items: center;
          gap: 0.3125rem;
          font-size: 0.875rem;
          color: #374151;
        }

        @media (prefers-color-scheme: dark) {
          .header-stats li {
            color: #d1d5db;
          }
        }

        .header-stats svg {
          color: #6b7280;
          opacity: 0.85;
        }

        .header-time {
          color: #6b7280;
          font-size: 0.875rem;
        }

        @media (prefers-color-scheme: dark) {
          .header-time {
            color: #9ca3af;
          }
        }

        .parts {
          display: flex;
          flex-direction: column;
          gap: 0.625rem;
        }

        .part {
          display: flex;
          gap: 0.625rem;
        }

        .part-decoration {
          flex: 0 0 auto;
          display: flex;
          flex-direction: column;
          gap: 0.625rem;
          align-items: center;
          justify-content: flex-start;
        }

        .part-anchor {
          width: 18px;
          opacity: 0.65;
        }

        .part-anchor svg {
          color: #6b7280;
          display: block;
        }

        .part-bar {
          width: 3px;
          min-height: 2rem;
          height: 100%;
          border-radius: 1px;
          background-color: #e5e7eb;
        }

        @media (prefers-color-scheme: dark) {
          .part-bar {
            background-color: #374151;
          }
        }

        .part-content {
          flex: 1 1 auto;
          min-width: 0;
          display: flex;
          flex-direction: column;
          gap: 1rem;
        }

        .user-text, .assistant-text {
          padding: 0.75rem;
          border-radius: 0.375rem;
          max-width: 40rem;
        }

        .user-text {
          background-color: #f9fafb;
          border: 1px solid #e5e7eb;
        }

        .assistant-text {
          background-color: #f0f9ff;
          border: 1px solid #0ea5e9;
        }

        @media (prefers-color-scheme: dark) {
          .user-text {
            background-color: #1f2937;
            border: 1px solid #374151;
          }

          .assistant-text {
            background-color: #0f172a;
            border: 1px solid #1e40af;
          }
        }

        .user-text pre, .assistant-text pre {
          margin: 0;
          font-size: 0.875rem;
          line-height: 1.5;
          white-space: pre-wrap;
          word-wrap: break-word;
          color: inherit;
        }

        .tool {
          max-width: 40rem;
        }

        .tool-title {
          display: flex;
          align-items: flex-start;
          gap: 0.375rem;
          font-size: 0.875rem;
          color: #6b7280;
          margin-bottom: 0.5rem;
        }

        .tool-name {
          font-weight: 500;
        }

        .tool-target {
          color: #1f2937;
          word-break: break-all;
          font-weight: 500;
        }

        @media (prefers-color-scheme: dark) {
          .tool-target {
            color: #f9fafb;
          }
        }

        .tool-result {
          background-color: #f9fafb;
          border: 1px solid #e5e7eb;
          border-radius: 0.25rem;
          padding: 0.75rem;
        }

        @media (prefers-color-scheme: dark) {
          .tool-result {
            background-color: #1f2937;
            border: 1px solid #374151;
          }
        }

        .tool-result pre {
          margin: 0;
          font-size: 0.75rem;
          line-height: 1.6;
          white-space: pre-wrap;
          word-break: break-word;
          color: inherit;
        }

        .summary {
          display: flex;
          gap: 0.625rem;
          margin-top: 2rem;
          padding-top: 1.5rem;
          border-top: 1px solid #e5e7eb;
        }

        @media (prefers-color-scheme: dark) {
          .summary {
            border-top: 1px solid #374151;
          }
        }

        .summary-decoration {
          flex: 0 0 auto;
          display: flex;
          align-items: center;
          justify-content: center;
          padding-top: 2px;
        }

        .summary-status {
          display: block;
          width: 14px;
          height: 14px;
          border-radius: 50%;
          background-color: #d1d5db;
        }

        .summary-status[data-status="connected"] {
          background-color: #10b981;
        }

        .summary-status[data-status="loading"] {
          background-color: #f59e0b;
        }

        .summary-status[data-status="error"] {
          background-color: #ef4444;
        }

        .summary-content {
          flex: 1 1 auto;
          display: flex;
          flex-direction: column;
          gap: 0.75rem;
        }

        .summary-text {
          font-size: 0.875rem;
          color: #6b7280;
          margin: 0;
        }

        .summary-stats {
          list-style: none;
          padding: 0;
          margin: 0;
          display: flex;
          gap: 0.5rem 1.5rem;
          flex-wrap: wrap;
        }

        .summary-stats li {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          font-size: 0.75rem;
          color: #374151;
        }

        @media (prefers-color-scheme: dark) {
          .summary-stats li {
            color: #d1d5db;
          }
        }

        .placeholder {
          color: #9ca3af;
        }

        @media (prefers-color-scheme: dark) {
          .placeholder {
            color: #6b7280;
          }
        }

        @media (max-width: 768px) {
          [data-component="share"] {
            padding: 1rem;
            gap: 2rem;
          }

          .header-title {
            font-size: 1.875rem;
          }

          .user-text, .assistant-text, .tool {
            max-width: 100%;
          }
        }
      </style>

      <div class="header">
        <h1 class="header-title">${escapeHtml(data.info.title)}</h1>
        <div>
          <ul class="header-stats">
            <li>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
              <span>v${data.info.version || "0.0.1"}</span>
            </li>
          </ul>
          <div class="header-time">
            ${formatDateTime(data.info.time.created)}
          </div>
        </div>
      </div>

      <div class="parts">
        ${messages.length > 0 ? messages.map((message) => renderMessageParts(message)).join("") : `<div class="part">
            <div class="part-decoration">
              <div class="part-anchor">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <div class="part-bar"></div>
            </div>
            <div class="part-content">
              <p>Waiting for messages...</p>
            </div>
          </div>`}

        <div class="summary">
          <div class="summary-decoration">
            <span 
              class="summary-status"
              data-status="${status}"
            ></span>
          </div>
          <div class="summary-content">
            <p class="summary-text">
              ${status === "loading" ? "Loading session..." : status === "connected" ? "Session loaded" : error || "Error loading session"}
            </p>
            <ul class="summary-stats">
              <li>
                <span>Cost</span>
                ${totals.cost > 0 ? `<span>${formatCost(totals.cost)}</span>` : `<span class="placeholder">\u2014</span>`}
              </li>
              <li>
                <span>Input Tokens</span>
                ${totals.tokens.input > 0 ? `<span>${formatTokens(totals.tokens.input)}</span>` : `<span class="placeholder">\u2014</span>`}
              </li>
              <li>
                <span>Output Tokens</span>
                ${totals.tokens.output > 0 ? `<span>${formatTokens(totals.tokens.output)}</span>` : `<span class="placeholder">\u2014</span>`}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </main>
  `;
}
function renderMessageParts(message) {
  const filteredParts = message.parts.filter((part) => {
    if (part.type === "step-start" || part.type === "snapshot" || part.type === "patch")
      return false;
    if (part.type === "step-finish") return false;
    if (part.type === "text" && part.synthetic === true) return false;
    if (part.type === "tool" && part.tool === "todoread") return false;
    if (part.type === "text" && !part.text) return false;
    if (part.type === "tool" && (part.state?.status === "pending" || part.state?.status === "running"))
      return false;
    return true;
  });
  return filteredParts.map(
    (part) => `
    <div class="part">
      <div class="part-decoration">
        <div class="part-anchor">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            ${message.role === "user" ? `<path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>` : `<path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>`}
          </svg>
        </div>
        <div class="part-bar"></div>
      </div>

      <div class="part-content">
        ${renderPartContent(part, message.role)}
      </div>
    </div>
  `
  ).join("");
}
function renderPartContent(part, messageRole) {
  switch (part.type) {
    case "text":
      if (!part.text) return "";
      const textClass = messageRole === "user" ? "user-text" : "assistant-text";
      return `<div class="${textClass}">
        <pre>${escapeHtml(part.text)}</pre>
      </div>`;
    case "tool":
      const toolName = part.tool || "unknown";
      let toolContent = `<div class="tool">
        <div class="tool-title">
          <span class="tool-name">${escapeHtml(toolName)}</span>`;
      if (part.state?.input && typeof part.state.input === "object") {
        if (part.state.input.filePath) {
          toolContent += `<span class="tool-target">${escapeHtml(part.state.input.filePath)}</span>`;
        } else if (part.state.input.command) {
          toolContent += `<span class="tool-target">${escapeHtml(part.state.input.command)}</span>`;
        } else if (part.state.input.pattern) {
          toolContent += `<span class="tool-target">"${escapeHtml(part.state.input.pattern)}"</span>`;
        }
      }
      toolContent += `</div>`;
      if (part.state?.output) {
        const output = typeof part.state.output === "string" ? part.state.output : JSON.stringify(part.state.output, null, 2);
        toolContent += `<div class="tool-result">
          <pre>${escapeHtml(output)}</pre>
        </div>`;
      }
      toolContent += `</div>`;
      return toolContent;
    case "file":
      return `<div class="file-attachment">
        \u{1F4C1} <strong>${escapeHtml(part.filename || "File")}</strong>
        ${part.mime ? `<span class="file-type">(${escapeHtml(part.mime)})</span>` : ""}
      </div>`;
    default:
      return `<div class="unknown-part">Unknown part type: ${escapeHtml(part.type)}</div>`;
  }
}
function escapeHtml(text) {
  return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;").replace(/\n/g, "<br>");
}
async function handleStatic(_request) {
  return new Response("Static assets not implemented yet", { status: 404 });
}
export {
  index_default as default
};
