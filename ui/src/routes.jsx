import React from 'react';
import Search from "./search";

import {
  Route,
  Switch,
  BrowserRouter,
} from 'react-router-dom';

const Routes = () => {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path="/" component={Search}/>
        <Route exact path="/:imageName" component={Search}/>
      </Switch>
    </BrowserRouter>
  )
};

export default Routes
