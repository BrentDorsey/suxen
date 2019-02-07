import React from 'react';
import ReactDOM from 'react-dom'
import App from './app';

// const App = React.lazy(() => import('./app'));

const wrapper = document.createElement('div')
document.body.appendChild(wrapper)

const Main = () => <React.Suspense fallback={<div>index.jsx Loading...</div>}>
  <App />
</React.Suspense>;

ReactDOM.render(<App />, wrapper)
