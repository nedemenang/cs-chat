const Reducer = (state, action) => {
    switch (action.type) {
      case "SET_LOGGED_IN_USER":
        return {
          ...state,
          user: action.payload,
        };
      case "SET_SELECTED_CHANNEL":
        return {
          ...state,
          selectedChannel: action.payload,
        };
      case "SET_IS_LOGGED_IN":
        return {
          ...state,
          isLoggedIn: action.payload,
        };
      case "SET_ERROR":
        return {
          ...state,
          error: action.payload,
        };
      default:
        return state;
    }
  };
  
  export default Reducer;
  