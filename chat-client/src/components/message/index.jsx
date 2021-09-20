import React from "react";
import "./style.css";
import { format } from "timeago.js";

export default function index({ message, isOwn }) {
  return (
    <div className={isOwn ? "message isOwn" : "message"}>
      <div className="messageTop">
        <p className="messageText">{message?.message}</p>
      </div>
      <div className="messageBottom">{format(message?.timeStamp)}</div>
    </div>
  );
}
