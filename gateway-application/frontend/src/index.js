import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import './styles/index.css';

import Home from './pages/Home';
import Account from './pages/Account';
import Chequing from './pages/Chequing';
import Payment from './pages/Payment';
import Deposit from './pages/Deposit';
import Withdraw from './pages/Withdraw';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Home />,
  },
  {
    path: '/account/:account_id',
    element: <Account />
  },
  {
    path: '/cheq/:account_id',
    element: <Chequing />
  },
  {
    path: '/payment',
    element: <Payment />
  },
  {
    path: '/deposit',
    element: <Deposit />
  },
  {
    path: '/withdraw',
    element: <Withdraw />
  }
]);

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router}>
    </RouterProvider>
  </React.StrictMode>
);
