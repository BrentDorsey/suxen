import React, {lazy, Suspense} from 'react';
import style from './app.scss';

import Suxen from './suxen';
const SuxenLazy = lazy(() => require('./suxen'));

const App = () => {
  return (
    <div className={style.root}>
      <div className={style.column}>
        <Suspense fallback={<div>app.jsx Loading...</div>}>
          <SuxenLazy />
        </Suspense>
      </div>
    </div>
  )
}

export default App
