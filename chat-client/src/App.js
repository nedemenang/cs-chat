import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import "./App.css";
import Store from "./Store.js";
import Header from "./components/header/index.jsx";
// import { useContext } from "react";
import SignupPage from "./views/Signup/index.jsx";
import LoginPage from "./views/Login/index.jsx";
import MessageDashboard from "./views/MessageDashboard/index.jsx";
import UserDashboard from "./views/UserDashboard/index.jsx";

function App() {
  // const state = useContext(Context);
  return (
    <Store>
      <Router>
        <Route component={Header} />
        <Switch>
          <Route path="/" exact component={LoginPage} />
          <Route path="/signup" component={SignupPage} />
          <Route path="/dashboard" component={UserDashboard} />
          <Route path="/message-dashboard" component={MessageDashboard} />
        </Switch>
        {/* <Footer /> */}
      </Router>
    </Store>
  );
}
export default App;
