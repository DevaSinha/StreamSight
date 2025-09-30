"use client";
import React, { useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";

const WS_URL = "ws://localhost:8080/ws/alerts"; // adjust as needed

type AlertType = {
  cameraName?: string;
  description?: string;
  timestamp?: string;
  [key: string]: any;
};

export default function AlertsFeed() {
  const [alerts, setAlerts] = useState<AlertType[]>([]);
  const { lastMessage, readyState } = useWebSocket(WS_URL, {
    shouldReconnect: () => true,
  });

  useEffect(() => {
    if (lastMessage !== null) {
      try {
        const data = JSON.parse(lastMessage.data);
        setAlerts((prev) => [data, ...prev]);
      } catch {
        setAlerts((prev) => [{ raw: lastMessage?.data }, ...prev]);
      }
    }
  }, [lastMessage]);

  return (
    <div>
      <h2>Live Alerts</h2>
      <div>
        {alerts.map((alert, idx) => (
          <div
            key={idx}
            style={{
              padding: "6px 0",
              marginBottom: "6px",
              borderBottom: "1px solid #ddd",
            }}
          >
            {alert.cameraName ? (
              <>
                <div>
                  <strong>{alert.cameraName}</strong>: {alert.description}
                </div>
                <div style={{ color: "#888", fontSize: "0.9em" }}>
                  {alert.timestamp}
                </div>
              </>
            ) : (
              <pre>{JSON.stringify(alert)}</pre>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
