import React, { useState, useContext } from "react";
import { Link, Redirect } from "react-router-dom";
import { notification, Spin } from "antd";
import { LoadingOutlined } from "@ant-design/icons";
import "./style.css";
import { Context } from "../../Store.js";
const { login } = require("../../services/index");

function Login() {
  const [state, dispatch] = useContext(Context);
  const [user, setUser] = useState({
    email: "",
    password: "",
    role: "",
    firstName: "",
    lastName: "",
  });

  const [redirect, setRedirect] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const signIn = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const data = await login(user);
      setIsLoading(false);
      dispatch({ type: "SET_LOGGED_IN_USER", payload: data });
      dispatch({ type: "SET_IS_LOGGED_IN", payload: true });
      notification.success({
        message: "Succesfully logged in!",
        description: "Successfully logged in, Redirecting you in a few!",
        placement: "topRight",
        duration: 2.0,
        onClose: () => {
          if (data.role === "ADMIN") {
            setRedirect("/message-dashboard");
          } else {
            setRedirect("/dashboard");
          }
          setUser({ email: "", password: "" });
        },
      });
    } catch (error) {
      notification.error({
        message: "Error logging in!",
        description: error.message,
        placement: "topRight",
        duration: 1.0,
        onClose: () => {
          setRedirect(null);
          setIsLoading(false);
        },
      });
    }
  };
  if (redirect) {
    return <Redirect to={redirect} />;
  } else {
    return (
      <React.Fragment>
        <div className="login">
          <div className="loginWrapper">
            <div className="loginLeft">
              <h3 className="loginLogo">CS Chat APP</h3>
            </div>
            <div className="loginRight">
              <form
                method="POST"
                className="loginBox"
                onSubmit={signIn}
                noValidate=""
              >
                <input
                  id="email"
                  className="loginInput"
                  placeholder="email"
                  onChange={(e) => setUser({ ...user, email: e.target.value })}
                  value={user.email}
                  name="email"
                  required
                  autoFocus
                />

                <input
                  id="password"
                  type="password"
                  onChange={(e) =>
                    setUser({ ...user, password: e.target.value })
                  }
                  value={user.password}
                  className="loginInput"
                  name="password"
                  placeholder="password"
                  required
                />
                <button type="submit" className="loginButton">
                  {isLoading ? (
                    <Spin
                      indicator={<LoadingOutlined style={{ fontSize: 24 }} />}
                    />
                  ) : (
                    "Login"
                  )}
                </button>
                <div className="loginSignup">
                  Don't have an account? <Link to="/signup">Sign Up</Link>
                </div>
              </form>
            </div>
          </div>
        </div>
      </React.Fragment>
    );
  }
}

export default Login;
