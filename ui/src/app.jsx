import React from 'react';
import ApolloClient from 'apollo-boost';
import {ApolloProvider} from 'react-apollo';
import style from './app.scss';
import Search from "./search";

import {
  Route,
  Switch,
  BrowserRouter,
} from 'react-router-dom';


const client = new ApolloClient({uri: '/query'});

const App = () => {
  return (
    <div className={style.root}>
      <div className={style.column}>
        <ApolloProvider client={client}>
          <BrowserRouter>
            <Switch>
              <Route exact path="/" component={Search}/>
              <Route exact path="/:imageName" component={Search}/>
            </Switch>
          </BrowserRouter>
        </ApolloProvider>
      </div>
    </div>
  )

}

export default App
