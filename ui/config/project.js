// Config API, for adding reducers and configuring our ReactQL app
import config from 'kit/config';

if (SERVER) {
  /* GRAPHQL */

  // Enable the GraphQL server by passing in the schema.  Since we're running
  // in an `if (SERVER)` block, we can be sure that `require()`'ing
  // the schema won't 'bloat' our browser bundle
  config.enableGraphQLServer(require('src/graphql/schema').default);
}
