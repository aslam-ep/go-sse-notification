let userID = "123";
let lastEventID = null;

function connectSSE() {
  const es = new EventSource("http://localhost:8081/notifications?userID=" + userID, {
    headers: lastEventID ? { "Last-Event-ID": lastEventID } : {}
  });

  es.onmessage = (event) => {
    console.log(`[SSE][${event.lastEventId || ""}] ${event.data}`);
    lastEventID = event.lastEventId;
  };

  es.onerror = (err) => {
    console.error("SSE connection error:", err);
    es.close();
    setTimeout(connectSSE, 2000); // reconnect
  };

  console.log("SSE connected for user", userID);
}

connectSSE();
