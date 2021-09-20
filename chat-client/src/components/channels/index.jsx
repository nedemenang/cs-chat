import React from "react";
import moment from "moment";
import "./style.css";

export default function Index({ channel }) {
  return (
    
    <div className="channel">
      <span className="channelName">
        {`${channel?.userFullName}`}
      </span>
      <span className="channelTime">
        {`${moment(channel.createdAt).format("DD/MM/YYYY")}`}
      </span>
    </div>
  );
}
