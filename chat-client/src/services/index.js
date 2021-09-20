import axios from "axios";

const login = async (user) => {
  const { email, password } = user;

  const { data } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/user/login`,
    {
      email,
      password,
    }
  );

  return data;
};

const register = async (user) => {
  const { email, password, firstName, lastName, role } = user;

  const { data } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/user`,
    {
      email,
      password,
      firstName,
      lastName,
      role,
    }
  );

  return data;
};

const getAdminMessages = async (user) => {
  const { token, email } = user;
  const { data } = await axios.get(
    `${process.env.REACT_APP_SERVER_URL}/channel?repEmail=${email}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const getUserMessages = async (user) => {
  const { token, email } = user;

  const { data } = await axios.get(
    `${process.env.REACT_APP_SERVER_URL}/channel?userEmail=${email}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const getActiveMessages = async (token) => {
  const { data } = await axios.get(
    `${process.env.REACT_APP_SERVER_URL}/channel?currentStatus=ACTIVE`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const getChannel = async (id, token) => {
  const { data } = await axios.get(
    `${process.env.REACT_APP_SERVER_URL}/channel/${id}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const updateChannelStatus = async (requestBody) => {
  const { data } = await axios.put(
    `${process.env.REACT_APP_SERVER_URL}/channel/${requestBody.id}`,
    requestBody,
    {
      headers: {
        Authorization: `Bearer ${requestBody.token}`,
      },
    }
  );

  return data;
};

const createChannel = async (requestBody, token) => {
  const { data } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/channel`,
    requestBody,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const createMessage = async (requestBody, token) => {
  const { data } = await axios.post(
    `${process.env.REACT_APP_SERVER_URL}/message`,
    requestBody,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

const getUser = async (email, token) => {
  const data = await axios.get(
    `${process.env.REACT_APP_SERVER_URL}/user/${email}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return data;
};

export {
  login,
  getUser,
  register,
  getAdminMessages,
  getUserMessages,
  getActiveMessages,
  getChannel,
  updateChannelStatus,
  createChannel,
  createMessage,
};
