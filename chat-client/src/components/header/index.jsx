import React, { useContext } from "react";
import { Link, useHistory, Redirect } from "react-router-dom";
import { Context } from "../../Store.js";
import "./styles.css";

const Index = () => {
  const history = useHistory();
  const [state, dispatch] = useContext(Context);

  const logout = async () => {
    try {
      dispatch({ type: "SET_LOGGED_IN_USER", payload: {} });
      dispatch({ type: "SET_IS_LOGGED_IN", payload: false });
      history.push("/");
    } catch (error) {
      console.log("error signing out: ", error);
    }
  };

  if (!state.user) {
    return <Redirect to={"/"} />;
  } else {
    return (
        <React.Fragment>
          <div className="headerContainer">
            <div className="headerLeft">
              <span className="headerTitle">CS Chat App</span>
            </div>
            <div className="headerCenter"></div>
            <div className="headerRight">
            {!state.user ? (
              <React.Fragment>
              <div>
                <Link to="/signup" className="headerLink">
                  Signup
                </Link>
                </div>
                <div>
                <Link to="/" className="headerLink">
                  Login
                </Link>
              </div>
          </React.Fragment>
          ) : (
            <React.Fragment>
            <div>
                <Link to="/" onClick={logout} className="headerLink">
                  Logout
                </Link>
                </div>
                </React.Fragment>
          )}
          </div>
          </div>
        </React.Fragment>
      );
    }
};

export default Index;
