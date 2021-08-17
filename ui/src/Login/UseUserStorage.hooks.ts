import { useState } from "react";
import Cookies from 'universal-cookie';

const key = 'user-storage';

interface AuthUser {
  id: string
  email: string
  firstName: string
  lastName: string
  isAdmin: string
  clubs: [{
    id: string
    clubName: string
  }]
}

// Hook
// TODO : Why cant TS detect the tuple type?
export function useAuthUserStorage(initialValue: any = null, maxAge: number = 60 * 60 * 4): [AuthUser, (value: AuthUser) => void, () => void] {
  // State to store our value
  // Pass initial state function to useState so logic is only executed once
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const cookies = new Cookies();
      // Get from local storage by key
      const item = cookies.get(key);
      // Parse stored json or if none return initialValue
      const returnValue = item ? item : initialValue;
      return returnValue as AuthUser
    } catch (error) {
      // If error also return initialValue
      console.log(error);
      return initialValue as AuthUser;
    }
  });
  
  // Return a wrapped version of useState's setter function that ...
  // ... persists the new value to cookies.
  const setValue = (value: AuthUser) => {
    try {
      const cookies = new Cookies();
      // Allow value to be a function so we have same API as useState
      const valueToStore =
        value instanceof Function ? value(storedValue) : value;
      // Save state
      setStoredValue(valueToStore);
      const options = (maxAge) ? 
        { maxAge }
        : {};
      // Save to local storage
      cookies.set(key, valueToStore, options);
    } catch (error) {
      // A more advanced implementation would handle the error case
      console.log(error);
    }
  };

  const reset = () => {
    setStoredValue(initialValue);
    const cookies = new Cookies();
    cookies.remove(key);
  };

  return [storedValue, setValue, reset];
}