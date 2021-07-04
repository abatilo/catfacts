import { lazy, Suspense } from "react";
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect,
} from "react-router-dom";

const HomePage = lazy(() => import("./Pages/HomePage"));

function App() {
  return (
    <Router>
      <Suspense fallback={<></>}>
        <Switch>
          <Route exact path="/">
            <HomePage />
          </Route>
          <Redirect to="/" />
        </Switch>
      </Suspense>
    </Router>
  );
}

export default App;
