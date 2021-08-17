import React from "react";
import "./App.scss";
import Dashboard from "./Dashboard/Dashboard";
import { BrowserRouter as Router, Redirect, Route, Switch } from "react-router-dom";
import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";
import Login from "./Login/Login";
import PrivateRoute from "./PrivateRoute";
import { CookiesProvider, useCookies } from "react-cookie";
import { useAuthUserStorage } from "./Login/UseUserStorage.hooks";

function App() {
  const [authUser] = useAuthUserStorage();

  const client = new ApolloClient({
    uri: "/api/query",
    cache: new InMemoryCache(),
  });

  const isLoggedIn = !!authUser;
  
  return (
    <ApolloProvider client={client}>
      <CookiesProvider>
        <Router>
          <Switch>
            <Route path="/login">
              {
                isLoggedIn
                  ? <Redirect to={{ pathname: '/' }} />
                  : <Login />
              }
            </Route>
            <PrivateRoute path="/" validateLogin={true} component={Dashboard}>
            </PrivateRoute>
          </Switch>
        </Router>
      </CookiesProvider>
    </ApolloProvider>
  );
}

export default App;
