import React, { useState, useEffect, useContext, useRef, useMemo } from "react";
import { Redirect } from "react-router-dom";
import { notification, Spin } from "antd";
import { LoadingOutlined } from "@ant-design/icons";
import "./style.css";
import Channels from "../../components/channels/index";
import Message from "../../components/message/index";
import { Context } from "../../Store.js";
import useWebSocket, { ReadyState } from 'react-use-websocket';

const {
  getUserMessages,
  getChannel,
  updateChannelStatus,
  createChannel,
  createMessage,
} = require("../../services/index");

export default function Index() {
  // const [inActiveMessages, setInActiveMessages] = useState([]);
  const [state, dispatch] = useContext(Context);
  const [myChannels, setMyChannels] = useState([]);
  const [socketUrl, setSocketUrl] = useState(process.env.REACT_APP_SOCKET_SERVER_URL);
  const [selectedChannel, setSelectedChannel] = useState({});
  const [selectedMessages, setSelectedMessages] = useState([]);
  const [arrivedMessage, setArrivedMessage] = useState({});
  const [newMessage, setNewMessage] = useState("");
  const [isNewMessage, setIsNewMessage] = useState(false);
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

  useEffect(() => {
    scrollRef.current?.scrollIntoView({behaviour: "smooth"});
  }, [selectedMessages]);

  useMemo(() => {
    if (lastMessage) {
      setArrivedMessage(JSON.parse(lastMessage.data));
    }
  }, [lastMessage]);

  const fetchMyMessages = async () => {
    setIsLoading(true);
    try {
      const requestBody = {
        email: state.user.email,
        token: state.user.token,
      };
      const resp = await getUserMessages(requestBody);
      setMyChannels(resp.data);
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

  useEffect(() => {
    const abortController = new AbortController();
    setSelectedMessages([]);
    fetchMyMessages();
    return () => {
      abortController.abort();
    };
  }, [state.user.email, state.user.role]);

  useEffect(() => {
    const abortController = new AbortController();
    const fetchChannel = async (channel) => {
      setIsLoading(true);
      try {
        if (
          selectedChannel?.id !== "" &&
          typeof selectedChannel?.id !== "undefined"
        ) {
          const data = await getChannel(selectedChannel?.id, state.user.token);
          dispatch({ type: "SET_SELECTED_CHANNEL", payload: data });
          setSelectedMessages(data.messages);
          setIsNewMessage(false);
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
    if (newMessage.trim() === "") {
      notification.error({
        message: "Could not send message",
        description: "Please type in a message",
        placement: "topRight",
        duration: 1.5,
        onClose: () => {
          setRedirect(null);
          setIsLoading(false);
        },
      });
      return;
    }
    if (!isNewMessage && selectedChannel?.id === undefined) {
      notification.error({
        message: "Could not send message",
        description: "Please select a conversation or click on new message to send",
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

    if (selectedChannel && !isNewMessage) {
      if (selectedChannel?.status === "INACTIVE") {
        const requestBody = {
          id: selectedChannel?.id,
          status: "ACTIVE",
          updatedBy: state.user.email,
        };
        try {
          await updateChannelStatus(requestBody);
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
      }
    } else {
      const channelRequest = {
        userEmail: state.user.email,
      };

      try {
        const addedChannel = await createChannel(
          channelRequest,
          state.user.token
        );
        if (newMessage.trim() !== "") {
          const messageRequest = {
            channelId: addedChannel.id,
            message: newMessage,
            messageFrom: state.user.email,
          };
          await createMessage(messageRequest, state.user.token);
          setSelectedChannel(addedChannel);
          setNewMessage("");
          setIsNewMessage(false);
          fetchMyMessages();
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
    }

    sendMessage(JSON.stringify(message));
  };

  const handleNewMessage = async (e) => {
    notification.success({
      message: "Add New Message!",
      description:
        "Type in the message box and click send to create a new message",
      placement: "topRight",
      duration: 2.0,
    });
    e.preventDefault();
    setIsNewMessage(true);
    setSelectedChannel({});
    setSelectedMessages([]);
  };

  if (state.user.role === "ADMIN") {
    return <Redirect to={"/message-dashboard"} />;
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
                  <Channels key={c.id} channel={c} />
                </div>
              ))}
              <button className="newMessageButton" onClick={handleNewMessage}>
                New Message
              </button>
            </div>
          </div>
          <div className="messageBox">
            <div className="messageBoxWrapper">
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
                        indicator={<LoadingOutlined style={{ fontSize: 24 }} />}
                      />
                    ) : (
                      "Send"
                    )}
                  </button>
                </div>
              </>
            </div>
          </div>
        </div>
      </React.Fragment>
    );
  }
}
