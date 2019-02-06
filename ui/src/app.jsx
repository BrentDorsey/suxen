import React from 'react';
import ApolloClient from 'apollo-client';
import {ApolloProvider} from 'react-apollo';
import { InMemoryCache } from 'apollo-cache-inmemory';
import style from './app.scss';
import Search from "./search";

import {
  Route,
  Switch,
  BrowserRouter,
} from 'react-router-dom';
import {HttpLink} from "apollo-boost";


const client = new ApolloClient({
 link: new HttpLink({ uri: "query" }),
 cache: new InMemoryCache(),
});

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
