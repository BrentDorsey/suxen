import React, {lazy, Suspense} from 'react';
import ApolloClient from 'apollo-client';
import {ApolloProvider} from 'react-apollo';
import {InMemoryCache} from 'apollo-cache-inmemory';
import {HttpLink} from "apollo-link-http";

const Routes = lazy(() => import('./Routes'));

const client = new ApolloClient({
  link: new HttpLink({uri: "/query"}),
  cache: new InMemoryCache(),
});

const Suxen = () => {
  return (
    <ApolloProvider client={client}>
      <Suspense fallback={<div>Loading...</div>}>
        <Routes />
      </Suspense>
    </ApolloProvider>
  )
}

export default Suxen
