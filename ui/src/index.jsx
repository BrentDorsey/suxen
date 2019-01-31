import React from 'react';
import ReactDOM from 'react-dom'
import App from './app'

const wrapper = document.createElement('div')
document.body.appendChild(wrapper)

ReactDOM.render(<App />, wrapper)
