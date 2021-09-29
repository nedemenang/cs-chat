import React, { useState, useEffect, useContext, useRef, useMemo } from "react";
import { notification, Spin } from "antd";
import { LoadingOutlined } from "@ant-design/icons";
import { Redirect } from "react-router-dom";
import useWebSocket, { ReadyState } from 'react-use-websocket';
import "./style.css";
import Channels from "../../components/channels/index";
import Message from "../../components/message/index";
import { Context } from "../../Store.js";
const {
  getAdminMessages,
  getActiveMessages,
  getChannel,
  updateChannelStatus,
} = require("../../services/index");

export default function Index() {
  const [state, dispatch] = useContext(Context);
  const [socketUrl, setSocketUrl] = useState(process.env.REACT_APP_SOCKET_SERVER_URL);
  const [myChannels, setMyChannels] = useState([]);
  const [activeChannels, setActiveChannels] = useState([]);
  const [selectedChannel, setSelectedChannel] = useState({});
  const [selectedMessages, setSelectedMessages] = useState([]);
  const [arrivedMessage, setArrivedMessage] = useState({});
  const [newMessage, setNewMessage] = useState("");
  const [redirect, setRedirect] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const scrollRef = useRef();

  const {
    sendMessage,
    lastMessage,
    readyState,
  } = useWebSocket(socketUrl, {
    onOpen: () => console.log('opened'),
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
  });

  useEffect(() => {
    const abortController = new AbortController();
    if (
      typeof arrivedMessage?.message !== "undefined" &&
      arrivedMessage?.message !== "" &&
      arrivedMessage?.channelId === selectedChannel?.id &&
      arrivedMessage?.messageFrom !== state.user.email
    ) {
      setSelectedMessages((prev) => [...prev, arrivedMessage]);
    }
    return () => {
      abortController.abort();
    };
  }, [arrivedMessage]);
  
  useMemo(() => {
    if (lastMessage) {
      setArrivedMessage(JSON.parse(lastMessage.data));
    }
  }, [lastMessage]);

  useEffect(() => {
    scrollRef.current?.scrollIntoView({behaviour: "smooth"});
  }, [selectedMessages]);

  useEffect(() => {
    const abortController = new AbortController();
    const fetchMyMessages = async () => {
      setIsLoading(true);
      try {
        const resp = await getAdminMessages({
          token: state.user.token,
          email: state.user.email,
        });
        const respActive = await getActiveMessages(state.user.token);

        setMyChannels(resp.data);
        setActiveChannels(respActive.data);
      } catch (error) {
        notification.error({
          message: "Error fetching messages!",
          description: error.message,
          placement: "topRight",
          duration: 1.5,
          onClose: () => {
            setRedirect(null);
            setIsLoading(false);
          },
        });
      }
      setIsLoading(false);
    };
    fetchMyMessages();
    return () => {
      abortController.abort();
    };
  }, []);

  useEffect(() => {
    const abortController = new AbortController();
    const fetchChannel = async () => {
      setIsLoading(true);
      try {
        if (selectedChannel?.id !== "" && typeof selectedChannel?.id !== "undefined") {
          const data = await getChannel(selectedChannel.id, state.user.token);
          dispatch({ type: "SET_SELECTED_CHANNEL", payload: data });
          setSelectedMessages(data.messages);
        }
      } catch (error) {
        notification.error({
          message: "An error occured!",
          description: error,
          placement: "topRight",
          duration: 2.0,
          onClose: () => {
            setRedirect(null);
            setIsLoading(false);
          },
        });
      }
      setIsLoading(false);
    };
    fetchChannel();
    return () => {
      abortController.abort();
    };
  }, [selectedChannel]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (selectedChannel?.id === undefined) {
      notification.error({
        message: "Could not send message",
        description: "Please select a conversation to send message",
        placement: "topRight",
        duration: 1.5,
        onClose: () => {
          setRedirect(null);
          setIsLoading(false);
        },
      });
      return;
    }
    const message = {
      channelId: selectedChannel?.id,
      message: newMessage,
      messageFrom: state.user.email,
    };
    setSelectedMessages([...selectedMessages, message]);
    setNewMessage("");
    if (selectedChannel?.currentStatus === "ACTIVE") {
      const requestBody = {
        id: selectedChannel?.id,
        status: "IN_PROGRESS",
        updatedBy: state.user.email,
        token: state.user.token,
      };
      await updateChannelStatus(requestBody);
    }

    sendMessage(JSON.stringify(message));
  };

  if (state.user.role === "USER") {
    return <Redirect to={"/dashboard"} />;
  } else {
    return (
      <React.Fragment>
        <div className="messenger">
          <div className="myMessages">
            <div className="myMessagesWrapper">
              <span>
                <h3>My Messages</h3>
              </span>
              {myChannels.map((c) => (
                <div onClick={() => setSelectedChannel(c)}>
                  <Channels
                    key={c.id}
                    channel={c}
                    isAdmin={state.user.role === "ADMIN"}
                  />
                </div>
              ))}
            </div>
          </div>
          <div className="messageBox">
            <div className="messageBoxWrapper">
              {selectedChannel ? (
                <>
                  <div className="messageBoxTop">
                    {selectedMessages.map((m) => (
                      <div ref={scrollRef}>
                      <Message
                        message={m}
                        isOwn={state.user.email === m?.messageFrom}
                      />
                      </div>
                    ))}
                  </div>
                  <div className="messageBoxBottom">
                    <textarea
                      required
                      className="messageBoxBottomTextarea"
                      placeholder="Type a message..."
                      onChange={(e) => setNewMessage(e.target.value)}
                      value={newMessage}
                    ></textarea>
                    <button disabled={readyState !== ReadyState.OPEN} className="chatSubmit" onClick={handleSubmit}>
                      {isLoading ? (
                        <Spin
                          indicator={
                            <LoadingOutlined style={{ fontSize: 24 }} />
                          }
                        />
                      ) : (
                        "Send"
                      )}
                    </button>
                  </div>
                </>
              ) : (
                <span className="noConversation">
                  Open a conversation to start
                </span>
              )}
            </div>
          </div>
            <div className="inactiveMessages">
              <div className="inactiveMessagesWrapper">
                <span>
                  <h3>Active Messages</h3>
                </span>
                {activeChannels.map((c) => (
                  <div onClick={() => setSelectedChannel(c)}>
                    <Channels key={c.id} channel={c} />
                  </div>
                ))}
              </div>
            </div>
        </div>
      </React.Fragment>
    );
  }
}
