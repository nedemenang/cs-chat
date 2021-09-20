import React, { createContext, useEffect, useReducer } from "react";
import Reducer from "./Reducer.js";

const initialState = {
  user: JSON.parse(localStorage.getItem("user")) || {},
  userChannels: [],
  selectedChannel: JSON.parse(localStorage.getItem("selectedChannel")) || {},
  isLoggedIn: false,
  error: null,
};

const Store = ({ children }) => {
  const [state, dispatch] = useReducer(Reducer, initialState);

  useEffect(()=>{
    localStorage.setItem("user", JSON.stringify(state.user))
  },[state.user])

  useEffect(()=>{
    localStorage.setItem("selectedChannel", JSON.stringify(state.selectedChannel))
  },[state.selectedChannel])

  return (
    <Context.Provider value={[state, dispatch]}>{children}</Context.Provider>
  );
};

export const Context = createContext(initialState);
export default Store;
