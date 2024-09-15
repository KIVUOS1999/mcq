import Question from "./Question/Question"
import WaitingRoom from "./Question/WaitingRoom";
import Room from './Rooms/Room'
import Submit from './Submit'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route exact path="/" element={<Room />}/>
          <Route path="/question/:roomId/:playerId/:time/:admin" element={<Question />}/>
          <Route path="/lobby/:roomId/:playerId/:isAdmin/:time/:questions" element={<WaitingRoom />}/>
          <Route path="/submit/:roomId" element={<Submit />}/>
        </Routes>
      </div>
    </Router>
  );
}

export default App;
