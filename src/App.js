import { React, useState, Component, useRef } from "react";
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import Articles from './components/articles';

function App () {
  const baseURL = "http://localhost:10000/article/"
  const get_keyword = useRef(null);
  const [getResult, setGetResult] = useState(null);

  async function getDataByKeyword() {

    const state = {
      articles: [],
    }

    const keyword = get_keyword.current.value;
    if (keyword) {
      try {
        const full_URL = 'http://localhost:10000/article/'+keyword;
        const res = await fetch(full_URL);
        if (!res.ok) {
          const message = "An error has occurred";
          throw new Error (message);
        }

        const data = await res.json();
        state.articles = data;
        setGetResult(state.articles)
    } catch (err) {
        setGetResult(err.message);
    }
  }
}

const clearGetOutput = () => {
    setGetResult(null);
  }

return (
    <div className="card">
      <div className="card-header">Search NYT Postgres</div>
      <div className="card-body">
        <div className="input-group input-group-sm">
          <input type="text" ref={get_keyword} className="form-control ml-2" placeholder="Keyword" />
          <div className="input-group-append">
            <button className="btn btn-sm btn-primary" onClick={getDataByKeyword}>Get by Keyword</button>
          </div>
          <button className="btn btn-sm btn-warning ml-2" onClick={clearGetOutput}>Clear</button>
        </div>
        { getResult && <Articles articles={getResult}/> }
      </div>
    </div>
  );
}

export default App;
