import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// 這個 store 用於儲存 JWT token。
// 我們從 localStorage 初始化它，這樣即使重新整理頁面，登入狀態也能保持。
const initialToken = browser ? window.localStorage.getItem('jwt_token') : null;
export const authToken = writable(initialToken);

// 當 authToken 的值發生變化時，自動更新 localStorage。
authToken.subscribe(token => {
  if (browser) {
    if (token) {
      window.localStorage.setItem('jwt_token', token);
    } else {
      window.localStorage.removeItem('jwt_token');
    }
  }
});

// 一個方便的 store，用來判斷使用者是否已登入。
export const isAuthenticated = writable(!!initialToken);
authToken.subscribe(token => {
    isAuthenticated.set(!!token);
});

// 登出函式
export function logout() {
    authToken.set(null);
}
