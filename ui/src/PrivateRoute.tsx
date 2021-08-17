// This is used to determine if a user is authenticated and
// if they are allowed to visit the page they navigated to.

// If they are: they proceed to the page
// If not: they are redirected to the login page.
import React from 'react'
import { useCookies } from 'react-cookie';
import { Redirect, Route } from 'react-router-dom'
import { useAuthUserStorage } from './Login/UseUserStorage.hooks';

const PrivateRoute = ({ component: Component, validateLogin, ...rest }) => {
  const [authUser] = useAuthUserStorage();

  const isLoggedIn = (authUser)
    ? Boolean(authUser)
    : Boolean(authUser)

  return (
    <Route
      {...rest}
      render={props =>
        isLoggedIn ? (
          <Component {...props} />
        ) : (
          <Redirect to={{ pathname: '/login', state: { from: props.location } }} />
        )
      }
    />
  )
}

export default PrivateRoute