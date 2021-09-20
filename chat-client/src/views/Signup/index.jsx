import React, { useState, useContext } from "react";
import { Link, Redirect } from "react-router-dom";
import { notification, Spin } from "antd";
import { LoadingOutlined } from "@ant-design/icons";
import axios from "axios";
import "./style.css";
import "antd/dist/antd.css";
import { Context } from "../../Store.js";
const { register } = require("../../services/index");

function Signup() {
  const [state, dispatch] = useContext(Context);
  const [user, setUser] = useState({
    userName: "",
    firstName: "",
    lastName: "",
    password: "",
    confirmPassword: "",
    email: "",
    role: "",
    redirect: null,
    isLoading: false,
  });

  const [redirect, setRedirect] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const signUp = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    const {
      password,
      email,
      firstName,
      lastName,
      confirmPassword,
      role,
    } = user;
    if (password !== confirmPassword) {
      notification.error({
        message: "Error signing up user!",
        description: "Password and confirm password do not match",
        placement: "topRight",
        duration: 2.0,
        onClose: () => {
          setRedirect(null);
          setIsLoading(false);
        },
      });
      return;
    }

    try {
      await register(user);
      setIsLoading(false);
      notification.success({
        message: "Succesfully signed up user!",
        description: "Account created successfully, Redirecting you in a few!",
        placement: "topRight",
        duration: 2.0,
        onClose: () => {
          setUser({
            email: "",
            firstName: "",
            lastName: "",
            password: "",
            confirmPassword: "",
            role: "",
          });
          setRedirect("/");
          setIsLoading(false);
        },
      });
    } catch (error) {
      notification.error({
        message: "Error signing up user!",
        description: error.message,
        placement: "topRight",
        duration: 2.0,
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
                className="registerBox"
                onSubmit={signUp}
                noValidate=""
              >
                <input
                  id="firstName"
                  onChange={(e) =>
                    setUser({ ...user, firstName: e.target.value })
                  }
                  type="text"
                  value={user.firstName}
                  className="loginInput"
                  name="firstName"
                  placeholder="First Name"
                  required
                  autoFocus
                />
                <input
                  id="lastName"
                  onChange={(e) =>
                    setUser({ ...user, lastName: e.target.value })
                  }
                  type="text"
                  value={user.lastName}
                  className="loginInput"
                  name="lastName"
                  placeholder="Last Name"
                  required
                />
                <input
                  id="email"
                  className="loginInput"
                  placeholder="Email"
                  onChange={(e) => setUser({ ...user, email: e.target.value })}
                  value={user.email}
                  type="email"
                  name="email"
                  required
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
                  placeholder="Password"
                  required
                />

                <input
                  id="confirmPassword"
                  onChange={(e) =>
                    setUser({ ...user, confirmPassword: e.target.value })
                  }
                  value={user.confirmPassword}
                  type="password"
                  className="loginInput"
                  name="confirmPassword"
                  placeholder="Confirm Password"
                  required
                />

                <select
                  id="role"
                  value={user.role}
                  onChange={(e) => setUser({ ...user, role: e.target.value })}
                  className="loginInput"
                >
                  <option selected>Select Role</option>
                  <option value="ADMIN">Admin</option>
                  <option value="USER">User</option>
                </select>

                <button type="submit" className="loginButton">
                  {isLoading ? (
                    <Spin
                      indicator={<LoadingOutlined style={{ fontSize: 24 }} />}
                    />
                  ) : (
                    "Register"
                  )}
                </button>
                <div className="loginSignup">
                  Don't have an account? <Link to="/">Sign Up</Link>
                </div>
              </form>
            </div>
          </div>
        </div>
      </React.Fragment>
    );
  }
}

export default Signup;
